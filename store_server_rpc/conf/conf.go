package conf

import (
	"time"
)

//store server config
type StoreServerRpcConfig struct {
	Mysql   string    `json:"mysql" yaml:"mysql"`
	MongoDb MongoDB   `json:"mongodb" yaml:"mongodb"`
	Es      EsConfig  `json:"es,omitempty" yaml:"es"`
	Es7     EsConfig  `json:"es7,omitempty" yaml:"es7"`
	Rpc     RpcConfig `json:"rpc"`
	RpcPort int       `json:"rpc_port,omitempty" yaml:"rpc_port"`
	Cls     ClsConfig `json:"cls" yaml:"cls"`
}

//mongo db config
type MongoDB struct {
	ModId          int    `json:"mod_id" yaml:"mod_id"`
	CmdId          int    `json:"cmd_id" yaml:"cmd_id"`
	Host           string `json:"host" yaml:"host"`
	Port           int    `json:"port" yaml:"port"`
	DbName         string `json:"db_name" yaml:"db_name"`
	User           string `json:"user" yaml:"user"`
	Passwd         string `json:"passwd" yaml:"passwd"`
	TimeOut        int    `json:"time_out" yaml:"time_out"`
	PoolSize       int    `json:"pool_size" yaml:"pool_size"`
	Direct         bool   `json:"direct" yaml:"direct"`
	AudioCol       string `json:"audio_col" yaml:"audio_col"`
	ExtResourceCol string `json:"external_resource_col" yaml:"external_resource_col"`
}

//rpc config
type RpcConfig struct {
	Listen string `json:"listen" yaml:"listen"`
	Debug  bool   `json:"debug" yaml:"debug"`
}

//es config
type EsConfig struct {
	Address       []string      `json:"address,omitempty" yaml:"address"`
	Sniff         bool          `json:"sniff,omitempty" yaml:"sniff"`
	DisableSync   bool          `json:"disable_sync,omitempty" yaml:"disable_sync"`
	FlushInterval time.Duration `json:"flush_interval,omitempty" yaml:"flush_interval"`
	Timeout       int           `json:"timeout" yaml:"timeout"`
	Auth          *EsServerAuth `json:"auth" yaml:"auth"`
	Index         string        `json:"index,omitempty" yaml:"index"`
	Type          string        `json:"type,omitempty" yaml:"type"`
	Proxy         string        `json:"proxy,omitempty" yaml:"proxy"`
}

//es auth config
type EsServerAuth struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

//cls config
type ClsConfig struct {
	SecretId  string `json:"secret_id" yaml:"secret_id"`
	SecretKey string `json:"secret_key" yaml:"secret_key"`
	Region    string `json:"region" yaml:"region"`
	TopicId   string `json:"topic_id" yaml:"topic_id"`
	Open      bool   `json:"open" yaml:"open"`
	Async     bool   `json:"async" yaml:"async"`
}

type ModOption func(*StoreServerRpcConfig) ModOption

func (c *StoreServerRpcConfig) Option(opts ...ModOption) (previous ModOption) {
	for _, opt := range opts {
		previous = opt(c)
	}
	return previous
}

func SetRpcConfig(cg RpcConfig) ModOption {
	return func(c *StoreServerRpcConfig) ModOption {
		previous := c.Rpc
		c.Rpc = cg
		return SetRpcConfig(previous)
	}
}

func SetMongoConfig(cg MongoDB) ModOption {
	return func(c *StoreServerRpcConfig) ModOption {
		previous := c.MongoDb
		c.MongoDb = cg
		return SetMongoConfig(previous)
	}
}

func SetEsConfig(cg EsConfig) ModOption {
	return func(c *StoreServerRpcConfig) ModOption {
		previous := c.Es
		c.Es = cg
		return SetEsConfig(previous)
	}
}

func NewDefaultStoreServerRpcConfig() *StoreServerRpcConfig {
	return &StoreServerRpcConfig{
		Rpc: RpcConfig{},
	}
}
