package rpcServer

import (
	//"go.mongodb.org/mongo-driver/bson"
	"github.com/store_server/dbtools/dblogic"
	"github.com/store_server/dbtools/models"
	"github.com/store_server/logger"
	lm "github.com/store_server/store_server_rpc/rpc/common"
	"github.com/store_server/utils/common"
	"net/http"
	"time"
)

//track rpc service
type TrackService struct{}

func (s *TrackService) CreateTrack(hr *http.Request, req *lm.CreateTrackRpcReq, rsp *lm.CommRpcRsp) (err error) {
	defer common.TimeCostTrack(time.Now(), "TrackService rpc", "CreateTrack", err)
	payload := &lm.CreateTrackRpcRsp{}
	if err := lm.CheckParamsIsNil(req); err != nil {
		lm.WrapRpcRsp(2, "", payload, rsp)
		return err
	}
	lastId, e := dblogic.TkDriver.InsertOneTrack(req.CreateDoc)
	if e != nil {
		err = e
		logger.Entry().Errorf("rpc to create track error: %v|%v", *req, err)
		payload.Id = -1
		lm.WrapRpcRsp(-1, "insert track record failed.", payload, rsp)
		return err
	}
	payload.Id = lastId
	lm.WrapRpcRsp(1, "succeed.", payload, rsp)
	return nil
}

func (s *TrackService) DeleteTrack(hr *http.Request, req *lm.DeleteTrackRpcReq, rsp *lm.CommRpcRsp) (err error) {
	defer common.TimeCostTrack(time.Now(), "TrackService rpc", "DeleteTrack", err)
	payload := &lm.DeleteTrackRpcRsp{Id: req.Id}
	if err := lm.CheckParamsIsNil(req); err != nil {
		lm.WrapRpcRsp(2, "", payload, rsp)
		return err
	}
	err = dblogic.TkDriver.DeleteOneTrack(req.Id)
	if err != nil {
		logger.Entry().Errorf("rpc to delete track error: %v|%v", *req, err)
		payload.Id = -1
		lm.WrapRpcRsp(-1, "delete track record failed.", payload, rsp)
		return err
	}
	payload.Deleted = true
	lm.WrapRpcRsp(1, "succeed.", payload, rsp)
	return nil
}

func (s *TrackService) UpdateTrack(hr *http.Request, req *lm.UpdateTrackRpcReq, rsp *lm.CommRpcRsp) (err error) {
	defer common.TimeCostTrack(time.Now(), "TrackService rpc", "UpdateTrack", err)
	payload := &lm.UpdateTrackRpcRsp{}
	if err = lm.CheckParamsIsNil(req); err != nil {
		lm.WrapRpcRsp(2, "", payload, rsp)
		return
	}
	_, err = dblogic.TkDriver.UpdateOneTrack(req.Id, req.UpdateDoc)
	if err != nil {
		logger.Entry().Errorf("rpc to update track error: %v|%v", *req, err)
		payload.Id, payload.Changed = -1, false
		lm.WrapRpcRsp(-1, "update track record failed.", payload, rsp)
		return
	}
	payload.Id, payload.Changed = req.Id, true
	lm.WrapRpcRsp(1, "succeed.", payload, rsp)
	return
}

func (s *TrackService) SearchTrack(hr *http.Request, req *lm.SearchTrackRpcReq, rsp *lm.CommRpcRsp) (err error) {
	defer common.TimeCostTrack(time.Now(), "TrackService rpc", "SearchTrack", err)
	payload := &lm.SearchTrackRpcRsp{}
	if err = lm.CheckParamsIsNil(req); err != nil {
		lm.WrapRpcRsp(2, "", payload, rsp)
		return
	}
	var results []*models.Track
	var total int64
	conds := make(map[string]interface{})
	if req.Filter != nil {
		if req.Filter.Id != 0 {
			conds["Ftrack_id"] = req.Filter.Id
		}
		if len(req.Filter.Name) != 0 {
			conds["Ftrack_name"] = req.Filter.Name
		}
		if req.Filter.SingerId != 0 {
			conds["Fsinger_id1"] = req.Filter.SingerId
		}
		if req.Filter.AlbumId != 0 {
			conds["Falbum_id"] = req.Filter.AlbumId
		}

	}
	results, total, err = dblogic.TkDriver.GetTracksByCondition(conds, req.Start, req.Count)
	if err != nil {
		logger.Entry().Errorf("rpc to search track error: %v|%v", *req, err)
		lm.WrapRpcRsp(-1, "search track record failed.", nil, rsp)
		return
	}
	if len(results) == 0 {
		logger.Entry().Errorf("rpc to search track is empty: %v", *req)
		lm.WrapRpcRsp(-1, "search track record is empty.", nil, rsp)
		return
	}
	payload.Data.Total, payload.Data.Start, payload.Data.Count = total, req.Start, req.Count
	payload.Data.List = results
	lm.WrapRpcRsp(1, "succeed.", payload.Data, rsp)
	return
}
