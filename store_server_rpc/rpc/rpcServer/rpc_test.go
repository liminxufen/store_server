package rpcServer

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/fvbock/endless"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/store_server/dbtools/driver"
	"github.com/store_server/dbtools/models"
	"github.com/store_server/logger"
	"github.com/store_server/store_server_rpc/conf"
	"github.com/store_server/store_server_rpc/g"
	"github.com/store_server/utils/common"
	"github.com/stretchr/testify/assert"

	log "github.com/store_server/logger"
	lm "github.com/store_server/store_server_rpc/rpc/common"
)

var (
	rpcServer       *RpcSvr
	getEnvOrDefault = common.GetEnvOrDefault
	startOk         = make(chan struct{})
)

func MockRpcServer() {
	ctx, _ := context.WithCancel(context.TODO())

	server := rpc.NewServer()
	server.RegisterCodec(json.NewCodec(), "application/json")

	server.RegisterService(new(TrackService), "track")

	ul, err := NewDefaultDBEnv(ctx) //初始化各DB环境
	if err != nil {
		logger.Entry().Errorf("init mysql and mongo db engine error in mock rpc server: %v", err)
	}
	ul.Start()

	router := mux.NewRouter()
	router.Handle("/rpc", server)

	logger.Entry().Info("mock rpc server listen and serving start...")
	if err := endless.ListenAndServe(g.Config().Rpc.Listen, router); err != nil {
		log.Fatalf("start mock rpc server error: %v", err)
	}
	startOk <- struct{}{}
}

func rpcSetup() {
	//driver.DefaultTestScheme.Setup()
	//driver.DefaultTestMongoScheme.Setup()
	rpcServer, _ = NewDefaultRpcServer()
	g.Config().Rpc.Listen = ":9090"
	g.Config().MongoDb = conf.MongoDB{
		Host:     "xxx.xxx.xxx.xxx",
		Port:     6006,
		User:     "music_cms",
		Passwd:   getEnvOrDefault("MONGO_TEST_PASSWD", ""),
		DbName:   "music_cms",
		PoolSize: 100,
		Direct:   false,
		TimeOut:  3,
	}
	g.Config().Mysql = driver.DefaultTestScheme.Mysql
	go rpcServer.Start()
	//go MockRpcServer()
}

func rpcCleanup() {
	rpcServer.closeDoneChan <- struct{}{}
	rpcServer.Stop()
	//driver.DefaultTestScheme.Cleanup()
	//driver.DefaultTestMongoScheme.Cleanup()
}

func MockRPCRequest(method string, req, rsp interface{}) (err error) {
	buf, err := json.EncodeClientRequest(method, req)
	if err != nil {
		return
	}
	body := bytes.NewBuffer(buf)
	fmt.Println("mock rpc request body: ", body.String())
	port := strings.TrimPrefix(g.Config().Rpc.Listen, ":")
	r, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%s/rpc", port), body)
	if err != nil {
		return
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("X-Forwarded-For", "127.0.0.1")
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	err = json.DecodeClientResponse(resp.Body, rsp)
	return
}

func genMockTrackData() *models.Track {
	current := models.TimeNormal{time.Now()}
	track := &models.Track{
		FtrackId:           int64(20),
		FuploadTime:        current,
		FmodifyTime:        current,
		FlastestModifyTime: current,
	}
	return track
}

func TestRpcTrackCreate(t *testing.T) {
	rpcSetup()
	defer rpcCleanup()
	select {
	case <-startOk:
		break
	default:
		fmt.Println("in select loop ...")
	}
	req := &lm.CreateTrackRpcReq{
		CreateDoc: genMockTrackData(),
	}
	rsp := new(lm.CommRpcRsp)

	err := MockRPCRequest("track.Create", req, rsp)
	assert.NoError(t, err)
	assert.Equal(t, rsp.Code, 1)
	assert.Equal(t, rsp.ErrMsg, "succeed.")
}

func TestRpcTrackUpdate(t *testing.T) {
	rpcSetup()
	defer rpcCleanup()
	req := &lm.UpdateTrackRpcReq{
		Id: int64(20),
	}
	updateDoc := genMockTrackData()
	updateDoc.FtrackName, updateDoc.Ftype, updateDoc.Flanguage = "test_image_2333", int64(10), int64(100)
	req.UpdateDoc = updateDoc
	rsp := new(lm.CommRpcRsp)

	err := MockRPCRequest("track.Update", req, rsp)
	assert.NoError(t, err)
	assert.Equal(t, rsp.Code, 1)
	assert.Equal(t, rsp.ErrMsg, "succeed.")
}

func TestRpcTrackSearch(t *testing.T) {
	defer rpcCleanup()
	req := &lm.SearchTrackRpcReq{
		Start: 1,
		Count: 10,
	}
	filter := &lm.SearchTrackRpcReq_FilterInfo{}
	sort := []*lm.SearchTrackRpcReq_SortInfo{}
	req.Filter, req.Sort = filter, sort
	rsp := new(lm.CommRpcRsp)

	err := MockRPCRequest("track.Search", req, rsp)
	assert.NoError(t, err)
	assert.Equal(t, rsp.Code, 1)
	assert.Equal(t, rsp.ErrMsg, "succeed.")
}
