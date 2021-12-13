package mongo

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/store_server/dbtools/driver"
	"github.com/store_server/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"go.mongodb.org/mongo-driver/mongo"

	m "github.com/store_server/dbtools/models"
)

//mongo driver
type MongoDriver struct {
	*driver.CMSDriver

	lock       sync.RWMutex
	collection string
}

func NewMongoDriver(cmsDriver *driver.CMSDriver) *MongoDriver {
	return &MongoDriver{cmsDriver, sync.RWMutex{}, ""}
}

var (
	MgDriver *MongoDriver
)

/************************ 通用方法 ************************/
func (md *MongoDriver) FindLastDoc(db, col string) (bson.Raw, error) {
	opts := options.FindOne()
	opts.SetSort(bson.D{{"_id", -1}})
	collection := md.MongoClient.Database(db).Collection(col)
	res := collection.FindOne(md.Ctx, bson.M{}, opts)
	if res.Err() != nil {
		return nil, res.Err()
	}
	return res.DecodeBytes()
}

func (md *MongoDriver) FindOneById(db, col string, id interface{}) (bson.Raw, error) {
	filter := bson.M{"_id": id}
	collection := md.MongoClient.Database(db).Collection(col)
	res := collection.FindOne(md.Ctx, filter)
	if res.Err() != nil {
		return nil, res.Err()
	}
	return res.DecodeBytes()
}

func (md *MongoDriver) FindOneByFilter(db, col string, filter interface{}) (bson.Raw, error) {
	collection := md.MongoClient.Database(db).Collection(col)
	res := collection.FindOne(md.Ctx, filter)
	if res.Err() != nil {
		return nil, res.Err()
	}
	return res.DecodeBytes()
}

func (md *MongoDriver) FindImportManyByFilter(db, col string, filter interface{},
	opt *options.FindOptions, results interface{}) error {
	if md.ImportMongo == nil {
		return fmt.Errorf("invalid import mongo client")
	}
	collection := md.ImportMongo.Database(db).Collection(col)
	cur, err := collection.Find(md.Ctx, filter, opt)
	if err != nil {
		return err
	}
	return cur.All(md.Ctx, results)
}

func (md *MongoDriver) FindOneImportByFilter(db, col string, filter interface{},
	opts ...*options.FindOneOptions) (bson.Raw, error) {
	if md.ImportMongo == nil {
		return nil, fmt.Errorf("invalid import mongo client")
	}
	collection := md.ImportMongo.Database(db).Collection(col)
	res := collection.FindOne(md.Ctx, filter, opts...)
	if res.Err() != nil {
		return nil, res.Err()
	}
	return res.DecodeBytes()
}

func (md *MongoDriver) FindManyByFilter(db, col string, filter interface{},
	opt *options.FindOptions, results interface{}) error {
	collection := md.MongoClient.Database(db).Collection(col)
	cur, err := collection.Find(md.Ctx, filter, opt)
	if err != nil {
		return err
	}
	return cur.All(md.Ctx, results)
}

func (md *MongoDriver) InsertOneDoc(db, col string, doc interface{}) error {
	collection := md.MongoClient.Database(db).Collection(col)
	_, err := collection.InsertOne(md.Ctx, doc)
	if err != nil {
		return err
	}
	return nil
}

func (md *MongoDriver) InsertManyDoc(db, col string, docs []interface{}) error {
	collection := md.MongoClient.Database(db).Collection(col)
	_, err := collection.InsertMany(md.Ctx, docs)
	if err != nil {
		return err
	}
	return nil
}

func (md *MongoDriver) UpdateOneByID(db, col string, id interface{}, update interface{}) error {
	filter := bson.M{"_id": id}
	collection := md.MongoClient.Database(db).Collection(col)
	_, err := collection.UpdateOne(md.Ctx, filter, update)
	if err != nil {
		return err
	}
	return err
}

func (md *MongoDriver) UpdateOneByFilter(db, col string, filter interface{}, update interface{}) error {
	collection := md.MongoClient.Database(db).Collection(col)
	//TODO query for sharded findAndModify must have shardkey
	res := collection.FindOneAndUpdate(md.Ctx, filter, update)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

func (md *MongoDriver) UpdateManyByFilter(db, col string, filter interface{}, update interface{}) error {
	collection := md.MongoClient.Database(db).Collection(col)
	_, err := collection.UpdateMany(md.Ctx, filter, update)
	if err != nil {
		return err
	}
	return err
}

func (md *MongoDriver) DeleteOneByID(db, col string, id interface{}) error {
	filter := bson.M{"_id": id}
	collection := md.MongoClient.Database(db).Collection(col)
	_, err := collection.DeleteOne(md.Ctx, filter)
	if err != nil {
		return err
	}
	return err
}

func (md *MongoDriver) DeleteOneByFilter(db, col string, filter interface{}) error {
	collection := md.MongoClient.Database(db).Collection(col)
	_, err := collection.DeleteOne(md.Ctx, filter)
	if err != nil {
		return err
	}
	return err
}

func (md *MongoDriver) DeleteManyByFilter(db, col string, filter interface{}) error {
	collection := md.MongoClient.Database(db).Collection(col)
	_, err := collection.DeleteMany(md.Ctx, filter)
	if err != nil {
		return err
	}
	return err
}

/*********************** 封装 **********************/
func (md *MongoDriver) GetDocById(db, col string, id, doc interface{}) error {
	raw, err := md.FindOneById(db, col, id)
	if err != nil {
		return err
	}
	docVal := reflect.ValueOf(doc)
	if docVal.Kind() != reflect.Ptr {
		return bson.Unmarshal(raw, &doc)
	}
	return bson.Unmarshal(raw, doc)
}

func (md *MongoDriver) GetDocByFilter(db, col string, filter, doc interface{}) error {
	raw, err := md.FindOneByFilter(db, col, filter)
	if err != nil {
		return err
	}
	docVal := reflect.ValueOf(doc)
	if docVal.Kind() != reflect.Ptr {
		return bson.Unmarshal(raw, &doc)
	}
	return bson.Unmarshal(raw, doc)
}

func (md *MongoDriver) GetImportDocByFilter(db, col string, filter, doc interface{},
	opts ...*options.FindOneOptions) error {
	raw, err := md.FindOneImportByFilter(db, col, filter, opts...)
	if err != nil {
		return err
	}
	docVal := reflect.ValueOf(doc)
	if docVal.Kind() != reflect.Ptr {
		return bson.Unmarshal(raw, &doc)
	}
	return bson.Unmarshal(raw, doc)
}

func (md *MongoDriver) GetDocsByFilter(db, col string, filter, docs interface{},
	args ...int64) error {
	var page, pagesize int64
	if len(args) >= 2 {
		page, pagesize = args[0], args[1]
		if page == 0 {
			page = 1
		}
	}
	opt := &options.FindOptions{}
	if page > 0 && pagesize > 0 {
		offset := (page - 1) * pagesize
		opt = opt.SetSort(bson.D{{"_id", -1}}).SetSkip(offset).SetLimit(pagesize)
	}
	docsVal := reflect.ValueOf(docs)
	if docsVal.Kind() != reflect.Ptr || docsVal.Kind() == reflect.Slice {
		return md.FindManyByFilter(db, col, filter, opt, &docs)
	}
	return md.FindManyByFilter(db, col, filter, opt, docs)
}

func (md *MongoDriver) FindLastDocId(db, col string, mod interface{}) (interface{}, error) {
	lastDoc, err := md.FindLastDoc(db, col)
	if err != nil {
		return -1, err
	}
	modVal := reflect.ValueOf(mod)
	if modVal.Kind() != reflect.Ptr {
		err = bson.Unmarshal(lastDoc, &mod)
	} else {
		err = bson.Unmarshal(lastDoc, mod)
	}
	if err != nil {
		return -1, err
	}
	var ok1, ok2, isPtr bool
	modType := reflect.TypeOf(mod)
	if modType.Kind() != reflect.Ptr {
		_, ok1 = modType.FieldByName("Id")
		_, ok2 = modType.FieldByName("ID")
	} else {
		isPtr = true
		modElmType := modType.Elem()
		_, ok1 = modElmType.FieldByName("Id")
		_, ok2 = modElmType.FieldByName("ID")
	}
	if isPtr {
		modVal = modVal.Elem()
	}
	if ok1 || ok2 {
		var idValue reflect.Value
		if ok1 {
			idValue = modVal.FieldByName("Id")
		}
		if ok2 {
			idValue = modVal.FieldByName("ID")
		}
		return idValue.Interface(), nil
	}
	return -1, fmt.Errorf("get last doc id error by model type")
}

/************************ track related ************************/
func (md *MongoDriver) GetTrack() error {

	return nil
}

/************************ auto_publish_album ************************/
func (md *MongoDriver) GetAutoPublishAlbumById(id interface{}) (*m.PublishedAlbum, error) {
	pa := &m.PublishedAlbum{}
	err := md.GetDocById("music_cms", "auto_publish_album", id, pa)
	return pa, err
}

func (md *MongoDriver) GetAutoPublishAlbum(filter interface{}) (*m.PublishedAlbum, error) {
	pa := &m.PublishedAlbum{}
	err := md.GetDocByFilter("music_cms", "auto_publish_album", filter, pa)
	return pa, err
}

func (md *MongoDriver) GetManyPublishAlbum(filter interface{},
	page, pagesize int64) ([]*m.PublishedAlbum, error) {
	pas := []*m.PublishedAlbum{}
	err := md.GetDocsByFilter("music_cms", "auto_publish_album", filter, &pas, page, pagesize)
	return pas, err
}

func (md *MongoDriver) InsertPublishAlbum(phAlbum *m.PublishedAlbum) (int64, error) {
	current := time.Now()
	phAlbum.CreateTime, phAlbum.ModifyTime = current, current
	phAlbum.Deleted = 0
	lastEr := &m.PublishedAlbum{}
	lastId, err := md.FindLastDocId("music_cms", "auto_publish_album", lastEr)
	if err != nil {
		logger.Entry().Errorf("find last doc id error: %v", err)
	}
	id := lastId.(int64)
	phAlbum.Id = id + 1
	err = md.InsertOneDoc("music_cms", "auto_publish_album", phAlbum)
	if err != nil {
		return -1, err
	}
	return phAlbum.Id, nil
}

func (md *MongoDriver) UpdatePublishAlbumById(id interface{}, update bson.M) error {
	update["modify_time"] = time.Now()
	update = bson.M{"$set": update}
	return md.UpdateOneByID("music_cms", "auto_publish_album", id, update)
}

func (md *MongoDriver) UpdatePublishAlbum(filter interface{}, update bson.M) error {
	update["modify_time"] = time.Now()
	update = bson.M{"$set": update}
	return md.UpdateOneByFilter("music_cms", "auto_publish_album", filter, update)
}

func (md *MongoDriver) UpdateManyPublishAlbum(filter interface{}, update bson.M) error {
	update["modify_time"] = time.Now()
	update = bson.M{"$set": update}
	return md.UpdateManyByFilter("music_cms", "auto_publish_album", filter, update)
}

func (md *MongoDriver) DeletePublishAlbumById(id interface{}) error {
	return md.DeleteOneByID("music_cms", "auto_publish_album", id)
}

func (md *MongoDriver) DeletePublishAlbum(filter interface{}) error {
	return md.DeleteOneByFilter("music_cms", "auto_publish_album", filter)
}

func (md *MongoDriver) DeleteManyPublishAlbum(filter interface{}) error {
	return md.DeleteManyByFilter("music_cms", "auto_publish_album", filter)
}

/************************ external_resources ************************/
func (md *MongoDriver) GetExternalResourcesById(id interface{}) (*m.ExternalResource, error) {
	rs := &m.ExternalResource{}
	err := md.GetDocById("music_cms", "external_resources", id, rs)
	return rs, err
}

func (md *MongoDriver) GetExternalResources(filter interface{}) (*m.ExternalResource, error) {
	rs := &m.ExternalResource{}
	err := md.GetDocByFilter("music_cms", "external_resources", filter, rs)
	return rs, err
}

func (md *MongoDriver) GetManyExternalResources(filter interface{},
	page, pagesize int64) ([]*m.ExternalResource, error) {
	ers := []*m.ExternalResource{}
	err := md.GetDocsByFilter("music_cms", "external_resources", filter, &ers, page, pagesize)
	return ers, err
}

func (md *MongoDriver) InsertExternalResources(extResource *m.ExternalResource) (int64, error) {
	current := time.Now()
	extResource.CreateTime, extResource.ModifyTime = current, current
	extResource.Deleted = 0
	lastEr := &m.ExternalResource{}
	lastId, err := md.FindLastDocId("music_cms", "external_resources", lastEr)
	if err != nil {
		logger.Entry().Errorf("find last doc id error: %v", err)
	}
	id := lastId.(int64)
	extResource.Id = id + 1
	err = md.InsertOneDoc("music_cms", "external_resources", extResource)
	if err != nil {
		return -1, err
	}
	return extResource.Id, nil
}

func (md *MongoDriver) UpdateExternalResourcesById(id interface{}, update bson.M) error {
	update["modify_time"] = time.Now()
	update = bson.M{"$set": update}
	return md.UpdateOneByID("music_cms", "external_resources", id, update)
}

func (md *MongoDriver) UpdateExternalResources(filter interface{}, update bson.M) error {
	update["modify_time"] = time.Now()
	update = bson.M{"$set": update}
	return md.UpdateOneByFilter("music_cms", "external_resources", filter, update)
}

func (md *MongoDriver) UpdateManyExternalResources(filter interface{}, update bson.M) error {
	update["modify_time"] = time.Now()
	update = bson.M{"$set": update}
	return md.UpdateManyByFilter("music_cms", "external_resources", filter, update)
}

func (md *MongoDriver) DeleteExternalResourcesById(id interface{}) error {
	return md.DeleteOneByID("music_cms", "external_resources", id)
}

func (md *MongoDriver) DeleteExternalResources(filter interface{}) error {
	return md.DeleteOneByFilter("music_cms", "external_resources", filter)
}

func (md *MongoDriver) DeleteManyExternalResources(filter interface{}) error {
	return md.DeleteManyByFilter("music_cms", "external_resources", filter)
}
