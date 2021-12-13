package op

import (
	"github.com/store_server/dbtools/dataplatform"
	"github.com/store_server/logger"
	"github.com/store_server/store_server_http/g"
	"github.com/store_server/store_server_http/kits"
)

/************************ 数据平台搜索相关 ***************************/
//search track request
type SearchTrackDReq struct {
	Page        int64                  `json:"page,omitempty"`
	PageSize    int64                  `json:"pageSize,omitempty"`
	Sql         string                 `json:"sql"`
	QueryFields []string               `json:"fields,omitempty"`
	Conditions  map[string]interface{} `json:"conditions,omitempty"`
	OrderBy     [][2]string            `json:"order,omitempty"`
	Scopes      [][3]interface{}       `json:"scopes,omitempty"`
}

//search track response
type SearchTrackDRsp struct {
	Items interface{} `json:"items"`
	Total int64       `json:"total"`
}

func TrackSearchFromDataplatForm(req *SearchTrackDReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.TrackSearchFromDataplatForm", &err, logger.Entry())
	ret := &SearchTrackDRsp{}
	var data interface{}
	var total int64
	if len(req.Sql) != 0 {
		data, total, err = dataplatform.DpDriver.SearchBySql(req.Sql, req.Page, req.PageSize)
	} else {
		data, total, err = dataplatform.DpDriver.SearchByCondition(g.Config().Dataplatform.TrackType,
			req.QueryFields, req.Conditions, req.Page, req.PageSize, req.OrderBy, req.Scopes)
	}
	if err != nil {
		logger.Entry().Errorf("search track from dataplatform error: %v|request: %v", err, *req)
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	ret.Items = data
	ret.Total = total
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

//search album request
type SearchAlbumDReq struct {
	Page        int64                  `json:"page,omitempty"`
	PageSize    int64                  `json:"pageSize,omitempty"`
	Sql         string                 `json:"sql"`
	QueryFields []string               `json:"fields,omitempty"`
	Conditions  map[string]interface{} `json:"conditions,omitempty"`
	OrderBy     [][2]string            `json:"order,omitempty"`
	Scopes      [][3]interface{}       `json:"scopes,omitempty"`
}

//search album response
type SearchAlbumDRsp struct {
	Items interface{} `json:"items"`
	Total int64       `json:"total"`
}

func AlbumSearchFromDataplatForm(req *SearchAlbumDReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.AlbumSearchFromDataplatForm", &err, logger.Entry())
	ret := &SearchAlbumDRsp{}
	var data interface{}
	var total int64
	if len(req.Sql) != 0 {
		data, total, err = dataplatform.DpDriver.SearchBySql(req.Sql, req.Page, req.PageSize)
	} else {
		data, total, err = dataplatform.DpDriver.SearchByCondition(g.Config().Dataplatform.AlbumType,
			req.QueryFields, req.Conditions, req.Page, req.PageSize, req.OrderBy, req.Scopes)
	}
	if err != nil {
		logger.Entry().Errorf("search album from dataplatform error: %v|request: %v", err, *req)
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	ret.Items = data
	ret.Total = total
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

//search singer request
type SearchSingerDReq struct {
	Page        int64                  `json:"page,omitempty"`
	PageSize    int64                  `json:"pageSize,omitempty"`
	Sql         string                 `json:"sql"`
	QueryFields []string               `json:"fields,omitempty"`
	Conditions  map[string]interface{} `json:"conditions,omitempty"`
	OrderBy     [][2]string            `json:"order,omitempty"`
	Scopes      [][3]interface{}       `json:"scopes,omitempty"`
}

//search singer response
type SearchSingerDRsp struct {
	Items interface{} `json:"items"`
	Total int64       `json:"total"`
}

func SingerSearchFromDataplatForm(req *SearchSingerDReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.SingerSearchFromDataplatForm", &err, logger.Entry())
	ret := &SearchAlbumDRsp{}
	var data interface{}
	var total int64
	if len(req.Sql) != 0 {
		data, total, err = dataplatform.DpDriver.SearchBySql(req.Sql, req.Page, req.PageSize)
	} else {
		data, total, err = dataplatform.DpDriver.SearchByCondition(g.Config().Dataplatform.SingerType,
			req.QueryFields, req.Conditions, req.Page, req.PageSize, req.OrderBy, req.Scopes)
	}
	if err != nil {
		logger.Entry().Errorf("search singer from dataplatform error: %v|request: %v", err, *req)
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	ret.Items = data
	ret.Total = total
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}
