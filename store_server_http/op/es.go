package op

import (
	"fmt"
	"github.com/olivere/elastic"
	elastic7 "github.com/olivere/elastic/v7"
	es "github.com/store_server/dbtools/elastic"
	es7 "github.com/store_server/dbtools/elastic7"
	"github.com/store_server/logger"
	"github.com/store_server/store_server_http/g"
	"github.com/store_server/store_server_http/kits"
)

var (
	IndexMap = map[string]string{
		"track":  "joox_music",
		"album":  "joox_music",
		"singer": "joox_music",
		"video":  "video",
	}
	TypeMap = map[string]string{
		"track":  "tracks",
		"album":  "albums",
		"singer": "singers",
		"video":  "music_mv",
	}
	NewIndexMap = map[string]string{
		"track":  "joox_tracks",
		"album":  "joox_albums",
		"singer": "joox_singers",
		"video1": "music_mv",
		"video2": "interview_mv",
	}
)

func processQuerys(terms, filter map[string]interface{}, rge map[string][2]interface{},
	query string, fields []string, multiMatch map[string][]string, should map[string]interface{},
	boosts map[string]float64, isnew bool, opts ...interface{}) ([]elastic.Query, []elastic7.Query,
	[]elastic.Query, []elastic7.Query) {
	var querys []elastic.Query
	var querys7 []elastic7.Query
	var shouldQuerys []elastic.Query
	var shouldQuerys7 []elastic7.Query
	if len(terms) > 0 {
		for k, v := range terms {
			if isnew {
				querys7 = append(querys7, es7.EsDriver.TermQuery(k, v))
			} else {
				querys = append(querys, es.EsDriver.TermQuery(k, v))
			}
		}
	}
	if len(filter) > 0 {
		for k, v := range filter {
			if isnew {
				querys7 = append(querys7, es7.EsDriver.MatchQuery(k, v))
			} else {
				querys = append(querys, es.EsDriver.MatchQuery(k, v))
			}
		}
	}
	if len(rge) > 0 {
		for k, v := range rge {
			if isnew {
				querys7 = append(querys7, es7.EsDriver.RangeQuery(k, v[0], v[1]))
			} else {
				querys = append(querys, es.EsDriver.RangeQuery(k, v[0], v[1]))
			}
		}
	}
	if len(query) > 0 {
		if isnew {
			querys7 = append(querys7, es7.EsDriver.StringQuery(query, true, fields...))
		} else {
			querys = append(querys, es.EsDriver.StringQuery(query, true, fields...))
		}
	}
	if len(multiMatch) > 0 {
		for k, v := range multiMatch {
			if isnew {
				querys7 = append(querys7, es7.EsDriver.MultiMatchQuery(k, v))
			} else {
				querys = append(querys, es.EsDriver.MultiMatchQuery(k, v))
			}
		}
	}
	if len(should) > 0 {
		for k, v := range should {
			if isnew {
				shouldQuerys7 = append(shouldQuerys7, es7.EsDriver.TermQuery(k, v))
			} else {
				shouldQuerys = append(shouldQuerys, es.EsDriver.TermQuery(k, v))
			}
		}
	}
	return querys, querys7, shouldQuerys, shouldQuerys7
}

/************************ track search相关 ***************************/
//search track request
type SearchTracksReq struct {
	Start      int                       `json:"start,omitempty"`
	Size       int                       `json:"count,omitempty"`
	Ids        []int64                   `json:"ids,omitempty"`
	Id         int64                     `json:"id,omitempty"`
	Region     *int                      `json:"region_id,omitempty"`
	Query      string                    `json:"query,omitempty"`
	Fields     []string                  `json:"fields,omitempty"`
	Terms      map[string]interface{}    `json:"terms,omitempty"`
	Filter     map[string]interface{}    `json:"filter,omitempty"`
	Should     map[string]interface{}    `json:"should,omitempty"`
	Range      map[string][2]interface{} `json:"range,omitempty"`
	MultiMatch map[string][]string       `json:"multi_match,omitempty"`
	Boosts     map[string]float64        `json:"boosts,omitempty"`
	SortBy     string                    `json:"sortby,omitempty"`
	Wildcard   bool                      `json:"wildcard,omitempty"`
	//标识是否使用新集群,下同
	New bool `json:"new,omitempty"`
	//标识是否使用泰国专用索引
	IsTh bool `json:"isth,omitempty"`
}

//search track response
type SearchTracksRsp struct {
	Tracks interface{} `json:"tracks"`
	Total  int64       `json:"total"`
}

// track search
func TracksSearch(req *SearchTracksReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.TracksSearch", &err, logger.Entry())
	ret := SearchTracksRsp{}
	var tks interface{}
	if len(req.Terms) != 0 || len(req.Filter) != 0 || len(req.MultiMatch) != 0 ||
		len(req.Range) != 0 || len(req.Query) != 0 || len(req.Should) != 0 {
		querys, querys7, shouldQuerys, shouldQuerys7 := processQuerys(req.Terms, req.Filter, req.Range, req.Query, req.Fields,
			req.MultiMatch, req.Should, req.Boosts, req.New)
		if req.New {
			ss := es7.EsDriver.SearchSource(
				es7.EsDriver.BoolQueryWithShould(querys7, shouldQuerys7), req.Start, req.Size, req.SortBy)
			if !req.IsTh {
				ret.Total, tks, err = es7.EsDriver.Search(ss, NewIndexMap["track"], "")
			} else {
				ret.Total, tks, err = es7.EsDriver.Search(ss, fmt.Sprintf("%s%s", NewIndexMap["track"], "_th"), "")
			}
		} else {
			ss := es.EsDriver.SearchSource(
				es.EsDriver.BoolQueryWithShould(querys, shouldQuerys), req.Start, req.Size, req.SortBy)
			if !req.IsTh {
				ret.Total, tks, err = es.EsDriver.Search(ss, IndexMap["track"], TypeMap["track"])
			} else {
				ret.Total, tks, err = es.EsDriver.Search(ss, fmt.Sprintf("%s%s", IndexMap["track"], "_th"), TypeMap["track"])
			}
		}
	} else if req.Id != 0 {
		if req.Region == nil {
			var ids []string
			for _, region := range g.Config().ValidRegions {
				ids = append(ids, fmt.Sprintf("track-%v-%v", region, req.Id))
			}
			if req.New {
				ret.Total, tks, err = es7.EsDriver.SearchByIds(NewIndexMap["track"], "", ids)
			} else {
				ret.Total, tks, err = es.EsDriver.SearchByIds(IndexMap["track"], TypeMap["track"], ids)
			}
		} else {
			id := fmt.Sprintf("track-%v-%v", *req.Region, req.Id)
			if req.New {
				ret.Total, tks, err = es7.EsDriver.SearchById(NewIndexMap["track"], "", id)
			} else {
				ret.Total, tks, err = es.EsDriver.SearchById(IndexMap["track"], TypeMap["track"], id)
			}
		}
	} else if len(req.Ids) != 0 {
		var ids []string
		for _, id := range req.Ids {
			if req.Region == nil {
				for _, region := range g.Config().ValidRegions {
					ids = append(ids, fmt.Sprintf("track-%v-%v", region, id))
				}
			} else {
				ids = append(ids, fmt.Sprintf("track-%v-%v", *req.Region, id))
			}
		}
		if req.New {
			ret.Total, tks, err = es7.EsDriver.SearchByIds(NewIndexMap["track"], "", ids)
		} else {
			ret.Total, tks, err = es.EsDriver.SearchByIds(IndexMap["track"], TypeMap["track"], ids)
		}
	} else {
		if len(req.Filter) == 0 && req.Start == 0 && req.Size == 0 {
			logger.Entry().Errorf("search tracks conditions is invalid")
			rsp = kits.APIWrapRsp(kits.ErrOther, "search tracks conditions is invalid", ret)
			return
		}
	}
	if err != nil {
		logger.Entry().Errorf("search tracks error: %v|request: %v", err, *req)
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	ret.Tracks = tks
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

/************************ album search相关 ***************************/
//search album request
type SearchAlbumsReq struct {
	Start      int                       `json:"start,omitempty"`
	Size       int                       `json:"count,omitempty"`
	Ids        []int64                   `json:"ids,omitempty"`
	Id         int64                     `json:"id,omitempty"`
	Region     *int                      `json:"region_id,omitempty"`
	Query      string                    `json:"query,omitempty"`
	Fields     []string                  `json:"fields,omitempty"`
	Terms      map[string]interface{}    `json:"terms,omitempty"`
	Filter     map[string]interface{}    `json:"filter,omitempty"`
	Should     map[string]interface{}    `json:"should,omitempty"`
	Range      map[string][2]interface{} `json:"range,omitmepty"`
	MultiMatch map[string][]string       `json:"multi_match,omitempty"`
	Boosts     map[string]float64        `json:"boosts,omitempty"`
	SortBy     string                    `json:"sortby,omitempty"`
	Wildcard   bool                      `json:"wildcard,omitempty"`
	New        bool                      `json:"new,omitempty"`
	//标识是否使用泰国专用索引
	IsTh bool `json:"isth,omitempty"`
}

//search album response
type SearchAlbumsRsp struct {
	Albums interface{} `json:"albums"`
	Total  int64       `json:"total"`
}

// album search
func AlbumsSearch(req *SearchAlbumsReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.AlbumsSearch", &err, logger.Entry())
	ret := SearchAlbumsRsp{}
	var ams interface{}
	if len(req.Terms) != 0 || len(req.Filter) != 0 || len(req.MultiMatch) != 0 ||
		len(req.Range) != 0 || len(req.Query) != 0 || len(req.Should) != 0 {
		querys, querys7, shouldQuerys, shouldQuerys7 := processQuerys(req.Terms, req.Filter, req.Range, req.Query, req.Fields,
			req.MultiMatch, req.Should, req.Boosts, req.New)
		if req.New {
			ss := es7.EsDriver.SearchSource(
				es7.EsDriver.BoolQueryWithShould(querys7, shouldQuerys7), req.Start, req.Size, req.SortBy)
			if !req.IsTh {
				ret.Total, ams, err = es7.EsDriver.Search(ss, NewIndexMap["album"], "")
			} else {
				ret.Total, ams, err = es7.EsDriver.Search(ss, fmt.Sprintf("%s%s", NewIndexMap["album"], "_th"), "")
			}
		} else {
			ss := es.EsDriver.SearchSource(
				es.EsDriver.BoolQueryWithShould(querys, shouldQuerys), req.Start, req.Size, req.SortBy)
			if !req.IsTh {
				ret.Total, ams, err = es.EsDriver.Search(ss, IndexMap["album"], TypeMap["album"])
			} else {
				ret.Total, ams, err = es.EsDriver.Search(ss, fmt.Sprintf("%s%s", IndexMap["album"], "_th"), TypeMap["album"])
			}
		}
	} else if req.Id != 0 {
		if req.Region == nil {
			var ids []string
			for _, region := range g.Config().ValidRegions {
				ids = append(ids, fmt.Sprintf("album-%v-%v", region, req.Id))
			}
			if req.New {
				ret.Total, ams, err = es7.EsDriver.SearchByIds(NewIndexMap["album"], "", ids)
			} else {
				ret.Total, ams, err = es.EsDriver.SearchByIds(IndexMap["album"], TypeMap["album"], ids)
			}
		} else {
			id := fmt.Sprintf("album-%v-%v", req.Region, req.Id)
			if req.New {
				ret.Total, ams, err = es7.EsDriver.SearchById(NewIndexMap["album"], "", id)
			} else {
				ret.Total, ams, err = es.EsDriver.SearchById(IndexMap["album"], TypeMap["album"], id)
			}
		}
	} else if len(req.Ids) != 0 {
		var ids []string
		for _, id := range req.Ids {
			if req.Region == nil {
				for _, region := range g.Config().ValidRegions {
					ids = append(ids, fmt.Sprintf("album-%v-%v", region, id))
				}
			} else {
				ids = append(ids, fmt.Sprintf("album-%v-%v", *req.Region, id))
			}
		}
		if req.New {
			ret.Total, ams, err = es7.EsDriver.SearchByIds(NewIndexMap["album"], "", ids)
		} else {
			ret.Total, ams, err = es.EsDriver.SearchByIds(IndexMap["album"], TypeMap["album"], ids)
		}
	} else {
		if len(req.Filter) == 0 && req.Start == 0 && req.Size == 0 {
			logger.Entry().Errorf("search albums conditions is invalid")
			rsp = kits.APIWrapRsp(kits.ErrOther, "search albums conditions is invalid", ret)
			return
		}
	}
	if err != nil {
		logger.Entry().Errorf("search albums error: %v|request: %v", err, *req)
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	ret.Albums = ams
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

/************************ singer search相关 ***************************/
//search singer request
type SearchSingersReq struct {
	Start      int                       `json:"start,omitempty"`
	Size       int                       `json:"count,omitempty"`
	Ids        []int64                   `json:"ids,omitempty"`
	Id         int64                     `json:"id,omitempty"`
	Region     *int                      `json:"region_id,omitempty"`
	Query      string                    `json:"query,omitempty"`
	Fields     []string                  `json:"fields,omitempty"`
	Terms      map[string]interface{}    `json:"terms,omitempty"`
	Filter     map[string]interface{}    `json:"filter,omitempty"`
	Should     map[string]interface{}    `json:"should,omitempty"`
	Range      map[string][2]interface{} `json:"range,omitempty"`
	MultiMatch map[string][]string       `json:"multi_match,omitempty"`
	Boosts     map[string]float64        `json:"boosts,omitempty"`
	SortBy     string                    `json:"sortby,omitempty"`
	Wildcard   bool                      `json:"wildcard,omitempty"`
	New        bool                      `json:"new,omitempty"`
	//标识是否使用泰国专用索引
	IsTh bool `json:"isth,omitempty"`
}

//search singer response
type SearchSingersRsp struct {
	Singers interface{} `json:"singers"`
	Total   int64       `json:"total"`
}

// singer search
func SingersSearch(req *SearchSingersReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.SingersSearch", &err, logger.Entry())
	ret := SearchSingersRsp{}
	var sgs interface{}
	if len(req.Terms) != 0 || len(req.Filter) != 0 || len(req.MultiMatch) != 0 ||
		len(req.Range) != 0 || len(req.Query) != 0 || len(req.Should) != 0 {
		querys, querys7, shouldQuerys, shouldQuerys7 := processQuerys(req.Terms, req.Filter, req.Range, req.Query, req.Fields,
			req.MultiMatch, req.Should, req.Boosts, req.New)
		if req.New {
			ss := es7.EsDriver.SearchSource(
				es7.EsDriver.BoolQueryWithShould(querys7, shouldQuerys7), req.Start, req.Size, req.SortBy)
			if !req.IsTh {
				ret.Total, sgs, err = es7.EsDriver.Search(ss, NewIndexMap["singer"], "")
			} else {
				ret.Total, sgs, err = es7.EsDriver.Search(ss, fmt.Sprintf("%s%s", NewIndexMap["singer"], "_th"), "")
			}
		} else {
			ss := es.EsDriver.SearchSource(
				es.EsDriver.BoolQueryWithShould(querys, shouldQuerys), req.Start, req.Size, req.SortBy)
			if !req.IsTh {
				ret.Total, sgs, err = es.EsDriver.Search(ss, IndexMap["singer"], TypeMap["singer"])
			} else {
				ret.Total, sgs, err = es.EsDriver.Search(ss, fmt.Sprintf("%s%s", IndexMap["singer"], "_th"), TypeMap["singer"])
			}
		}
	} else if req.Id != 0 {
		if req.Region == nil {
			var ids []string
			for _, region := range g.Config().ValidRegions {
				ids = append(ids, fmt.Sprintf("singer-%v-%v", region, req.Id))
			}
			if req.New {
				ret.Total, sgs, err = es7.EsDriver.SearchByIds(NewIndexMap["singer"], "", ids)
			} else {
				ret.Total, sgs, err = es.EsDriver.SearchByIds(IndexMap["singer"], TypeMap["singer"], ids)
			}
		} else {
			id := fmt.Sprintf("singer-%v-%v", *req.Region, req.Id)
			if req.New {
				ret.Total, sgs, err = es7.EsDriver.SearchById(NewIndexMap["singer"], "", id)
			} else {
				ret.Total, sgs, err = es.EsDriver.SearchById(IndexMap["singer"], TypeMap["singer"], id)
			}
		}
	} else if len(req.Ids) != 0 {
		var ids []string
		for _, id := range req.Ids {
			if req.Region == nil {
				for _, region := range g.Config().ValidRegions {
					ids = append(ids, fmt.Sprintf("singer-%v-%v", region, id))
				}
			} else {
				ids = append(ids, fmt.Sprintf("singer-%v-%v", *req.Region, id))
			}
		}
		if req.New {
			ret.Total, sgs, err = es7.EsDriver.SearchByIds(NewIndexMap["singer"], "", ids)
		} else {
			ret.Total, sgs, err = es.EsDriver.SearchByIds(IndexMap["singer"], TypeMap["singer"], ids)
		}
	} else {
		if len(req.Filter) == 0 && req.Start == 0 && req.Size == 0 {
			logger.Entry().Errorf("search singers conditions is invalid")
			rsp = kits.APIWrapRsp(kits.ErrOther, "search singers conditions is invalid", ret)
			return
		}
	}
	if err != nil {
		logger.Entry().Errorf("search singers error: %v|request: %v", err, *req)
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	ret.Singers = sgs
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

/************************ video search相关 ***************************/
//search video request
type SearchVideosReq struct {
	Start int   `json:"start,omitempty"`
	Size  int   `json:"count,omitempty"`
	Id    int64 `json:"id,omitempty"`
	//video type
	Type       int                       `json:"type,omitempty"`
	Query      string                    `json:"query,omitempty"`
	Fields     []string                  `json:"fields,omitempty"`
	Terms      map[string]interface{}    `json:"terms,omitempty"`
	Filter     map[string]interface{}    `json:"filter,omitempty"`
	Should     map[string]interface{}    `json:"should,omitempty"`
	Range      map[string][2]interface{} `json:"range,omitempty"`
	MultiMatch map[string][]string       `json:"multi_match,omitempty"`
	Boosts     map[string]float64        `json:"boosts,omitempty"`
	SortBy     string                    `json:"sortby,omitempty"`
	New        bool                      `json:"new,omitempty"`
}

//search video response
type SearchVideosRsp struct {
	Videos interface{} `json:"videos"`
	Total  int64       `json:"total"`
}

// video search
func VideosSearch(req *SearchVideosReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.VideosSearch", &err, logger.Entry())
	ret := SearchVideosRsp{}
	var vos interface{}
	index := NewIndexMap["video1"]
	vtype := TypeMap["video"]
	if req.Type != 0 {
		index = NewIndexMap["video2"]
		vtype = "interview_mv"
	}
	if len(req.Terms) != 0 || len(req.Filter) != 0 || len(req.MultiMatch) != 0 ||
		len(req.Range) != 0 || len(req.Should) != 0 {
		querys, querys7, shouldQuerys, shouldQuerys7 := processQuerys(req.Terms, req.Filter, req.Range, req.Query, req.Fields,
			req.MultiMatch, req.Should, req.Boosts, req.New)
		if req.New {
			ss := es7.EsDriver.SearchSource(
				es7.EsDriver.BoolQueryWithShould(querys7, shouldQuerys7), req.Start, req.Size, req.SortBy)
			ret.Total, vos, err = es7.EsDriver.Search(ss, index, "")
			if err != nil {
				ret.Total, vos, err = es7.EsDriver.Search(ss, fmt.Sprintf("%s%s", index, "_th"), "")
			}
		} else {
			ss := es.EsDriver.SearchSource(
				es.EsDriver.BoolQueryWithShould(querys, shouldQuerys), req.Start, req.Size, req.SortBy)
			ret.Total, vos, err = es.EsDriver.Search(ss, IndexMap["video"], vtype)
			if err != nil {
				ret.Total, vos, err = es.EsDriver.Search(ss, fmt.Sprintf("%s%s", IndexMap["video"], "_th"), vtype)
			}
		}
	} else if req.Id != 0 {
		id := fmt.Sprintf("%v", req.Id)
		if req.New {
			ret.Total, vos, err = es7.EsDriver.SearchById(index, "", id)
		} else {
			ret.Total, vos, err = es.EsDriver.SearchById(IndexMap["video"], vtype, id)
		}
	} else {
		if len(req.Filter) == 0 && len(req.Terms) == 0 && req.Start == 0 && req.Size == 0 {
			logger.Entry().Errorf("search videos conditions is invalid")
			rsp = kits.APIWrapRsp(kits.ErrOther, "search videos conditions is invalid", ret)
			return
		}
	}
	if err != nil {
		logger.Entry().Errorf("search videos error: %v|request: %v", err, *req)
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	ret.Videos = vos
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

/************************ track update or insert 相关 ***************************/
//update or insert track request
type UpsertTracksReq struct {
	Id        int64                  `json:"id,omitempty"`
	Region    int                    `json:"region_id,omitempty"`
	Filter    map[string]interface{} `json:"filter,omitempty"`
	UpsertDoc map[string]interface{} `json:"doc"`
	Sync      bool                   `json:"sync"`
	New       bool                   `json:"new,omitempty"`
}

//update or insert track response
type UpsertTracksRsp struct {
	Total int64 `json:"total"`
}

//track upsert
func TracksUpsert(req *UpsertTracksReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.TracksUpsert", &err, logger.Entry())
	ret := UpsertTracksRsp{}
	if req.Id != 0 {
		id := fmt.Sprintf("track-%v-%v", req.Region, req.Id)
		ddl := es.EsDriver.NewDocDecl(IndexMap["track"], TypeMap["track"], id, req.UpsertDoc)
		if req.Sync {
			err = es.EsDriver.UpsertOne(IndexMap["track"], TypeMap["track"], id, req.UpsertDoc)
		} else {
			es.EsDriver.AddOneToBulk(ddl)
		}
	} else {
		//TODO first search docs by filter, and then do update or insert
	}
	if err != nil {
		logger.Entry().Errorf("update or insert track error: %v|request: %v", err, *req)
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

/************************ album update or insert 相关 ***************************/
//update or insert album request
type UpsertAlbumsReq struct {
	Id        int64                  `json:"id,omitempty"`
	Region    int                    `json:"region_id,omitempty"`
	Filter    map[string]interface{} `json:"filter,omitempty"`
	UpsertDoc map[string]interface{} `json:"doc"`
	Sync      bool                   `json:"sync"`
	New       bool                   `json:"new,omitempty"`
}

//update or insert album response
type UpsertAlbumsRsp struct {
	Total int64 `json:"total"`
}

//album upsert
func AlbumsUpsert(req *UpsertAlbumsReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.AlbumsUpsert", &err, logger.Entry())
	ret := UpsertAlbumsRsp{}
	if req.Id != 0 {
		id := fmt.Sprintf("album-%v-%v", req.Region, req.Id)
		ddl := es.EsDriver.NewDocDecl(IndexMap["album"], TypeMap["album"], id, req.UpsertDoc)
		if req.Sync {
			err = es.EsDriver.UpsertOne(IndexMap["album"], TypeMap["album"], id, req.UpsertDoc)
		} else {
			es.EsDriver.AddOneToBulk(ddl)
		}
	} else {
		//TODO first search docs by filter, and then do update or insert
	}
	if err != nil {
		logger.Entry().Errorf("update or insert album error: %v|request: %v", err, *req)
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

/************************ singer update or insert 相关 ***************************/
//update or insert singer request
type UpsertSingersReq struct {
	Id        int64                  `json:"id,omitempty"`
	Region    int                    `json:"region_id,omitempty"`
	Filter    map[string]interface{} `json:"filter,omitempty"`
	UpsertDoc map[string]interface{} `json:"doc"`
	Sync      bool                   `json:"sync"`
	New       bool                   `json:"new,omitempty"`
}

//update or insert singer response
type UpsertSingersRsp struct {
	Total int64 `json:"total"`
}

//singer upsert
func SingersUpsert(req *UpsertSingersReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.SingersUpsert", &err, logger.Entry())
	ret := UpsertSingersRsp{}
	if req.Id != 0 {
		id := fmt.Sprintf("singer-%v-%v", req.Region, req.Id)
		ddl := es.EsDriver.NewDocDecl(IndexMap["singer"], TypeMap["singer"], id, req.UpsertDoc)
		if req.Sync {
			err = es.EsDriver.UpsertOne(IndexMap["singer"], TypeMap["singer"], id, req.UpsertDoc)
		} else {
			es.EsDriver.AddOneToBulk(ddl)
		}
	} else {
		//TODO first search docs by filter, and then do update or insert
	}
	if err != nil {
		logger.Entry().Errorf("update or insert singer error: %v|request: %v", err, *req)
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

/************************ video update or insert 相关 ***************************/
//update or insert video request
type UpsertVideosReq struct {
	Id        int64                  `json:"id,omitempty"`
	Filter    map[string]interface{} `json:"filter,omitempty"`
	UpsertDoc map[string]interface{} `json:"doc"`
	Sync      bool                   `json:"sync"`
	New       bool                   `json:"new,omitempty"`
}

//update or insert video response
type UpsertVideosRsp struct {
	Total int64 `json:"total"`
}

//video upsert
func VideosUpsert(req *UpsertVideosReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.VideosUpsert", &err, logger.Entry())
	ret := UpsertVideosRsp{}
	if req.Id != 0 {
		id := fmt.Sprintf("%v", req.Id)
		ddl := es.EsDriver.NewDocDecl(IndexMap["video"], TypeMap["video"], id, req.UpsertDoc)
		if req.Sync {
			err = es.EsDriver.UpsertOne(IndexMap["video"], TypeMap["video"], id, req.UpsertDoc)
		} else {
			es.EsDriver.AddOneToBulk(ddl)
		}
	} else {
		//TODO first search docs by filter, and then do update or insert
	}
	if err != nil {
		logger.Entry().Errorf("update or insert video error: %v|request: %v", err, *req)
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

/************************ track delete相关 ***************************/
//delete track request
type DeleteTrackDocReq struct {
	Id     int64                  `json:"id,omitempty"`
	Region int                    `json:"region_id,omitempty"`
	Filter map[string]interface{} `json:"filter,omitempty"`
	Sync   bool                   `json:"sync"`
	New    bool                   `json:"new,omitempty"`
}

//delete track response
type DeleteTrackDocRsp struct {
	Total int64 `json:"total"`
}

//track delete
func TrackDocDelete(req *DeleteTrackDocReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.TrackDocDelete", &err, logger.Entry())
	ret := DeleteTrackDocRsp{}
	if req.Id != 0 {
		id := fmt.Sprintf("track-%v-%v", req.Region, req.Id)
		ddl := es.EsDriver.NewDocDecl(IndexMap["track"], TypeMap["track"], id, nil, true)
		if req.Sync {
			err = es.EsDriver.DeleteOne(IndexMap["track"], TypeMap["track"], id)
		} else {
			es.EsDriver.AddOneToBulk(ddl)
		}
	} else {
		//TODO first search docs by filter, and then do delete
	}
	if err != nil {
		logger.Entry().Errorf("delete track doc error: %v|request: %v", err, *req)
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

/************************ album delete相关 ***************************/
//delete album request
type DeleteAlbumDocReq struct {
	Id     int64                  `json:"id,omitempty"`
	Region int                    `json:"region_id,omitempty"`
	Filter map[string]interface{} `json:"filter,omitempty"`
	Sync   bool                   `json:"sync"`
	New    bool                   `json:"new,omitempty"`
}

//delete album response
type DeleteAlbumDocRsp struct {
	Total int64 `json:"total"`
}

//album delete
func AlbumDocDelete(req *DeleteAlbumDocReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.AlbumDocDelete", &err, logger.Entry())
	ret := DeleteAlbumDocRsp{}
	if req.Id != 0 {
		id := fmt.Sprintf("album-%v-%v", req.Region, req.Id)
		ddl := es.EsDriver.NewDocDecl(IndexMap["album"], TypeMap["album"], id, nil, true)
		if req.Sync {
			err = es.EsDriver.DeleteOne(IndexMap["album"], TypeMap["album"], id)
		} else {
			es.EsDriver.AddOneToBulk(ddl)
		}
	} else {
		//TODO first search docs by filter, and then do delete
	}
	if err != nil {
		logger.Entry().Errorf("delete album doc error: %v|request: %v", err, *req)
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

/************************ singer delete相关 ***************************/
//delete singer request
type DeleteSingerDocReq struct {
	Id     int64                  `json:"id,omitempty"`
	Region int                    `json:"region_id,omitempty"`
	Filter map[string]interface{} `json:"filter,omitempty"`
	Sync   bool                   `json:"sync"`
	New    bool                   `json:"new,omitempty"`
}

//delete singer response
type DeleteSingerDocRsp struct {
	Total int64 `json:"total"`
}

//singer delete
func SingerDocDelete(req *DeleteSingerDocReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.SingerDocDelete", &err, logger.Entry())
	ret := DeleteSingerDocRsp{}
	if req.Id != 0 {
		id := fmt.Sprintf("singer-%v-%v", req.Region, req.Id)
		ddl := es.EsDriver.NewDocDecl(IndexMap["singer"], TypeMap["singer"], id, nil, true)
		if req.Sync {
			err = es.EsDriver.DeleteOne(IndexMap["singer"], TypeMap["singer"], id)
		} else {
			es.EsDriver.AddOneToBulk(ddl)
		}
	} else {
		//TODO first search docs by filter, and then do delete
	}
	if err != nil {
		logger.Entry().Errorf("delete singer doc error: %v|request: %v", err, *req)
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}

/************************ video delete相关 ***************************/
//delete video request
type DeleteVideoDocReq struct {
	Id     int64                  `json:"id,omitempty"`
	Filter map[string]interface{} `json:"filter,omitempty"`
	Sync   bool                   `json:"sync"`
	New    bool                   `json:"new,omitempty"`
}

//delete video response
type DeleteVideoDocRsp struct {
	Total int64 `json:"total"`
}

//video delete
func VideoDocDelete(req *DeleteVideoDocReq) (rsp *kits.WrapRsp, err error) {
	defer kits.CatchErr("http.VideoDocDelete", &err, logger.Entry())
	ret := DeleteVideoDocRsp{}
	if req.Id != 0 {
		id := fmt.Sprintf("%v", req.Id)
		ddl := es.EsDriver.NewDocDecl(IndexMap["video"], TypeMap["video"], id, nil, true)
		if req.Sync {
			err = es.EsDriver.DeleteOne(IndexMap["video"], TypeMap["video"], id)
		} else {
			es.EsDriver.AddOneToBulk(ddl)
		}
	} else {
		//TODO first search docs by filter, and then do delete
	}
	if err != nil {
		logger.Entry().Errorf("delete video doc error: %v|request: %v", err, *req)
		rsp = kits.APIWrapRsp(kits.ErrOther, err.Error(), ret)
		return
	}
	rsp = kits.APIWrapRsp(0, "ok", ret)
	return
}
