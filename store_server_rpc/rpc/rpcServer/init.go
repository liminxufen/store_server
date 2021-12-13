package rpcServer

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/store_server/dbtools/dblogic"
	"github.com/store_server/dbtools/driver"
	ies "github.com/store_server/dbtools/elastic"
	ies7 "github.com/store_server/dbtools/elastic7"
	im "github.com/store_server/dbtools/mongo"
	"github.com/store_server/logger"
	"github.com/store_server/store_server_rpc/g"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

//db util define
type DBUtil struct {
	dbs          map[string]*gorm.DB
	rawdbs       map[string]*sql.DB
	mgoclient    *mongo.Client
	importclient *mongo.Client
	esclient     *ies.ESClient
	esclient7    *ies7.ESClient
	sync.RWMutex
	ctx context.Context
}

func InitDB(cfg string) (db *gorm.DB, err error) {
	return driver.CreateDB(cfg)
}

func InitRawDB(cfg string) (db *sql.DB, err error) {
	return driver.CreateRawDB(cfg)
}

func NewMongoClientOpts(host, database, user, passwd string) (opts *options.ClientOptions) { //mongo client配置参数
	addrs := fmt.Sprintf("%s:%d", g.Config().MongoDb.Host, g.Config().MongoDb.Port)
	if len(host) != 0 {
		addrs = host
	}
	dur := time.Duration(g.Config().MongoDb.TimeOut) * time.Second
	poolSize := uint64(g.Config().MongoDb.PoolSize)
	direct := g.Config().MongoDb.Direct
	opts = &options.ClientOptions{
		Hosts: []string{addrs},
		Auth: &options.Credential{
			//AuthMechanism: "MONGODB-CR",
			AuthMechanism: "SCRAM-SHA-1",
			AuthSource:    database,
			Username:      user,
			Password:      passwd,
		},
		MaxPoolSize:    &poolSize,
		Direct:         &direct,
		ConnectTimeout: &dur,
	}
	return opts
}

func InitMongo(modId, cmdId int, database, user, passwd string) (client *mongo.Client, err error) {
	var addrs string
	if err != nil { //优先从负载均衡获取mongodb服务器地址
		logger.Entry().Errorf("get cmongo server address by mod_id:%v |and cmd_id: %v|error: %v",
			modId, cmdId, err)
		addrs = fmt.Sprintf("%s:%d", g.Config().MongoDb.Host, g.Config().MongoDb.Port)
	} else {
		addrs = ""
	}
	opts := NewMongoClientOpts(addrs, database, user, passwd)
	return driver.CreateMongo(opts)
}

func InitElastic(ctx context.Context) (client *ies.ESClient, err error) {
	ec := g.Config().Es
	args := make([]string, 0)
	if auth := ec.Auth; auth != nil {
		args = append(args, auth.Username, auth.Password)
	}
	return ies.NewEsClient(ctx, ec.Address, ec.Timeout, ec.Sniff, ec.Proxy, args...)
}

func InitElastic7(ctx context.Context) (client *ies7.ESClient, err error) {
	ec := g.Config().Es7
	args := make([]string, 0)
	if auth := ec.Auth; auth != nil {
		args = append(args, auth.Username, auth.Password)
	}
	return ies7.NewEsClient(ctx, ec.Address, ec.Timeout, ec.Sniff, ec.Proxy, args...)
}

func NewDefaultDBEnv(ctx context.Context) (ul *DBUtil, err error) {
	ul = &DBUtil{
		dbs:    make(map[string]*gorm.DB),
		rawdbs: make(map[string]*sql.DB),
		ctx:    ctx}
	defer func() {
		if err != nil {
			ul.Stop()
		}
	}()
	select {
	case <-ctx.Done():
		ul.Stop()
	default:
		ul.dbs["musicDB"], err = InitDB(g.Config().Mysql)
		if err != nil {
			return
		}
		ul.rawdbs["rawDB"], err = InitRawDB(g.Config().Mysql)
		if err != nil {
			return
		}
		ul.mgoclient, err = InitMongo(g.Config().MongoDb.ModId, g.Config().MongoDb.CmdId,
			g.Config().MongoDb.DbName, g.Config().MongoDb.User, g.Config().MongoDb.Passwd)
		if err != nil {
			return
		}
		ul.esclient, err = InitElastic(ctx)
		if err != nil {
			return
		}
		ul.esclient7, err = InitElastic7(ctx)
		if err != nil {
			return
		}
	}
	return
}

func (ul *DBUtil) Start() (err error) {
	driver.CmsDriver, err = driver.NewCMSDriver(ul.ctx, ul.dbs["musicDB"], ul.dbs["importDB"],
		ul.dbs["KtrackDB"], ul.dbs["klyricDB"], ul.dbs["lyricDB"], ul.rawdbs["rawDB"], ul.mgoclient, ul.importclient)
	if err != nil {
		return
	}
	dblogic.TkDriver = dblogic.NewTracksDriver(driver.CmsDriver)
	dblogic.VoDriver = dblogic.NewVideosDriver(driver.CmsDriver)
	im.MgDriver = im.NewMongoDriver(driver.CmsDriver)
	ies.EsDriver = ul.esclient
	ies7.EsDriver = ul.esclient7
	go ies.EsDriver.Run()
	go ies7.EsDriver.Run()
	return nil
}

func (ul *DBUtil) Stop() (err error) {
	if driver.CmsDriver != nil {
		driver.CmsDriver.Close()
	}
	if ul.esclient != nil {
		ul.esclient.Close()
	}
	if ul.esclient7 != nil {
		ul.esclient7.Close()
	}
	return nil
}
