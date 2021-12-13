package op

import (
	"context"
	"fmt"
	"github.com/store_server/dbtools/dblogic"
	m "github.com/store_server/dbtools/models"
	"github.com/store_server/logger"
	"github.com/store_server/store_server_http/g"
	"github.com/store_server/store_server_http/kits"
	"github.com/store_server/utils/errors"
	"reflect"
	"strings"
	"sync"
	"time"
)

/************************ 视频查询相关 ***************************/
//query video request
type QueryVideoReq struct {
	RawSql   string                 `json:"rawSql"`
	Ids      []int64                `json:"ids"`
	Page     int64                  `json:"page,omitempty"`
	PageSize int64                  `json:"pageSize,omitempty"`
	Fields   map[string]interface{} `json:"fields,omitempty"`
}

//query video response
type QueryVideoRsp struct {
	Videos []*m.Video `json:"videos"`
	Total  int64      `json:"total,omitempty"`
}

func VideosQuery(req *QueryVideoReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.VideosQuery", &err, logger.Entry())
	ret := QueryVideoRsp{}
	var videos []*m.Video
	if len(req.RawSql) != 0 {
		videos, ret.Total, err = dblogic.VoDriver.ExecRawQuerySql4Video(req.RawSql, req.Page, req.PageSize)
	} else if len(req.Ids) != 0 {
		videos, ret.Total, err = dblogic.VoDriver.GetVideosByIds(req.Ids)
	} else { //others query condition
		if len(req.Fields) == 0 && req.Page == 0 && req.PageSize == 0 {
			logger.Entry().Errorf("query videos fields conditions is nil")
			rsp = kits.APIWrapRsp(kits.ErrOther, "query videos fields conditions is invalid", ret)
			return
		}
		videos, ret.Total, err = dblogic.VoDriver.GetVideosByCondition(req.Fields, req.Page, req.PageSize)
	}
	if err != nil {
		logger.Entry().Errorf("query videos error: %v|request: %v", err, *req)
		ret.Total = 0
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	ret.Videos = videos
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

//query video extra os request
type QueryVideoExtraOsReq struct {
	RawSql   string                 `json:"rawSql"`
	Id       int64                  `json:"id"`
	Region   int64                  `json:"region"`
	Page     int64                  `json:"page,omitempty"`
	PageSize int64                  `json:"pageSize,omitempty"`
	Fields   map[string]interface{} `json:"fields,omitempty"`
}

//query video extra os response
type QueryVideoExtraOsRsp struct {
	Videos []*m.VideoExtraOs `json:"videoExtraOs"`
	Total  int64             `json:"total,omitempty"`
}

func VideoExtraOsQuery(req *QueryVideoExtraOsReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.VideoExtraOsQuery", &err, logger.Entry())
	ret := QueryVideoExtraOsRsp{}
	var videos []*m.VideoExtraOs
	if len(req.RawSql) != 0 {
		videos, ret.Total, err = dblogic.VoDriver.ExecRawQuerySql4VideoExtraOs(req.RawSql, req.Page, req.PageSize)
	} else if req.Id != 0 {
		videos, ret.Total, err = dblogic.VoDriver.GetVideoExtraOs(req.Id, req.Region)
	} else {
		if len(req.Fields) == 0 && req.Page == 0 && req.PageSize == 0 {
			logger.Entry().Errorf("query video extra os fields conditions is nil")
			rsp = kits.APIWrapRsp(kits.ErrOther, "query video extra os fields conditions is invalid", ret)
			return
		}
		videos, ret.Total, err = dblogic.VoDriver.GetVideoExtraOsByCondition(req.Fields, req.Page, req.PageSize)
	}
	if err != nil {
		logger.Entry().Errorf("query video extra os error: %v|request: %v", err, *req)
		ret.Total = 0
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	ret.Videos = videos
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

//query video singer track request
type QueryVideoSingerTrackReq struct {
	RawSql   string                 `json:"rawSql"`
	Id       int64                  `json:"id"`
	Region   int64                  `json:"region"`
	Page     int64                  `json:"page,omitempty"`
	PageSize int64                  `json:"pageSize,omitempty"`
	Fields   map[string]interface{} `json:"fields,omitempty"`
}

//query video singer track response
type QueryVideoSingerTrackRsp struct {
	Videos []*m.VideoSingerTrack `json:"videoSingerTrack"`
	Total  int64                 `json:"total,omitempty"`
}

func VideoSingerTrackQuery(req *QueryVideoSingerTrackReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.VideoSingerTrackQuery", &err, logger.Entry())
	ret := QueryVideoSingerTrackRsp{}
	var videos []*m.VideoSingerTrack
	if len(req.RawSql) != 0 {
		videos, ret.Total, err = dblogic.VoDriver.ExecRawQuerySql4VideoSingerTrack(req.RawSql, req.Page, req.PageSize)
	} else if req.Id != 0 {
		videos, ret.Total, err = dblogic.VoDriver.GetVideoSingerTrack(req.Id, req.Region)
	} else {
		if len(req.Fields) == 0 && req.Page == 0 && req.PageSize == 0 {
			logger.Entry().Errorf("query video singer track fields conditions is nil")
			rsp = kits.APIWrapRsp(kits.ErrOther, "query video singer track fields conditions is invalid", ret)
			return
		}
		videos, ret.Total, err = dblogic.VoDriver.GetVideoSingerTrackByCondition(req.Fields, req.Page, req.PageSize)
	}
	if err != nil {
		logger.Entry().Errorf("query video singer track error: %v|request: %v", err, *req)
		ret.Total = 0
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	ret.Videos = videos
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

/************************ 视频联合查询相关 ***************************/
//join query video request
type JoinQueryVideoReq struct {
	RawSql   string `json:"rawSql"`
	Page     int64  `json:"page,omitempty"`
	PageSize int64  `json:"pageSize,omitempty"`
}

//join query video response
type JoinQueryVideoRsp struct {
	Results [][]interface{} `json:"results"`
	Total   int64           `json:"total,omitempty"`
}

func VideosJoinQuery(req *JoinQueryVideoReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.VideosJoinQuery", &err, logger.Entry())
	ret := JoinQueryVideoRsp{}
	results, err := dblogic.VoDriver.JoinQueryWithRawSql(req.RawSql, req.Page, req.PageSize)
	if err != nil {
		logger.Entry().Errorf("join query videos error: %v|request: %v", err, *req)
		ret.Total = 0
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
	}
	ret.Results, ret.Total = results, int64(len(results))
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

/************************ 视频匹配相关 ***************************/
var (
	fullTracks            = []*matchTrackInfo{}
	fullSingers           = []*matchSingerInfo{}
	fullVideoSingerTracks = []*matchVideoSingerTrackInfo{}
)

const (
	mStatusNoMatch        = "to be matched"        //待匹配
	mStatusOnlyArtist     = "only artist"          //仅艺人匹配
	mStatusArtistAndTrack = "both artist and song" //歌曲和艺人匹配
)

//导出曲库全量歌曲数据
func exportAllTracks() error {
	sql := "select Ftrack_id, Ftrack_name, Fsinger_id1, Fsinger_id2, Fsinger_id3 from t_track;"
	res, err := dblogic.TkDriver.ExportAllTracks(sql)
	if err != nil {
		logger.Entry().Errorf("export all tracks timely error: %v", err)
		return err
	}
	fullTracks = fullTracks[0:0]
	for _, tk := range res {
		if len(tk) == 5 {
			tkInfo := &matchTrackInfo{
				Id:        tk[0],
				Name:      tk[1],
				SingerId1: tk[2],
				SingerId2: tk[3],
				SingerId3: tk[4],
			}
			fullTracks = append(fullTracks, tkInfo)
		}
	}
	logger.Entry().Debugf("--- export all tracks count: %v ---", len(fullTracks))
	return nil
}

//导出曲库全量艺人数据
func exportAllSingers() error {
	res := []string{}
	fullSingers = fullSingers[0:0]
	for _, sr := range res {
		if len(sr) == 3 {
			srInfo := &matchSingerInfo{
				Id:        sr[0],
				Name:      sr[1],
				AliasName: sr[2],
			}
			fullSingers = append(fullSingers, srInfo)
		}
	}
	logger.Entry().Debugf("--- export all singers count: %v ---", len(fullSingers))
	return nil
}

//导出曲库全量视频关联歌曲艺人数据
func exportAllVideoSingerTracks() error {
	sql := "select Flocal_v_id, Fsinger_id, Ftrack_id from t_video_singer_track;"
	res, err := dblogic.VoDriver.ExportAllVideoSingerTracks(sql)
	if err != nil {
		logger.Entry().Errorf("export all video singer tracks timely error: %v", err)
		return err
	}
	fullVideoSingerTracks = fullVideoSingerTracks[0:0]
	for _, vst := range res {
		if len(vst) == 3 {
			vstInfo := &matchVideoSingerTrackInfo{
				VideoId:  vst[0],
				SingerId: vst[1],
				TrackId:  vst[2],
			}
			fullVideoSingerTracks = append(fullVideoSingerTracks, vstInfo)
		}
	}
	logger.Entry().Debugf("--- export all video singer tracks count: %v ---", len(fullVideoSingerTracks))
	return nil
}

//定期导出全量曲库数据
func ExportAllData(ctx context.Context) {
	//服务初始化时是否导出全量数据到内存
	if g.Config().ExportAllOpen {
		exportAllTracks()
		exportAllSingers()
		exportAllVideoSingerTracks()
	}
	tk := time.NewTicker(1 * 24 * time.Hour)
	defer tk.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tk.C:
			exportAllTracks()
			exportAllSingers()
			exportAllVideoSingerTracks()
		}
	}
}

func MatchedTrackIsRelateMV(id interface{}) (bool, error) {
	/*conds := map[string]interface{}{"Ftrack_id": id}
	_, total, err := dblogic.VoDriver.GetVideoSingerTrackByCondition(conds, 1, 10000)
	if err != nil {
		logger.Entry().Errorf("get video singer track by condition: %v|error: %v", conds, err)
		return false, err
	}
	if total <= 0 {
		return false, nil
	}
	return true, nil*/
	relate := false
	for _, vst := range fullVideoSingerTracks {
		if reflect.DeepEqual(vst.TrackId, id) {
			relate = true
			break
		}
	}
	if relate {
		return true, nil
	}
	return false, nil
}

func getAllTracksWithSinger(sid interface{}) []*matchTrackInfo {
	tracks := make([]*matchTrackInfo, 0)
	/*inCh := make(chan *matchTrackInfo, 1000)
	outCh := make(chan *matchTrackInfo)
	go func() {
		for i := 0; i < len(fullTracks); i++ {
			inCh <- fullTracks[i]
		}
		close(inCh)
	}()
	go matchAllTracksWithSinger(sid, inCh, outCh)
	for tk := range outCh {
		tracks = append(tracks, tk)
	}*/
	for _, tk := range fullTracks {
		if reflect.DeepEqual(sid, tk.SingerId1) || reflect.DeepEqual(sid, tk.SingerId2) ||
			reflect.DeepEqual(sid, tk.SingerId3) {
			tracks = append(tracks, tk)
		}
	}
	return tracks
}

func matchAllTracksWithSinger(sid interface{}, inCh, outCh chan *matchTrackInfo) error {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for tk := range inCh {
				if reflect.DeepEqual(sid, tk.SingerId1) || reflect.DeepEqual(sid, tk.SingerId2) ||
					reflect.DeepEqual(sid, tk.SingerId3) {
					outCh <- tk
				}
			}
		}()
	}
	wg.Wait()
	close(outCh)
	return nil
}

func doMatchTrack(name string, inCh, outCh chan *matchTrackInfo) error {
	var wg sync.WaitGroup
	name = strings.TrimSpace(name)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for tk := range inCh {
				tkName := ""
				if n, ok := tk.Name.(string); ok {
					tkName = n
				}
				if name == strings.TrimSpace(tkName) { //匹配
					//排除关联mv的歌曲
					isRelate, _ := MatchedTrackIsRelateMV(tk.Id)
					if !isRelate {
						outCh <- tk
						logger.Entry().Debugf("matched track[id: %v|name: %v]", tk.Id, tk.Name)
					} else {
						logger.Entry().Debugf("matched track[id: %v|name: %v] related by other videos", tk.Id, tk.Name)
					}
				} else { //不匹配

				}
			}
		}()
	}
	wg.Wait()
	close(outCh)
	return nil
}

//歌曲匹配操作
func matchTrackOperation(name string) ([]*matchTrackInfo, error) {
	matchInfo := make([]*matchTrackInfo, 0)
	inCh := make(chan *matchTrackInfo, 1000)
	outCh := make(chan *matchTrackInfo)
	go func() {
		for i := 0; i < len(fullTracks); i++ {
			inCh <- fullTracks[i]
		}
		close(inCh)
	}()
	go doMatchTrack(name, inCh, outCh)
	for matched := range outCh {
		matchInfo = append(matchInfo, matched)
	}
	return matchInfo, nil
}

func doMatchSinger(name string, inCh, outCh chan *matchSingerInfo) error {
	var wg sync.WaitGroup
	name = strings.TrimSpace(name)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for sr := range inCh {
				srName := ""
				if n, ok := sr.Name.(string); ok {
					srName = n
				}
				if name == strings.TrimSpace(srName) { //匹配
					outCh <- sr
					logger.Entry().Debugf("matched singer[id: %v|name: %v]", sr.Id, sr.Name)
				} else { //不匹配

				}
			}
		}()
	}
	wg.Wait()
	close(outCh)
	return nil
}

//艺人匹配操作
func matchSingerOperation(name string) ([]*matchSingerInfo, error) {
	matchInfo := make([]*matchSingerInfo, 0)
	inCh := make(chan *matchSingerInfo, 1000)
	outCh := make(chan *matchSingerInfo)
	go func() {
		for i := 0; i < len(fullSingers); i++ {
			inCh <- fullSingers[i]
		}
		close(inCh)
	}()
	go doMatchSinger(name, inCh, outCh)
	for matched := range outCh {
		matchInfo = append(matchInfo, matched)
	}
	return matchInfo, nil
}

//仅mv进行歌曲匹配
func matchOnlyMVTrack(names []string) ([]*matchTrackInfo, error) {
	tks := make([]*matchTrackInfo, 0)
	for _, name := range names {
		matchInfo, _ := matchTrackOperation(name)
		if len(matchInfo) > 0 { //匹配
			tks = append(tks, matchInfo...)
		}
	}
	return tks, nil
}

//全部视频进行艺人匹配
func matchAllVideoSinger(names []string) ([]*matchSingerInfo, map[string][]*matchSingerInfo, error) {
	sgs := make([]*matchSingerInfo, 0)
	unmatchInfos := make(map[string][]*matchSingerInfo)
	for _, name := range names {
		matchInfo, _ := matchSingerOperation(name)
		if len(matchInfo) == 1 { //匹配  >1 or ==0 不匹配
			sgs = append(sgs, matchInfo...)
		} else if len(matchInfo) > 0 {
			if _, ok := unmatchInfos[name]; !ok {
				unmatchInfos[name] = make([]*matchSingerInfo, 0)
			}
			unmatchInfos[name] = append(unmatchInfos[name], matchInfo...)
		}
	}
	return sgs, unmatchInfos, nil
}

//对匹配的艺人和歌曲进行校验
func checkMatchedSingerAndTrack(vType string, matchSingers *[]*matchSingerInfo,
	matchTracks *[]*matchTrackInfo) ([]*matchVideoInfo, string, error) {
	if matchSingers == nil || matchTracks == nil {
		return nil, mStatusNoMatch, errors.New("invalid match singers or tracks info pointer")
	}
	matchInfo := make([]*matchVideoInfo, 0)
	if len(*matchSingers) == 0 { //无艺人匹配
		return matchInfo, mStatusNoMatch, nil
	}
	resultStatus := mStatusOnlyArtist
	if len(*matchSingers) > 1 {
		resultStatus = mStatusNoMatch
	}
	allMatch := false
	for _, msi := range *matchSingers { //在艺人匹配的情况下校验歌曲是否匹配
		if msi == nil {
			continue
		}
		mvInfo := &matchVideoInfo{
			matchSingerInfo: matchSingerInfo{
				Id:   msi.Id,
				Name: msi.Name,
			},
			Tracks: make([]*matchTrackInfo, 0),
		}
		status := mStatusOnlyArtist
		/*sid, err := common.Interface2Int64(msi.Id)
		if err != nil {
			logger.Entry().Errorf("singer id[%v] from interface to int64 error: %v", msi.Id, err)
			continue
		}
		tks, err := dblogic.SrDriver.ExportAllTracksWithSinger(sid)
		if err != nil {
			logger.Entry().Errorf("export all track with singer id[%v] error: %v", msi.Id, err)
		}*/
		tks := getAllTracksWithSinger(msi.Id)
		singleMatch := false
		for _, tk := range tks {
			for _, mtk := range *matchTracks {
				/*mid, err := common.Interface2Int64(mtk.Id)
				if err != nil {
					logger.Entry().Errorf("matched track id[%v] from interface to int64, error: %v", mtk.Id, err)
				}
				if tk.FtrackId == mid { //歌曲和艺人匹配*/
				if reflect.DeepEqual(tk.Id, mtk.Id) { //歌曲和艺人匹配
					singleMatch = true
					status = mStatusArtistAndTrack
					mvInfo.Tracks = append(mvInfo.Tracks, mtk)
				} else { //歌曲和艺人不匹配

				}
			}
		}
		if singleMatch {
			allMatch = true
		}
		mvInfo.Status = status
		matchInfo = append(matchInfo, mvInfo)
	}
	if allMatch {
		resultStatus = mStatusArtistAndTrack
	}
	return matchInfo, resultStatus, nil
}

//对未匹配的艺人进行校验(多个匹配情形)
func checkUnmatchedSingers(vType string, matchTracks *[]*matchTrackInfo, matchInfos *[]*matchVideoInfo,
	unmatchInfos map[string][]*matchSingerInfo) (string, error) {
	if vType != "mv" {
		return "", nil
	}
	if matchTracks == nil || matchInfos == nil {
		return "", nil
	}
	if len(unmatchInfos) == 0 {
		return "", nil
	}
	allMatch, resultStatus := false, ""
	for _, unmatchSingers := range unmatchInfos {
		for _, msi := range unmatchSingers {
			if msi == nil {
				continue
			}
			mvInfo := &matchVideoInfo{
				matchSingerInfo: matchSingerInfo{
					Id:   msi.Id,
					Name: msi.Name,
				},
				Tracks: make([]*matchTrackInfo, 0),
			}
			status := mStatusOnlyArtist
			/*sid, err := common.Interface2Int64(msi.Id)
			if err != nil {
				logger.Entry().Errorf("singer id[%v] from interface to int64 error: %v", msi.Id, err)
				continue
			}
			tks, err := dblogic.SrDriver.ExportAllTracksWithSinger(sid)
			if err != nil {
				logger.Entry().Errorf("export all track with singer id[%v] error: %v", msi.Id, err)
			}*/
			tks := getAllTracksWithSinger(msi.Id)
			singleMatch := false
			for _, tk := range tks {
				for _, mtk := range *matchTracks {
					/*mid, err := common.Interface2Int64(mtk.Id)
					if err != nil {
						logger.Entry().Errorf("matched track id[%v] from interface to int64, error: %v", mtk.Id, err)
					}
					if tk.FtrackId == mid { //歌曲和艺人匹配*/
					if reflect.DeepEqual(tk.Id, mtk.Id) { //歌曲和艺人匹配
						singleMatch = true
						status = mStatusArtistAndTrack
						mvInfo.Tracks = append(mvInfo.Tracks, mtk)
					}
				}
			}
			if singleMatch {
				allMatch = true
			}
			mvInfo.Status = status
			*matchInfos = append(*matchInfos, mvInfo)
		}
	}
	if allMatch {
		resultStatus = mStatusArtistAndTrack
	}
	return resultStatus, nil
}

type matchVideoSingerTrackInfo struct {
	VideoId  interface{} `json:"-"`
	SingerId interface{} `json:"-"`
	TrackId  interface{} `json:"-"`
}

//query video match info request
type QueryVideoMatchInfoReq struct {
	Vid    int64    `json:"vid"`
	Type   string   `json:"type"`
	Vname  []string `json:"vname"`
	Singer []string `json:"singer"`
}

type matchTrackInfo struct {
	Id        interface{} `json:"id"`
	Name      interface{} `json:"name"`
	SingerId1 interface{} `json:"-"`
	SingerId2 interface{} `json:"-"`
	SingerId3 interface{} `json:"-"`
}

type matchSingerInfo struct {
	Id        interface{} `json:"id"`
	Name      interface{} `json:"name"`
	AliasName interface{} `json:"-"`
}

type matchVideoInfo struct {
	matchSingerInfo
	Tracks []*matchTrackInfo `json:"tracks"`
	Status string            `json:"status"`
}

//query video match info response
type QueryVideoMatchInfoRsp struct {
	MatchInfos []*matchVideoInfo `json:"matchInfos"`
	Status     string            `json:"status,omitempty"` //both artist and song, only artist
	Vid        int64             `json:"vid"`
}

func VideoMatchInfoQuery(req *QueryVideoMatchInfoReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.VideoMatchInfoQuery", &err, logger.Entry())
	ret := QueryVideoMatchInfoRsp{Vid: req.Vid}
	if len(req.Type) == 0 || len(req.Singer) == 0 {
		rsp = kits.APIWrapRsp(kits.ErrParams, "query video match info conditions is nil", ret)
		return
	}
	matchSingers, unmatchSingers, err := matchAllVideoSinger(req.Singer)
	if err != nil {
		rsp = kits.APIWrapRsp(kits.ErrInnerServer, fmt.Sprintf("match all singer info with singer error: %v", err), ret)
		return
	}
	matchTracks := []*matchTrackInfo{}
	if req.Type == "mv" {
		matchTracks, err = matchOnlyMVTrack(req.Vname)
		if err != nil {
			rsp = kits.APIWrapRsp(kits.ErrInnerServer, fmt.Sprintf("match all track info with singer error: %v", err), ret)
			return
		}
	}
	matchVideos, resultStatus, err := checkMatchedSingerAndTrack(req.Type, &matchSingers, &matchTracks)
	if err != nil {
		rsp = kits.APIWrapRsp(kits.ErrInnerServer,
			fmt.Sprintf("check matched singer and track info from video error: %v", err), ret)
		return
	}
	status, _ := checkUnmatchedSingers(req.Type, &matchTracks, &matchVideos, unmatchSingers)
	if len(status) != 0 {
		resultStatus = status
	}
	ret.MatchInfos, ret.Status = matchVideos, resultStatus
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}
