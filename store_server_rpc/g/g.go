package g

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"sync"

	"github.com/store_server/store_server_rpc/conf"
	"github.com/store_server/utils/common"
)

var (
	config = conf.NewDefaultStoreServerRpcConfig()
	lock   = sync.RWMutex{}
)

func ParseConfig(data []byte) (err error) {

	c := conf.NewDefaultStoreServerRpcConfig()
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return
	}
	lock.Lock()
	config = c
	lock.Unlock()
	return
}

func Config() *conf.StoreServerRpcConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func InitConf(fileName string, debug bool) {
	// 优先从配置中心加载
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

func LoadConfigFromRainbow(appId, group string) error {
	return nil
}
