package g

import (
	"bytes"
	"fmt"
	"github.com/store_server/store_server_http/conf"
	"github.com/store_server/utils/common"
	"gopkg.in/yaml.v2"
	"net/http"
	//"github.com/mitchellh/mapstructure"
	"os"
	"sync"
)

var (
	config = conf.NewDefaultStoreServerHttpConfig()
	lock   = sync.RWMutex{}
)

func ParseConfig(data []byte) (err error) {
	c := conf.NewDefaultStoreServerHttpConfig()
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return
	}
	lock.Lock()
	config = c
	lock.Unlock()
	return
}

func Config() *conf.StoreServerHttpConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func InitConf(fileName string, debug bool) {
	// 优先从公司配置中心加载
	f := func() {
		err := common.ParseYamlConfigFile(fileName, config)
		if err != nil {
			panic(err)
		}
	}
	if !debug {
		err := LoadConfigFromRainbow("", "")
		if err != nil { //从配置中心加载失败
			os.Stderr.WriteString(fmt.Sprintf("load config from rainbow error: %v\n", err))
			f()
		}
	} else {
		f()
	}
}

func DoPost(api string, payload []byte) (code int, err error) {
	var (
		req *http.Request
		rsp *http.Response
	)
	req, err = http.NewRequest(http.MethodPost, api, bytes.NewBuffer(payload))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	rsp, err = http.DefaultClient.Do(req)
	if err != nil {
		code = rsp.StatusCode
		return
	}
	defer rsp.Body.Close()
	return
}

func LoadConfigFromRainbow(appId, group string) error {
	return nil
}
