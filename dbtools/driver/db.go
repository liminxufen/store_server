package driver

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //初始化mysql driver
	"github.com/store_server/dbtools/models"
	"github.com/store_server/logger"
	log "github.com/store_server/logger"
	"github.com/store_server/utils/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"sync"
	"time"
)

var getEnvOrDefault = common.GetEnvOrDefault

var (
	mysql string = fmt.Sprintf("%s:%s@tcp(%s:%s)/music?charset=utf8&parseTime=True&loc=Local",
		getEnvOrDefault("MYSQL_USER", "root"), getEnvOrDefault("MYSQL_PASSWORD", ""),
		getEnvOrDefault("MYSQL_HOST", "127.0.0.1"), getEnvOrDefault("MYSQL_PORT", "3306"))

	mongo_conf string = fmt.Sprintf("mongodb://%s:%s@%s:%s?authMechanism=MONGODB-CR&authSource=music_cms",
		getEnvOrDefault("MONGO_USER", "music_cms"), getEnvOrDefault("MONGO_PASSWORD", ""),
		getEnvOrDefault("MONGO_HOST", "localhost"), getEnvOrDefault("MONGO_PORT", "27017"))

	mongo_db_name = getEnvOrDefault("MONGO_DBNAME", "music_cms")
)

var Tables = []interface{}{
	&models.Track{}, &models.TrackExtraOs{},
}

func CreateDB(cf string) (db *gorm.DB, err error) {
	if len(cf) == 0 {
		err = fmt.Errorf("mysql config is empty.")
		return
	}
	db, err = gorm.Open("mysql", cf)
	if err != nil {
		return
	}
	db.SingularTable(true) //全局禁用表名复数
	db.DB().SetConnMaxLifetime(time.Minute * 10)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(200)
	for _, _ = range Tables { //no permission for product db
		//db.AutoMigrate(item)
	}
	return
}

func CreateRawDB(cf string) (db *sql.DB, err error) {
	if len(cf) == 0 {
		err = fmt.Errorf("mysql config is empty.")
		return
	}
	db, err = sql.Open("mysql", cf)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(200)
	return db, nil
}

func CreateMongo(opts *options.ClientOptions) (client *mongo.Client, err error) {
	client, err = mongo.NewClient(opts)
	if err != nil {
		logger.Entry().Errorf("new mongo client config opts: %v|error: %v", *opts, err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return
	}
	return client, nil
}

/*封装测试MYSQL DB环境*/
type Scheme struct {
	localDB *gorm.DB
	Mysql   string
	sync.RWMutex
}

var DefaultTestScheme = &Scheme{Mysql: mysql}

func (s *Scheme) DB() *gorm.DB {
	s.createDB()
	return s.localDB
}

func (s *Scheme) createDB() {
	s.Lock()
	defer s.Unlock()
	if s.localDB != nil {
		return
	}

	db, err := CreateDB(s.Mysql)
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range Tables { //删除旧表
		if db.HasTable(item) {
			if err := db.DropTable(item).Error; err != nil {
				log.Fatal(err)
			}
		}
	}
	db = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8")
	for _, item := range Tables { //根据model自动更新表结构
		if err := db.AutoMigrate(item).Error; err != nil {
			log.Fatal(err)
		}
	}
	s.localDB = db
	if os.Getenv("not_sql") != "true" {
		s.localDB = s.localDB.LogMode(true)
	}
}

func (s *Scheme) Setup() {
	s.createDB()
}

func (s *Scheme) Cleanup() {
	/*db := s.localDB
	for _, item := range Tables {
		if err := db.Delete(item).Error; err != nil {
			log.Fatal(err)
		}
	}*/
}

var Collections = []string{"auto_publish_album_test", "singer_info_test"}

/*封装测试MONGO环境*/
type MongoScheme struct {
	localMDB   *mongo.Client
	dbs        *mongo.Database
	collection *mongo.Collection
	Mongo      string
	opts       *options.ClientOptions
	sync.RWMutex
	Ctx context.Context
}

var DefaultTestMongoScheme = &MongoScheme{opts: mockClientopts()}

func mockClientopts() *options.ClientOptions {
	//mock mongo client options
	dur := time.Duration(3) * time.Second
	poolSize := uint64(100)
	direct := false
	opts := &options.ClientOptions{
		Hosts: []string{},
		Auth: &options.Credential{
			//AuthMechanism: "MONGODB-CR",
			AuthMechanism: "SCRAM-SHA-1",
			AuthSource:    "music_cms",
			Username:      "music_cms",
			Password:      getEnvOrDefault("MONGO_TEST_PASSWD", ""),
		},
		MaxPoolSize:    &poolSize,
		ConnectTimeout: &dur,
		Direct:         &direct,
	}
	return opts
}

func (s *MongoScheme) DB() *mongo.Client {
	s.createClient()
	return s.localMDB
}

func (s *MongoScheme) Collection(col string) *mongo.Collection {
	s.collection = s.localMDB.Database(mongo_db_name).Collection(col)
	return s.collection
}

func (s *MongoScheme) Database() *mongo.Database {
	s.dbs = s.localMDB.Database(mongo_db_name)
	return s.dbs
}

func (s *MongoScheme) createClient() {
	s.Lock()
	defer s.Unlock()
	if s.localMDB != nil {
		return
	}
	s.Ctx = context.Background()
	mdb, err := CreateMongo(s.opts)
	if err != nil {
		fmt.Printf("----------------- create mongo db client error: %v --------------------\n", err)
		log.Fatal(err)
	}
	//删除旧文档集合
	for _, col := range Collections {
		collect := mdb.Database(mongo_db_name).Collection(col)
		collect.DeleteMany(s.Ctx, bson.M{})
	}
	//创建新文档集合
	for _, col := range Collections {
		mdb.Database(mongo_db_name).Collection(col)
	}
	s.localMDB = mdb
}

func (s *MongoScheme) Setup() {
	s.createClient()
}

func (s *MongoScheme) Cleanup() {
	/*mdb := s.localMDB
	for _, col := range Collections {
		collect := mdb.Database(g.Config().MongoDb.DbName).Collection(col)
		collect.Drop(s.Ctx)
	}*/
}

//从环境变量获取mysql和mongo配置
func init() {
	if os.Getenv("mysql_env") != "" {
		mysql = os.Getenv("mysql_env")
	}
	if os.Getenv("mongo_env") != "" {
		mongo_conf = os.Getenv("mongo_env")
	}
}
