package conf

import (
	"time"
)

//store server config
type StoreServerHttpConfig struct {
	Http         HttpConfig            `json:"http,omitempty" yaml:"http"`
	IpWhiteList  string                `json:"ip_white_list,omitempty" yaml:"ip_white_list"`
	Mysql        string                `json:"mysql,omitempty" yaml:"mysql"`
	MongoDb      MongoDB               `json:"mongodb" yaml:"mongodb"`
	Es           EsConfig              `json:"es,omitempty" yaml:"es"`
	Es7          EsConfig              `json:"es7,omitempty" yaml:"es7"`
	Influx       InfluxDB              `json:"influxdb,omitempty" yaml:"influxdb"`
	Dataplatform DataplatformSearchAPI `json:"dataplatform_search" yaml:"dataplatform_search"`
	CMQ          CMQConfig             `json:"cmq" yaml:"cmq"`
	//初始化时全量数据导出开关
	ExportAllOpen bool      `json:"export_all_open" yaml:"export_all_open"`
	Cls           ClsConfig `json:"cls" yaml:"cls"`
	ValidRegions  []int     `json:"valid_regions" yaml:"valid_regions"`
}

//http config
type HttpConfig struct {
	Listen      string `json:"listen,omitempty" yaml:"listen"`
	Debug       bool   `json:"debug,omitempty" yaml:"debug"`
	HttpTimeout int    `json:"http_timeout,string" yaml:"http_timeout"`
}

//mongo config
type MongoDB struct {
	ModId    int    `json:"mod_id" yaml:"mod_id"`
	CmdId    int    `json:"cmd_id" yaml:"cmd_id"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	DbName   string `json:"db_name" yaml:"db_name"`
	User     string `json:"user" yaml:"user"`
	Passwd   string `json:"passwd" yaml:"passwd"`
	TimeOut  int    `json:"time_out" yaml:"time_out"`
	PoolSize int    `json:"pool_size" yaml:"pool_size"`
	Direct   bool   `json:"direct" yaml:"direct"`
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

//数据平台配置
type DataplatformSearchAPI struct {
	ModId      int    `json:"mod_id" yaml:"mod_id"`
	CmdId      int    `json:"cmd_id" yaml:"cmd_id"`
	Api        string `json:"api" yaml:"api"`
	Host       string `json:"host" yaml:"host"`
	Port       int    `json:"port" yaml:"port"`
	TrackType  int    `json:"track_type" yaml:"track_type"`
	AlbumType  int    `json:"album_type" yaml:"album_type"`
	SingerType int    `json:"singer_type" yaml:"singer_type"`
}

//cmq config
type CMQConfig struct {
	Domain    string   `json:"domain" yaml:"domain"`
	Path      string   `json:"path" yaml:"path"`
	SecretID  string   `json:"secret_id" yaml:"secret_id"`
	SecretKey string   `json:"secret_key" yaml:"secret_key"`
	Region    string   `json:"region" yaml:"region"`
	Queue     []string `json:"queue" yaml:"queue"`
	Topic     string   `json:"topic" yaml:"topic"`
	Delay     int      `json:"delay" yaml:"delay"`
	Timeout   int      `json:"timeout" yaml:"timeout"`
	MsgEnv    int      `json:"msg_env" yaml:"msg_env"`
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

//influx db config
type InfluxDB struct {
	Host        string `json:"host,omitempty" yaml:"host"`
	Port        int    `json:"port,string" yaml:"port"`
	Database    string `json:"database,omitempty" yaml:"database"`
	Measurement string `json:"measurement,omitempty" yaml:"measurement"`
}

type ModOption func(*StoreServerHttpConfig) ModOption

func (c *StoreServerHttpConfig) Option(opts ...ModOption) (previous ModOption) {
	for _, opt := range opts {
		previous = opt(c)
	}
	return previous
}

func SetHttpConfig(cg HttpConfig) ModOption {
	return func(c *StoreServerHttpConfig) ModOption {
		previous := c.Http
		c.Http = cg
		return SetHttpConfig(previous)
	}
}

func SetMongoConfig(cg MongoDB) ModOption {
	return func(c *StoreServerHttpConfig) ModOption {
		previous := c.MongoDb
		c.MongoDb = cg
		return SetMongoConfig(previous)
	}
}

func SetEsConfig(cg EsConfig) ModOption {
	return func(c *StoreServerHttpConfig) ModOption {
		previous := c.Es
		c.Es = cg
		return SetEsConfig(previous)
	}
}

func SetInfluxConfig(cg InfluxDB) ModOption {
	return func(c *StoreServerHttpConfig) ModOption {
		previous := c.Influx
		c.Influx = cg
		return SetInfluxConfig(previous)
	}
}

func NewDefaultStoreServerHttpConfig() *StoreServerHttpConfig {
	return &StoreServerHttpConfig{
		Http: HttpConfig{},
		CMQ:  CMQConfig{},
	}
}
