package consul

/*封装consul客户端工具，支持配置注册与服务发现*/

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/store_server/logger"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

//consul driver
type ConsulDriver struct {
	sync.RWMutex
	ctx        context.Context
	client     *api.Client
	session    *api.Session
	ServiceMap map[string][]*ServiceInfo
}

//service info
type ServiceInfo struct {
	ServiceID string
	IP        string
	Port      int
}

func NewConsulDriver(ctx context.Context, address, scheme string) (cd *ConsulDriver, err error) {
	if scheme == "" {
		scheme = "http" //default http
	}
	conf := &api.Config{
		Address: address,
		Scheme:  scheme,
	}
	cd = &ConsulDriver{}
	cd.Lock()
	defer cd.Unlock()
	cd.ctx, cd.ServiceMap = ctx, make(map[string][]*ServiceInfo)
	cd.client, err = api.NewClient(conf)
	if err != nil {
		return
	}
	cd.session = cd.client.Session()
	return
}

func (cd *ConsulDriver) RegisterKV(key string, val []byte) (err error) { //注册key value到consul
	p := &api.KVPair{Key: key, Value: val}
	_, err = cd.client.KV().Put(p, &api.WriteOptions{})
	if err != nil {
		return
	}
	return
}

func (cd *ConsulDriver) LoadKV(key string) (val []byte, err error) { //从consul加载指定key对应的value
	var pair *api.KVPair
	pair, _, err = cd.client.KV().Get(key, &api.QueryOptions{})
	if err != nil {
		return
	}
	if pair.Key != key {
		err = fmt.Errorf("load kv pair from consul not match, query key: %v|load key: %v", key, pair.Key)
		return
	}
	val = pair.Value
	return
}

func (cd *ConsulDriver) UpdateKV(key string, val []byte) (err error) { //更新key value
	var (
		ok      bool
		session string
	)
	session, _, err = cd.session.Create(nil, &api.WriteOptions{})
	if err != nil {
		return
	}
	pair := &api.KVPair{
		Key:     key,
		Session: session,
	}
	ok, _, err = cd.client.KV().Acquire(pair, &api.WriteOptions{})
	if !ok {
		return
	}
	defer func() {
		ok, _, err = cd.client.KV().Release(pair, &api.WriteOptions{})
		if !ok {
			return
		}
	}()
	err = cd.DeleteKV(key)
	if err != nil {
		return
	}
	err = cd.RegisterKV(key, val)
	return
}

func (cd *ConsulDriver) DeleteKV(key string) (err error) { //删除key指定的KV Pair
	_, err = cd.client.KV().Delete(key, &api.WriteOptions{})
	return
}

func (cd *ConsulDriver) RegisterConfigFile(file, key string) (err error) { //将指定的配置文件内容注册到KV, key默认为配置文件名
	if _, err = os.Stat(file); os.IsNotExist(err) {
		return
	}
	ext := filepath.Ext(file)
	fn := filepath.Base(file)
	if len(key) == 0 {
		key = strings.TrimSuffix(fn, ext)
	}
	var content []byte
	content, err = ioutil.ReadFile(file)
	if err != nil {
		return
	}
	return cd.RegisterKV(key, content)
}

func (cd *ConsulDriver) WatchConfigChange(key string, fn func(*[]byte) error) (err error) { //监测配置变更，方便重新加载更新后的配置
	var lastV []byte
	lastV, err = cd.LoadKV(key) //获取初始配置值
	if err != nil {
		return
	}
	tk := time.NewTicker(time.Second) //每秒获取一次，对比配置内容
	defer tk.Stop()
	go func(initV []byte) {
		for {
			select {
			case <-tk.C:
			case <-cd.ctx.Done():
				return
			}

			var v []byte
			v, err = cd.LoadKV(key)
			if err != nil {
				logger.Entry().Errorf("load config content from consul by key: %v|error: %v", key, err)
				continue
			}
			if len(v) == 0 {
				logger.Entry().Errorf("load config content is empty by key: %v .....", key)
				continue
			}
			if md5.Sum(v) != md5.Sum(initV) {
				initV = v //有更新，赋予最新值
				if fn != nil {
					err = fn(&v) //通知监测端重载配置
					if err != nil {
						logger.Entry().Errorf("notify listen client to reload config by key: %v|error: %v", key, err)
					}
				}
			}
		}
	}(lastV)
	return
}

func (cd *ConsulDriver) RegisterService(host string, port int, name string,
	interval time.Duration, opts ...string) (err error) { //注册服务
	svr := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%v-%v-%v", name, host, port),
		Name:    name,
		Address: host,
		Port:    port,
		Check: &api.AgentServiceCheck{ //health check
			Interval: interval.String(),
			HTTP:     fmt.Sprintf("%v:%v/%v", host, port, "status"),
		},
	}
	tag := ""
	if len(opts) > 0 {
		tag = opts[0]
	}
	svr.Tags = []string{tag}
	err = cd.client.Agent().ServiceRegister(svr)
	return
}

func (cd *ConsulDriver) DeregisterService(id string) (err error) { //注销服务
	if id == "" {
		err = fmt.Errorf("can't deregister service id which is empty...")
	}
	err = cd.client.Agent().ServiceDeregister(id)
	return
}

func (cd *ConsulDriver) DiscoverService(service_name string, healthOnly bool) error { //发现服务
	if _, ok := cd.ServiceMap[service_name]; ok { //已缓存该服务
		return nil
	}
	services, _, err := cd.client.Catalog().Services(&api.QueryOptions{})
	if err != nil {
		return err
	}
	svrs := make([]*ServiceInfo, 0)
	for name := range services {
		servicesData, _, e := cd.client.Health().Service(name, "", healthOnly, &api.QueryOptions{})
		if e != nil {
			return e
		}
		for _, entry := range servicesData {
			if service_name != entry.Service.Service {
				continue
			}
			for _, health := range entry.Checks {
				if health.ServiceName != service_name {
					continue
				}
				logger.Entry().Info("health nodeid: ", health.Node, " service_name: ", health.ServiceName,
					" service_id: ", health.ServiceID, " status: ", health.Status, " ip: ", entry.Service.Address,
					" port: ", entry.Service.Port)
				node := &ServiceInfo{
					IP:        entry.Service.Address,
					Port:      entry.Service.Port,
					ServiceID: health.ServiceID,
				}
				svrs = append(svrs, node)
			}
		}
	}
	cd.Lock()
	cd.ServiceMap[service_name] = svrs //同名服务，通过ID区分
	cd.Unlock()
	return nil
}
