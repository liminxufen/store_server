package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/store_server/logger"
	"github.com/store_server/store_server_http/kits"
	"github.com/store_server/store_server_http/op"
	"github.com/store_server/utils/common"
)

/*HTTP Server相关api*/
func configServerAPI() {
	//router.Use(kits.ApiLog())         //api request log, 暂时禁用该中间件, 可能引起log组件长时间锁等待
	router.Use(kits.HttpQPSCounter()) //api request statistics
	//router.Use(kits.ApiAccessAuth())  //api access auth for whitelist

	router.GET("/store_server/config/reload", func(c *gin.Context) { //update config when need
		var (
			path string
			err  error
			ok   bool
		)
		if path, ok = c.GetQuery("path"); !ok {
			path = ""
		}
		rsp, err := ReloadConfig(path)
		if err != nil {
			logger.Entry().Error(err)
		}
		c.JSON(http.StatusOK, rsp)
	})

	router.POST("/rpc", gin.WrapH(kits.NewProxyHandle(""))) //封装对rpc的请求方式为http的反向代理

	configTracksAPI()
	configVideosAPI()
	configMatchesAPI()
	configMongosAPI()
	configEsAPI()
	configDataplatformAPI()
}

func configTracksAPI() {
	configTracksQueryAPI()
	configTracksUpdateAPI()
	configTracksInsertAPI()
	configTracksDeleteAPI()
	configTracksJoinAPI()
	configTracksOtherAPI()
}

func configEsAPI() {
	configEsSearchAPI()
	configEsUpsertAPI()
	configEsDeleteAPI()
}

//歌曲数据存储操作API定义
func configTracksQueryAPI() {
	tqr := router.Group("/store_server/tracks/query")
	{
		tqr.POST("/track", func(c *gin.Context) {
			queryReq := &op.QueryTrackReq{}
			if err := c.BindJSON(queryReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			//rsp, err := op.TracksQuery(c.Request)
			rsp, err := op.TracksQuery(queryReq)
			if err != nil {
				logger.Entry().Errorf("query track error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
		tqr.POST("/track_extra_os", func(c *gin.Context) {
			queryReq := &op.QueryTrackExtraOsReq{}
			if err := c.BindJSON(queryReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.TrackExtraOsQuery(queryReq)
			if err != nil {
				logger.Entry().Errorf("query track extra os error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
		tqr.POST("/playinfo", func(c *gin.Context) {
			pReq := &op.QueryTrackPlayReq{}
			if err := c.BindJSON(pReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.TrackPlayQuery(pReq)
			if err != nil {
				logger.Entry().Errorf("query track play info error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
		tqr.POST("/singerinfo", func(c *gin.Context) {
			sReq := &op.QueryTrackSingerReq{}
			if err := c.BindJSON(sReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.TrackSingerQuery(sReq)
			if err != nil {
				logger.Entry().Errorf("query track singer info error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
	}
}

func configTracksUpdateAPI() {
	tur := router.Group("/store_server/tracks/update")
	{
		tur.POST("/track", func(c *gin.Context) {
			updateReq := &op.UpdateTrackReq{}
			if err := c.BindJSON(updateReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.TracksUpdate(updateReq)
			if err != nil {
				logger.Entry().Errorf("update track error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
		tur.POST("/track_extra_os", func(c *gin.Context) {
			updateReq := &op.UpdateTrackExtraOsReq{}
			if err := c.BindJSON(updateReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.TrackExtraOsUpdate(updateReq)
			if err != nil {
				logger.Entry().Errorf("update track extra os error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
	}
}

func configTracksDeleteAPI() {
	tdr := router.Group("/store_server/tracks/delete")
	{
		tdr.POST("/track", func(c *gin.Context) {
			deleteReq := &op.DeleteTrackReq{}
			if err := c.BindJSON(deleteReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.TracksDelete(deleteReq)
			if err != nil {
				logger.Entry().Errorf("delete track error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
		tdr.POST("/track_extra_os", func(c *gin.Context) {
			deleteReq := &op.DeleteTrackExtraOsReq{}
			if err := c.BindJSON(deleteReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.TrackExtraOsDelete(deleteReq)
			if err != nil {
				logger.Entry().Errorf("delete track extra os error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
	}
}

func configTracksInsertAPI() {
	tir := router.Group("/store_server/tracks/insert")
	{
		tir.POST("/track", func(c *gin.Context) {
			insertReq := &op.InsertTrackReq{}
			if err := c.BindJSON(insertReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.TracksInsert(insertReq)
			if err != nil {
				logger.Entry().Errorf("insert track error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
	}
}

func configTracksJoinAPI() {
	tjg := router.Group("/store_server/tracks/join")
	{
		tjg.POST("/query", func(c *gin.Context) {
			jqReq := &op.JoinQueryTrackReq{}
			if err := c.BindJSON(jqReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.TracksJoinQuery(jqReq)
			if err != nil {
				logger.Entry().Errorf("join query track error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
	}
}

func configTracksOtherAPI() {
	tog := router.Group("/store_server/tracks")
	{
		tog.POST("/url", func(c *gin.Context) {
			qReq := &op.QueryTrackPlayURLReq{}
			if err := c.BindJSON(qReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.TrackPlayURLQuery(qReq)
			if err != nil {
				logger.Entry().Errorf("query track play url error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
	}
}

//视频数据存储操作API定义
func configVideosAPI() {
	vqr := router.Group("/store_server/videos/query")
	{
		vqr.POST("/video", func(c *gin.Context) {
			queryReq := &op.QueryVideoReq{}
			if err := c.BindJSON(queryReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.VideosQuery(queryReq)
			if err != nil {
				logger.Entry().Errorf("query video error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
		vqr.POST("/video_extra_os", func(c *gin.Context) {
			queryReq := &op.QueryVideoExtraOsReq{}
			if err := c.BindJSON(queryReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.VideoExtraOsQuery(queryReq)
			if err != nil {
				logger.Entry().Errorf("query video extra os error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
		vqr.POST("/video_singer_track", func(c *gin.Context) {
			queryReq := &op.QueryVideoSingerTrackReq{}
			if err := c.BindJSON(queryReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.VideoSingerTrackQuery(queryReq)
			if err != nil {
				logger.Entry().Errorf("query video singer track error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
	}
	vjg := router.Group("/store_server/videos/join")
	{
		vjg.POST("/query", func(c *gin.Context) {
			jqReq := &op.JoinQueryVideoReq{}
			if err := c.BindJSON(jqReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.VideosJoinQuery(jqReq)
			if err != nil {
				logger.Entry().Errorf("join query video error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
	}
}

//视频、歌曲、艺人等匹配逻辑API定义
func configMatchesAPI() {
	mrs := router.Group("/store_server/matches/query")
	{
		mrs.POST("/video", func(c *gin.Context) {
			queryReq := &op.QueryVideoMatchInfoReq{}
			if err := c.BindJSON(queryReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.VideoMatchInfoQuery(queryReq)
			if err != nil {
				logger.Entry().Errorf("query video match info error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
		mrs.POST("/track", func(c *gin.Context) {
			if err := c.BindJSON(""); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			c.JSON(http.StatusOK, "")
		})
		mrs.POST("/album", func(c *gin.Context) {
			if err := c.BindJSON(""); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			c.JSON(http.StatusOK, "")
		})
		mrs.POST("/singer", func(c *gin.Context) {
			if err := c.BindJSON(""); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			c.JSON(http.StatusOK, "")
		})
	}
}

//mongo数据操作API定义
func configMongosAPI() {
	mqs := router.Group("/store_server/mongo/query")
	{
		mqs.POST("/external_resources", func(c *gin.Context) {
			queryReq := &op.QueryExternalResourcesReq{}
			if err := c.BindJSON(queryReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.ExternalResourcesQuery(queryReq)
			if err != nil {
				logger.Entry().Errorf("query external resources info error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
		mqs.POST("/publish_albums", func(c *gin.Context) {
			if err := c.BindJSON(""); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			c.JSON(http.StatusOK, "")
		})
		mqs.POST("/track", func(c *gin.Context) {
			if err := c.BindJSON(""); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			c.JSON(http.StatusOK, "")
		})
		mqs.POST("/singer", func(c *gin.Context) {
			if err := c.BindJSON(""); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			c.JSON(http.StatusOK, "")
		})
	}
	mus := router.Group("/store_server/mongo/update")
	{
		mus.POST("/external_resources", func(c *gin.Context) {
			updateReq := &op.UpdateExternalResourcesReq{}
			if err := c.BindJSON(updateReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.ExternalResourcesUpdate(updateReq)
			if err != nil {
				logger.Entry().Errorf("update external resources info error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
		mus.POST("/publish_albums", func(c *gin.Context) {
			if err := c.BindJSON(""); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			c.JSON(http.StatusOK, "")
		})
		mus.POST("/track", func(c *gin.Context) {
			if err := c.BindJSON(""); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			c.JSON(http.StatusOK, "")
		})
		mus.POST("/singer", func(c *gin.Context) {
			if err := c.BindJSON(""); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			c.JSON(http.StatusOK, "")
		})
	}
	mcs := router.Group("/store_server/mongo/insert")
	{
		mcs.POST("/external_resources", func(c *gin.Context) {
			insertReq := &op.InsertExternalResourcesReq{}
			if err := c.BindJSON(insertReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.ExternalResourcesInsert(insertReq)
			if err != nil {
				logger.Entry().Errorf("insert external resources info error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
		mcs.POST("/publish_albums", func(c *gin.Context) {
			if err := c.BindJSON(""); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			c.JSON(http.StatusOK, "")
		})
		mcs.POST("/track", func(c *gin.Context) {
			if err := c.BindJSON(""); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			c.JSON(http.StatusOK, "")
		})
		mcs.POST("/singer", func(c *gin.Context) {
			if err := c.BindJSON(""); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			c.JSON(http.StatusOK, "")
		})
	}
	mds := router.Group("/store_server/mongo/delete")
	{
		mds.POST("/external_resources", func(c *gin.Context) {
			deleteReq := &op.DeleteExternalResourcesReq{}
			if err := c.BindJSON(deleteReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.ExternalResourcesDelete(deleteReq)
			if err != nil {
				logger.Entry().Errorf("delete external resources info error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
		mds.POST("/publish_albums", func(c *gin.Context) {
			if err := c.BindJSON(""); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			c.JSON(http.StatusOK, "")
		})
		mds.POST("/track", func(c *gin.Context) {
			if err := c.BindJSON(""); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			c.JSON(http.StatusOK, "")
		})
		mds.POST("/singer", func(c *gin.Context) {
			if err := c.BindJSON(""); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			c.JSON(http.StatusOK, "")
		})
	}
}

//elasticsearch数据操作API定义
func configEsSearchAPI() {
	ess := router.Group("/store_server/es/search")
	{
		ess.POST("/tracks", func(c *gin.Context) {
			searchReq := &op.SearchTracksReq{}
			if err := c.BindJSON(searchReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.TracksSearch(searchReq)
			if err != nil {
				logger.Entry().Errorf("search tracks error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
		ess.POST("/albums", func(c *gin.Context) {
			searchReq := &op.SearchAlbumsReq{}
			if err := c.BindJSON(searchReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.AlbumsSearch(searchReq)
			if err != nil {
				logger.Entry().Errorf("search albums error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
		ess.POST("/singers", func(c *gin.Context) {
			searchReq := &op.SearchSingersReq{}
			if err := c.BindJSON(searchReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.SingersSearch(searchReq)
			if err != nil {
				logger.Entry().Errorf("search singers error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
		ess.POST("/videos", func(c *gin.Context) {
			searchReq := &op.SearchVideosReq{}
			if err := c.BindJSON(searchReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.VideosSearch(searchReq)
			if err != nil {
				logger.Entry().Errorf("search videos error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
	}
}

func configEsUpsertAPI() {
	esu := router.Group("/store_server/es/upsert")
	{
		esu.POST("/tracks", func(c *gin.Context) {
			upsertReq := &op.UpsertTracksReq{}
			if err := c.BindJSON(upsertReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.TracksUpsert(upsertReq)
			if err != nil {
				logger.Entry().Errorf("upsert tracks error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
		esu.POST("/albums", func(c *gin.Context) {
			upsertReq := &op.UpsertAlbumsReq{}
			if err := c.BindJSON(upsertReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.AlbumsUpsert(upsertReq)
			if err != nil {
				logger.Entry().Errorf("upsert albums error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
		esu.POST("/singers", func(c *gin.Context) {
			upsertReq := &op.UpsertSingersReq{}
			if err := c.BindJSON(upsertReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.SingersUpsert(upsertReq)
			if err != nil {
				logger.Entry().Errorf("upsert singers error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
		esu.POST("/videos", func(c *gin.Context) {
			upsertReq := &op.UpsertVideosReq{}
			if err := c.BindJSON(upsertReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.VideosUpsert(upsertReq)
			if err != nil {
				logger.Entry().Errorf("upsert videos error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
	}
}

func configEsDeleteAPI() {
	esd := router.Group("/store_server/es/delete")
	{
		esd.POST("/tracks", func(c *gin.Context) {
			deleteReq := &op.DeleteTrackDocReq{}
			if err := c.BindJSON(deleteReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.TrackDocDelete(deleteReq)
			if err != nil {
				logger.Entry().Errorf("delete tracks error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
		esd.POST("/albums", func(c *gin.Context) {
			deleteReq := &op.DeleteAlbumDocReq{}
			if err := c.BindJSON(deleteReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.AlbumDocDelete(deleteReq)
			if err != nil {
				logger.Entry().Errorf("delete albums error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
		esd.POST("/singers", func(c *gin.Context) {
			deleteReq := &op.DeleteSingerDocReq{}
			if err := c.BindJSON(deleteReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.SingerDocDelete(deleteReq)
			if err != nil {
				logger.Entry().Errorf("delete singers error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
		esd.POST("/videos", func(c *gin.Context) {
			deleteReq := &op.DeleteVideoDocReq{}
			if err := c.BindJSON(deleteReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.VideoDocDelete(deleteReq)
			if err != nil {
				logger.Entry().Errorf("delete videos error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
	}
}

//dataplatform数据操作API定义
func configDataplatformAPI() {
	dps := router.Group("/store_server/dataplatform/search")
	{
		dps.POST("/tracks", func(c *gin.Context) {
			searchReq := &op.SearchTrackDReq{}
			//if err := c.BindJSON(searchReq); err != nil {
			data, _ := c.GetRawData()
			if err := common.UnmarshalWithNumber(data, searchReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.TrackSearchFromDataplatForm(searchReq)
			if err != nil {
				logger.Entry().Errorf("search tracks from dataplatform error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
		dps.POST("/albums", func(c *gin.Context) {
			searchReq := &op.SearchAlbumDReq{}
			//if err := c.BindJSON(searchReq); err != nil {
			data, _ := c.GetRawData()
			if err := common.UnmarshalWithNumber(data, searchReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.AlbumSearchFromDataplatForm(searchReq)
			if err != nil {
				logger.Entry().Errorf("search albums from dataplatform error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
		dps.POST("/singers", func(c *gin.Context) {
			searchReq := &op.SearchSingerDReq{}
			//if err := c.BindJSON(searchReq); err != nil {
			data, _ := c.GetRawData()
			if err := common.UnmarshalWithNumber(data, searchReq); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			rsp, err := op.SingerSearchFromDataplatForm(searchReq)
			if err != nil {
				logger.Entry().Errorf("search singers from dataplatform error: %v", err)
			}
			c.JSON(http.StatusOK, rsp)
		})
	}
}
