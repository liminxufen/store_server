package driver

import (
	"context"
	"database/sql"

	"github.com/jinzhu/gorm"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/store_server/logger"
	"sync"
)

//JOOX CMS DB驱动封装
var (
	CmsDriver *CMSDriver
)

//cms driver
type CMSDriver struct {
	MusicDB     *gorm.DB
	ImportDB    *gorm.DB
	KtrackDB    *gorm.DB
	KlyricDB    *gorm.DB
	LyricDB     *gorm.DB
	RawDB       *sql.DB
	MongoClient *mongo.Client
	ImportMongo *mongo.Client
	sync.RWMutex
	Ctx    context.Context
	Cancel func()
}

func NewCMSDriver(ct context.Context, musicDb, importDb, ktrackDb, klyricDb, lyricDb *gorm.DB,
	rawDb *sql.DB, client, importClient *mongo.Client) (cd *CMSDriver, err error) {
	cd = &CMSDriver{}
	cd.Lock()
	defer cd.Unlock()
	cd.MusicDB, cd.ImportDB, cd.KtrackDB, cd.KlyricDB = musicDb, importDb, ktrackDb, klyricDb
	cd.LyricDB, cd.RawDB, cd.MongoClient, cd.ImportMongo, cd.Ctx = lyricDb, rawDb, client, importClient, ct
	return
}

func (cd *CMSDriver) clone() *CMSDriver {
	cv := &CMSDriver{
		MusicDB: cd.MusicDB.Begin(),
		Ctx:     cd.Ctx,
	}
	return cv
}

func (cd *CMSDriver) Begin() (cv *CMSDriver, err error) { //开始事务
	cv = cd.clone()
	err = cv.MusicDB.Error
	logger.Entry().Debugf("CMSDriver[db: %p] begin transaction...", cv.MusicDB)
	return
}

func (cd *CMSDriver) Rollback() { //支持回滚
	logger.Entry().Debugf("CMSDriver[db: %p] rollback...", cd.MusicDB)
	if err := cd.MusicDB.Rollback().Error; err != nil {
		logger.Entry().Errorf("CMSDriver rollback err: %v", err)
	}
}

func (cd *CMSDriver) Commit() { //提交事务
	logger.Entry().Debugf("CMSDriver[db: %p] commit...", cd.MusicDB)
	if err := cd.MusicDB.Commit().Error; err != nil {
		logger.Entry().Errorf("CMSDriver commit err: %v", err)
		cd.Rollback()
	}
}

func (cd *CMSDriver) Close() {
	if cd == nil {
		return
	}
	if cd.MusicDB != nil {
		cd.MusicDB.Close()
	}
	if cd.ImportDB != nil {
		cd.ImportDB.Close()
	}
	if cd.MongoClient != nil {
		cd.MongoClient.Disconnect(cd.Ctx)
	}
	if cd.ImportMongo != nil {
		cd.ImportMongo.Disconnect(cd.Ctx)
	}
}
