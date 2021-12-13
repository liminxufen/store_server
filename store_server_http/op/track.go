package op

import (
	"github.com/store_server/dbtools/dblogic"
	m "github.com/store_server/dbtools/models"
	"github.com/store_server/logger"
	"github.com/store_server/store_server_http/kits"
)

/************************ 歌曲查询相关 ***************************/
//query track request
type QueryTrackReq struct {
	RawSql   string                 `json:"rawSql"`
	Ids      []int64                `json:"ids"`
	Page     int64                  `json:"page,omitempty"`
	PageSize int64                  `json:"pageSize,omitempty"`
	Fields   map[string]interface{} `json:"fields,omitempty"`
}

//query track response
type QueryTrackRsp struct {
	Tracks []*m.Track `json:"tracks"`
	Total  int64      `json:"total,omitempty"`
}

func TracksQuery(req *QueryTrackReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.TracksQuery", &err, logger.Entry())
	ret := QueryTrackRsp{}
	var tracks []*m.Track
	if len(req.RawSql) != 0 {
		tracks, ret.Total, err = dblogic.TkDriver.ExecRawQuerySql4Track(req.RawSql, req.Page, req.PageSize)
	} else if len(req.Ids) != 0 && req.Ids[0] != 0 {
		tracks, ret.Total, err = dblogic.TkDriver.GetTracksByIds(req.Ids)
	} else { //others query condition
		if len(req.Fields) == 0 && req.Page == 0 && req.PageSize == 0 {
			logger.Entry().Errorf("query tracks fields conditions is nil")
			rsp = kits.APIWrapRsp(kits.ErrOther, "query tracks fields conditions is invalid", ret)
			return
		}
		tracks, ret.Total, err = dblogic.TkDriver.GetTracksByCondition(req.Fields, req.Page, req.PageSize)
	}
	if err != nil {
		logger.Entry().Errorf("query tracks error: %v|request: %v", err, *req)
		ret.Total = 0
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	ret.Tracks = tracks
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

//query track extra os request
type QueryTrackExtraOsReq struct {
	RawSql   string                 `json:"rawSql"`
	Id       int64                  `json:"id"`
	Region   int64                  `json:"region"`
	Page     int64                  `json:"page,omitempty"`
	PageSize int64                  `json:"pageSize,omitempty"`
	Fields   map[string]interface{} `json:"fields,omitempty"`
}

//query track extra os response
type QueryTrackExtraOsRsp struct {
	Tracks []*m.TrackExtraOs `json:"trackExtraOs"`
	Total  int64             `json:"total,omitempty"`
}

func TrackExtraOsQuery(req *QueryTrackExtraOsReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.TrackExtraOsQuery", &err, logger.Entry())
	ret := QueryTrackExtraOsRsp{}
	var tracks []*m.TrackExtraOs
	if len(req.RawSql) != 0 {
		tracks, ret.Total, err = dblogic.TkDriver.ExecRawQuerySql4TrackExtraOs(req.RawSql, req.Page, req.PageSize)
	} else if req.Id != 0 {
		var track *m.TrackExtraOs
		track, err = dblogic.TkDriver.GetOneTrackExtraOs(req.Id, req.Region)
		tracks = []*m.TrackExtraOs{track}
		ret.Total = 1
	} else {
		if len(req.Fields) == 0 && req.Page == 0 && req.PageSize == 0 {
			logger.Entry().Errorf("query track extra os fields conditions is nil")
			rsp = kits.APIWrapRsp(kits.ErrOther, "query track extra os fields conditions is invalid", ret)
			return
		}
		tracks, ret.Total, err = dblogic.TkDriver.GetTrackExtraOsByCondition(req.Fields, req.Page, req.PageSize)
	}
	if err != nil {
		logger.Entry().Errorf("query track extra os error: %v|request: %v", err, *req)
		ret.Total = 0
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	ret.Tracks = tracks
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

/************************ 歌曲更新相关 ***************************/
//update track request
type UpdateTrackReq struct {
	Ids    []int64                `json:"ids"`
	Conds  map[string]interface{} `json:"conditions,omitempty"`
	Fields map[string]interface{} `json:"updateFields,omitempty"`
}

//update track response
type UpdateTrackRsp struct {
	Affected int64 `json:"affected,omitempty"`
}

func TracksUpdate(req *UpdateTrackReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.TracksUpdate", &err, logger.Entry())
	ret := UpdateTrackRsp{}
	ret.Affected, err = dblogic.TkDriver.UpdateTracksAttr(req.Ids, req.Conds, req.Fields)
	if err != nil {
		logger.Entry().Errorf("update tracks error: %v|request: %v", err, *req)
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

//update track extra os request
type UpdateTrackExtraOsReq struct {
	Ids    []int64                `json:"ids"`
	Conds  map[string]interface{} `json:"conditions,omitempty"`
	Fields map[string]interface{} `json:"updateFields,omitempty"`
}

//update track extra os response
type UpdateTrackExtraOsRsp struct {
	Affected int64 `json:"affected,omitempty"`
}

func TrackExtraOsUpdate(req *UpdateTrackExtraOsReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.TrackExtraOsUpdate", &err, logger.Entry())
	ret := UpdateTrackExtraOsRsp{}
	ret.Affected, err = dblogic.TkDriver.UpdateTrackExtraOsAttr(req.Ids, req.Conds, req.Fields)
	if err != nil {
		logger.Entry().Errorf("update track extra os error: %v|request: %v", err, *req)
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

/************************ 歌曲删除相关 ***************************/
//delete track request
type DeleteTrackReq struct {
	Ids   []int64                `json:"ids"`
	Conds map[string]interface{} `json:"conditions,omitempty"`
}

//delete track response
type DeleteTrackRsp struct {
	Affected int64 `json:"affected,omitempty"`
}

func TracksDelete(req *DeleteTrackReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.TracksDelete", &err, logger.Entry())
	ret := DeleteTrackRsp{}
	ret.Affected, err = dblogic.TkDriver.DeleteTracks(req.Ids, req.Conds)
	if err != nil {
		logger.Entry().Errorf("delete tracks error: %v|request: %v", err, *req)
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

//delete track extra os request
type DeleteTrackExtraOsReq struct {
	Ids   []int64                `json:"ids"`
	Conds map[string]interface{} `json:"conditions,omitempty"`
}

//delete track extra os response
type DeleteTrackExtraOsRsp struct {
	Affected int64 `json:"affected,omitempty"`
}

func TrackExtraOsDelete(req *DeleteTrackExtraOsReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.TrackExtraOsDelete", &err, logger.Entry())
	ret := DeleteTrackExtraOsRsp{}
	ret.Affected, err = dblogic.TkDriver.DeleteTrackExtraOs(req.Ids, req.Conds)
	if err != nil {
		logger.Entry().Errorf("delete track extra os error: %v|request: %v", err, *req)
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

/************************ 歌曲插入相关 ***************************/
//insert track request
type InsertTrackReq struct {
	Tracks       []*m.Track        `json:"tracks,omitempty"`
	TrackExtraOs []*m.TrackExtraOs `json:"trackExtraOs,omitempty"`
}

//insert track response
type InsertTrackRsp struct {
	Affected int64 `json:"affected,omitempty"`
}

func TracksInsert(req *InsertTrackReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.TracksInsert", &err, logger.Entry())
	ret := InsertTrackRsp{}
	var table string
	if len(req.Tracks) != 0 {
		ret.Affected, err = dblogic.TkDriver.InsertTracks(req.Tracks)
		table = "t_track"
	} else if len(req.TrackExtraOs) != 0 {
		ret.Affected, err = dblogic.TkDriver.InsertTrackExtraOs(req.TrackExtraOs)
		table = "t_track_extra_os"
	} else {
		logger.Entry().Errorf("insert tracks error: %v|table: %s|request: %v", err, table, *req)
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	if err != nil {
		logger.Entry().Errorf("invalid insert params")
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

/************************ 歌曲联合查询相关 ***************************/
//join query track request
type JoinQueryTrackReq struct {
	RawSql   string `json:"rawSql"`
	Page     int64  `json:"page,omitempty"`
	PageSize int64  `json:"pageSize,omitempty"`
}

//join query track response
type JoinQueryTrackRsp struct {
	Results [][]interface{} `json:"results"`
	Total   int64           `json:"total,omitempty"`
}

func TracksJoinQuery(req *JoinQueryTrackReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.TracksJoinQuery", &err, logger.Entry())
	ret := JoinQueryTrackRsp{}
	results, err := dblogic.TkDriver.JoinQueryWithRawSql(req.RawSql, req.Page, req.PageSize)
	if err != nil {
		logger.Entry().Errorf("join query tracks error: %v|request: %v", err, *req)
		ret.Total = 0
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
	}
	ret.Results, ret.Total = results, int64(len(results))
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

/************************ 查询歌曲是否能播放等信息 ***************************/
//query track play info request
type QueryTrackPlayReq struct {
	Id int64 `json:"id"`
}

//query track play info response
type QueryTrackPlayRsp struct {
	CanPlay bool `json:"can_play"`
}

func TrackPlayQuery(req *QueryTrackPlayReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.TrackPlayQuery", &err, logger.Entry())
	ret := QueryTrackPlayRsp{}
	if req.Id == 0 {
		logger.Entry().Errorf("query track play info conditions is nil")
		rsp = kits.APIWrapRsp(kits.ErrOther, "query track play info conditions is invalid", ret)
		return
	}
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

/************************ 查询歌曲是否有关联艺人信息 ***************************/
//query track singer info request
type QueryTrackSingerReq struct {
	Id     int64   `json:"id"`
	Region int64   `json:"region"`
	Ids    []int64 `json:"ids,omitempty"`
}

//query track singer info response
type QueryTrackSingerRsp struct {
	InvalidIds []int64 `json:"invalid_ids"`
	NoSinger   bool    `json:"no_singer"`
}

func TrackSingerQuery(req *QueryTrackSingerReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.TrackSingerQuery", &err, logger.Entry())
	ret := QueryTrackSingerRsp{}
	if req.Id == 0 && req.Region == 0 && len(req.Ids) == 0 {
		logger.Entry().Errorf("query track singer info conditions is nil")
		rsp = kits.APIWrapRsp(kits.ErrOther, "query track singer info conditions is invalid", ret)
		return
	}
	var hasSinger bool
	if req.Id != 0 {
	} else if len(req.Ids) != 0 {
	}
	if err != nil {
		logger.Entry().Errorf("query track singer info error: %v|request: %v", err, *req)
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	ret.NoSinger = !hasSinger
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

/************************ 查询歌曲播放链接 ***************************/
//query track play url request
type QueryTrackPlayURLReq struct {
	Id     int64 `json:"id"`
	Region int64 `json:"region"`
}

//query track play url response
type QueryTrackPlayURLRsp struct {
	URL string `json:"url"`
}

func TrackPlayURLQuery(req *QueryTrackPlayURLReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.TrackPlayURLQuery", &err, logger.Entry())
	ret := QueryTrackPlayURLRsp{}
	if req.Id == 0 {
		logger.Entry().Errorf("query track play url conditions is nil")
		rsp = kits.APIWrapRsp(kits.ErrOther, "query track play url conditions is invalid", ret)
		return
	}
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}
