package rpcClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/store_server/dbtools/models"
	"github.com/store_server/logger"
	"github.com/store_server/store_server_rpc/g"
	"github.com/store_server/utils/common"

	rpcjson "github.com/gorilla/rpc/json"
	lm "github.com/store_server/store_server_rpc/rpc/common"
)

//track rpc client
type TrackRpcClient struct {
	ServiceName string
	Host        string
}

func NewTrackRpcClient() *TrackRpcClient {
	client := &TrackRpcClient{
		ServiceName: "track",
	}
	host := common.GetLocalIP()
	if len(host) == 0 {
		host = "127.0.0.1"
	} else {
		logger.Entry().Debugf("local host ip: %v", host)
	}
	client.Host = host
	return client
}

func (client *TrackRpcClient) rpcRequest(method string, req interface{}) (data []byte, err error) {
	rsp := new(lm.CommRpcRsp)
	msg, err := rpcjson.EncodeClientRequest(fmt.Sprintf("%s.%s", client.ServiceName, method), req)
	if err != nil {
		return
	}
	url := fmt.Sprintf("http://%s:%d/rpc", client.Host, g.Config().RpcPort)
	resp, err := http.Post(url, "application/json", bytes.NewReader(msg))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	err = rpcjson.DecodeClientResponse(resp.Body, rsp)
	if err != nil {
		logger.Entry().Errorf("rpcjson.DecodeClientResponse error: %v", err)
		return
	}
	if rsp.Code != 1 {
		err = fmt.Errorf("rpc code: %v, error msg: %v", rsp.Code, rsp.ErrMsg)
		return
	}
	data, err = json.Marshal(rsp.Data)
	return
}

func (client *TrackRpcClient) GenCreateTrackReq(user int, createDoc *models.Track) (req *lm.CreateTrackRpcReq) {
	req = &lm.CreateTrackRpcReq{
		CreateDoc: createDoc,
	}
	return
}

func (client *TrackRpcClient) Create(req *lm.CreateTrackRpcReq) (rsp *lm.CreateTrackRpcRsp, err error) {
	rsp = new(lm.CreateTrackRpcRsp)
	defer func() {
		if rE := recover(); rE != nil {
			err = fmt.Errorf("panic|%v", rE.(error))
		}
		if err != nil {
			logger.Entry().Errorf("track rpc client[Create]|req:%v|error:%v", *req, err)
			err = fmt.Errorf("track rpc client[Create]|req:%v|error:%v", *req, err)
		}
	}()
	var body []byte
	body, err = client.rpcRequest("Create", req)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, rsp)
	return
}

func (client *TrackRpcClient) GenUpdateTrackReq(id int64, user int,
	updateDoc *models.Track) (req *lm.UpdateTrackRpcReq) {
	req = &lm.UpdateTrackRpcReq{
		Id:        id,
		UpdateDoc: updateDoc,
	}
	return
}

func (client *TrackRpcClient) Update(req *lm.UpdateTrackRpcReq) (rsp *lm.UpdateTrackRpcRsp, err error) {
	rsp = new(lm.UpdateTrackRpcRsp)
	defer func() {
		if rE := recover(); rE != nil {
			err = fmt.Errorf("panic|%v", rE.(error))
		}
		if err != nil {
			logger.Entry().Errorf("track rpc client[Update]|req:%v|error:%v", *req, err)
			err = fmt.Errorf("track rpc client[Update]|req:%v|error:%v", *req, err)
		}
	}()
	var body []byte
	body, err = client.rpcRequest("Update", req)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, rsp)
	return
}

func (client *TrackRpcClient) GenSearchTrackReq(args map[string]interface{}) (req *lm.SearchTrackRpcReq) {
	req = &lm.SearchTrackRpcReq{}
	if startI, ok := args["start"]; ok {
		if start, ok := startI.(int64); ok {
			req.Start = start
		}
	}
	if countI, ok := args["count"]; ok {
		if count, ok := countI.(int64); ok {
			req.Count = count
		}
	}
	if filterI, ok := args["filter"]; ok {
		if filter, ok := filterI.(*lm.SearchTrackRpcReq_FilterInfo); ok {
			req.Filter = filter
		}
	}
	if sortI, ok := args["sort"]; ok {
		if sort, ok := sortI.([]*lm.SearchTrackRpcReq_SortInfo); ok {
			req.Sort = sort
		}
	}
	return
}

func (client *TrackRpcClient) Search(req *lm.SearchTrackRpcReq) (rsp *lm.SearchTrackRpcRsp, err error) {
	rsp = new(lm.SearchTrackRpcRsp)
	defer func() {
		if rE := recover(); rE != nil {
			err = fmt.Errorf("panic|%v", rE.(error))
		}
		if err != nil {
			logger.Entry().Errorf("track rpc client[Search]|req:%v|error:%v", *req, err)
			err = fmt.Errorf("track rpc client[search]|req:%v|error:%v", *req, err)
		}
	}()
	var body []byte
	body, err = client.rpcRequest("Search", req)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, rsp)
	return
}
