package http

/*负责与CMS平台的交互，包括音视频文件的上传、更新、查找、搜索等*/
import (
	"context"
	"fmt"
	"github.com/DeanThompson/ginpprof"
	"github.com/fvbock/endless"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/contrib/expvar"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/store_server/logger"
	log "github.com/store_server/logger"
	"github.com/store_server/store_server_http/g"
	"github.com/store_server/store_server_http/kits"
	"github.com/store_server/store_server_http/op"
	"github.com/zsais/go-gin-prometheus"
	_ "net/http/pprof" //初始化pprof
	"strings"
	"time"
)

var (
	router *gin.Engine
)

//初始化白名单
func InitIpWhiteList(s string) {
	if len(s) == 0 {
		return
	}
	kits.IPWhiteLst = strings.Split(s, "|")
}

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

//store server define
type StoreServerHttp struct {
	closeChan     chan struct{}
	closeDoneChan chan struct{}
}

func NewDefaultStoreServerHttp() (sh *StoreServerHttp, err error) {
	sh = &StoreServerHttp{
		closeChan:     make(chan struct{}, 1),
		closeDoneChan: make(chan struct{}, 1),
	}
	return sh, nil
}

func (sh *StoreServerHttp) Stop() (err error) {
	sh.closeChan <- struct{}{}
	select {
	case <-sh.closeDoneChan:
	}
	return
}

//启动HTTP服务
func (sh *StoreServerHttp) Start() (err error) {
	ctx, cancel := context.WithCancel(context.TODO())
	if !g.Config().Http.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	if len(g.Config().IpWhiteList) == 0 {
		logger.Entry().Info("【please set ip white list for store_server api access permission】")
	} else {
		logger.Entry().Infof("store http server start, ip white list is: %v", g.Config().IpWhiteList)
		InitIpWhiteList(g.Config().IpWhiteList)
	}
	router = gin.Default()
	ginpprof.Wrap(router)

	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)
	router.GET("/debug/vars", expvar.Handler())
	router.GET("/store_server/prometheus_metrics", gin.WrapH(promhttp.Handler()))

	ul, err := NewDefaultDBEnv(ctx)
	if err != nil {
		logger.Entry().Errorf("init mysql, mongo and elastic engine error in starting store server http: %v", err)
	}
	ul.Start()

	/*err = kits.InitInfluxEnv(ctx)
	if err != nil {
		logger.Entry().Errorf("init influxdb client error in starting http server: %v", err)
	}*/

	kits.HTTPCounter = kits.NewCounterService(ctx) //开启QPS统计

	configServerAPI()
	go op.ExportAllData(ctx) //导出数据

	addr := g.Config().Http.Listen
	if addr == "" {
		err = fmt.Errorf("address in config is empty")
		return
	}

	go func() {
		select {
		case <-sh.closeChan:
			logger.Entry().Info("http recv sigout and exit...")
			cancel()
			ul.Stop()
			sh.closeDoneChan <- struct{}{}
		}
	}()

	//go router.Run(addr)
	if err := endless.ListenAndServe(addr, router); err != nil {
		sentry.CaptureException(err)
		sentry.Flush(time.Second)
		log.Fatalf("start store server http error: %v", err)
	}
	return
}
