package op

import (
	m "github.com/store_server/dbtools/models"
	"github.com/store_server/dbtools/mongo"
	"github.com/store_server/logger"
	"github.com/store_server/store_server_http/kits"
)

/************************ ExternalResources查询相关 ***************************/
//query external resource request
type QueryExternalResourcesReq struct {
	Id       int64                  `json:"id,omitempty"`
	Page     int64                  `json:"page,omitempty"`
	PageSize int64                  `json:"pageSize,omitempty"`
	Filter   map[string]interface{} `json:"filter,omitempty"`
}

//query external resource response
type QueryExternalResourcesRsp struct {
	ExtResources interface{} `json:"external_resources"`
}

func ExternalResourcesQuery(req *QueryExternalResourcesReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.ExternalResourcesQuery", &err, logger.Entry())
	ret := QueryExternalResourcesRsp{}
	var ers interface{}
	if req.Id != 0 {
		ers, err = mongo.MgDriver.GetExternalResourcesById(req.Id)
	} else { //others query condition
		if len(req.Filter) == 0 && req.Page == 0 && req.PageSize == 0 {
			logger.Entry().Errorf("query external resources filter conditions is invalid")
			rsp = kits.APIWrapRsp(kits.ErrOther, "query external resources filter conditions is invalid", ret)
			return
		}
		ers, err = mongo.MgDriver.GetManyExternalResources(req.Filter, req.Page, req.PageSize)
	}
	if err != nil {
		logger.Entry().Errorf("query external resources error: %v|request: %v", err, *req)
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	ret.ExtResources = ers
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

/************************ ExternalResources更新相关 ***************************/
//update external resource request
type UpdateExternalResourcesReq struct {
	Id     int64                  `json:"id"`
	Conds  map[string]interface{} `json:"condition,omitempty"`
	Fields map[string]interface{} `json:"updateDoc,omitempty"`
}

//update external resource response
type UpdateExternalResourcesRsp struct {
	Affected int64 `json:"affected,omitempty"`
}

func ExternalResourcesUpdate(req *UpdateExternalResourcesReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.ExternalResourcesUpdate", &err, logger.Entry())
	ret := UpdateExternalResourcesRsp{}
	if req.Id != 0 {
		err = mongo.MgDriver.UpdateExternalResourcesById(req.Id, req.Fields)
	} else {
		err = mongo.MgDriver.UpdateExternalResources(req.Conds, req.Fields)
	}
	if err != nil {
		logger.Entry().Errorf("update external resources error: %v|request: %v", err, *req)
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

/************************ ExternalResources创建相关 ***************************/
//insert external resource request
type InsertExternalResourcesReq struct {
	ExternalResource *m.ExternalResource `json:"externalResource"`
}

//insert external resource response
type InsertExternalResourcesRsp struct {
	Id int64 `json:"id"`
}

func ExternalResourcesInsert(req *InsertExternalResourcesReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.ExternalResourcesInsert", &err, logger.Entry())
	ret := InsertExternalResourcesRsp{-1}
	var id int64
	id, err = mongo.MgDriver.InsertExternalResources(req.ExternalResource)
	if err != nil {
		logger.Entry().Errorf("insert external resources error: %v|request: %v", err, *req)
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	ret.Id = id
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

/************************ ExternalResources删除相关 ***************************/
//delete external resource request
type DeleteExternalResourcesReq struct {
	Id    int64                  `json:"id"`
	Conds map[string]interface{} `json:"condition,omitempty"`
}

//delete external resource response
type DeleteExternalResourcesRsp struct {
	Affected int64 `json:"affected,omitempty"`
}

func ExternalResourcesDelete(req *DeleteExternalResourcesReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.ExternalResourcesDelete", &err, logger.Entry())
	ret := DeleteExternalResourcesRsp{}
	if req.Id != 0 {
		err = mongo.MgDriver.DeleteExternalResourcesById(req.Id)
	} else {
		err = mongo.MgDriver.DeleteExternalResources(req.Conds)
	}
	if err != nil {
		logger.Entry().Errorf("delete external resources error: %v|request: %v", err, *req)
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}
