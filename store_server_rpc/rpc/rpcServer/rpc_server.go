package rpcServer

import (
	"context"
	"time"

	"github.com/fvbock/endless"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/store_server/logger"
	"github.com/store_server/store_server_rpc/g"

	log "github.com/store_server/logger"
)

func flushSentry(ctx context.Context) {
	tk := time.NewTicker(3 * time.Second)
	defer tk.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tk.C:
			sentry.Flush(10 * time.Second)
		}
	}
}

//rpc server
type RpcSvr struct {
	closeChan     chan struct{}
	closeDoneChan chan struct{}
}

func NewDefaultRpcServer() (s *RpcSvr, err error) {
	s = &RpcSvr{
		closeChan:     make(chan struct{}, 1),
		closeDoneChan: make(chan struct{}, 1),
	}
	return s, nil
}

func (s *RpcSvr) Stop() (err error) {
	s.closeChan <- struct{}{}
	select {
	case <-s.closeDoneChan:
	}
	return
}

func (s *RpcSvr) Start() (err error) {
	ctx, cancel := context.WithCancel(context.TODO())

	server := rpc.NewServer()
	server.RegisterCodec(json.NewCodec(), "application/json")

	server.RegisterService(new(TrackService), "track")

	ul, err := NewDefaultDBEnv(ctx) //初始化各DB环境
	if err != nil {
		logger.Entry().Errorf("init mysql and mongo db engine error in starting rpc server: %v", err)
	}
	ul.Start()

	router := mux.NewRouter()
	router.Handle("/rpc", server)
	router.Handle("/store_server/prometheus_metrics", promhttp.Handler())

	go func() {
		select {
		case <-s.closeChan:
			logger.Entry().Info("rpc server recv sigout and exit...")
			cancel()
			//ul.Stop(ctx)
			s.closeDoneChan <- struct{}{}
		}
	}()

	logger.Entry().Info("media rpc server listen and serving start...")
	if err := endless.ListenAndServe(g.Config().Rpc.Listen, router); err != nil {
		sentry.CaptureException(err)
		sentry.Flush(time.Second)
		log.Fatalf("start media rpc server error: %v", err)
	}
	return
}
