package elastic7

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/store_server/logger"
	"github.com/store_server/utils/common"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	bulkSize = 500
)

var (
	EsDriver  *ESClient
	sortOrder = true
	ALLID     = regexp.MustCompile(`^(\d+)`)
	SUBS      = regexp.MustCompile(`([;|,\n\t]|\s{2,})`)
)

//es响应数据结构
type EsSearchRsp struct {
	Took    float64     `json:"took"`
	TimeOut bool        `json:"time_out"`
	Shard   interface{} `json:"_shards"`
	Hits    EsOuterHits `json:"hits"`
}

//es outer hits
type EsOuterHits struct {
	Total    int            `json:"total"`
	MaxScore float64        `json:"max_score"`
	Hits     []*EsInnerHits `json:"hits"`
}

//es shard
type EsShard struct {
	Total   int `json:"total"`
	Success int `json:"successful"`
	Failed  int `json:"failed"`
}

//es inner hits
type EsInnerHits struct {
	Index  string      `json:"_index"`
	Type   string      `json:"_type"`
	Id     string      `json:"_id"`
	Score  float64     `json:"_score"`
	Source interface{} `json:"_source,omitempty"`
}

//es client definition
type ESClient struct {
	ctx           context.Context
	cancel        context.CancelFunc
	client        *elastic.Client
	bulkService   *elastic.BulkService
	getService    *elastic.GetService
	scrollService *elastic.ScrollService
	bulkCh        chan *DocDecl
	lock          sync.RWMutex
	closed        bool

	index   string
	docType string
}

//doc declaration
type DocDecl struct {
	Index  string
	Type   string
	Id     string
	Doc    interface{}
	Delete bool
}

//es client propertion definition
func (c *ESClient) SetIndex(index string) *ESClient {
	c.index = index
	return c
}

func (c *ESClient) SetType(docType string) *ESClient {
	if len(docType) == 0 {
		docType = "_doc"
	}
	c.docType = docType
	return c
}

//_type参数为兼容7.x以下版本数据
func (c *ESClient) checkType(tptr *string) {
	if tptr == nil {
		return
	}
	if len(*tptr) == 0 {
		*tptr = "_doc"
	}
}

func (c *ESClient) Client() *elastic.Client {
	return c.client
}

func (c *ESClient) NewDocDecl(index, _type, id string, doc interface{}, args ...bool) *DocDecl {
	deleted := false
	if len(args) > 0 {
		deleted = args[0]
	}
	c.checkType(&_type)
	return &DocDecl{index, _type, id, doc, deleted}
}

//校验查询条件是否全部为ID(支持多ID查询)
func (c *ESClient) checkIsAllID(condition string) (query string, isAll bool) {
	isAll = true
	query = SUBS.ReplaceAllString(condition, " ")
	strs := strings.Split(query, " ")
	tmp := make([]string, 0, len(strs))
	for _, str := range strs {
		if len(str) <= 0 {
			continue
		}
		if !ALLID.MatchString(str) {
			isAll = false
		}
		t := "\"" + str + "\""
		tmp = append(tmp, t)
	}
	//全文检索匹配
	query = strings.Join(tmp, " ")
	return
}

//字符串查询, opts represent fields to run query string against
func (c *ESClient) StringQuery(query string, isAll bool, opts ...string) *elastic.QueryStringQuery {
	strQuery := elastic.NewQueryStringQuery(query).TieBreaker(0.3).MinimumShouldMatch("95%")
	var defaultOpt string = "and"
	if isAll {
		defaultOpt = "or"
	}
	for _, opt := range opts {
		strQuery = strQuery.Field(opt)
	}
	return strQuery.DefaultOperator(defaultOpt).AnalyzeWildcard(true).Escape(true)
}

//term查询
func (c *ESClient) TermQuery(field string, val interface{}, opts ...interface{}) *elastic.TermQuery {
	termQuery := elastic.NewTermQuery(field, val)
	if len(opts) > 0 {
		termQuery = termQuery.Boost(opts[0].(float64))
	}
	return termQuery
}

//match查询
func (c *ESClient) MatchQuery(field string, val interface{}, opts ...interface{}) *elastic.MatchQuery {
	matchQuery := elastic.NewMatchQuery(field, val)
	if len(opts) > 0 {
		matchQuery = matchQuery.Boost(opts[0].(float64))
	}
	return matchQuery.MinimumShouldMatch("95%").Lenient(true)
}

//range查询
func (c *ESClient) RangeQuery(field string, lower, upper interface{}) *elastic.RangeQuery {
	rangeQuery := elastic.NewRangeQuery(field).Gte(lower).Lte(upper)
	return rangeQuery
}

//multi-match查询
func (c *ESClient) MultiMatchQuery(val interface{}, fields []string,
	boosts ...map[string]float64) *elastic.MultiMatchQuery {
	multiMatchQuery := elastic.NewMultiMatchQuery(val, fields...)
	if boosts != nil && len(boosts) > 0 {
		for _, boost := range boosts {
			for field, bst := range boost {
				multiMatchQuery = multiMatchQuery.FieldWithBoost(field, bst)
			}
		}
	}
	return multiMatchQuery.TieBreaker(0.3).MinimumShouldMatch("95%")
}

//bool查询
func (c *ESClient) BoolQuery() *elastic.BoolQuery {
	//if no need of doc score, can set query to filter of bool query
	boolQuery := elastic.NewBoolQuery()
	return boolQuery
}

//bool查询with should
func (c *ESClient) BoolQueryWithShould(mustQuery []elastic.Query,
	shouldQuery []elastic.Query) *elastic.BoolQuery {
	boolQuery := c.BoolQuery().Must(mustQuery...)
	if len(shouldQuery) > 0 {
		boolQuery = boolQuery.Should(shouldQuery...).MinimumNumberShouldMatch(1)
	}
	return boolQuery
}

//terms查询
func (c *ESClient) TermsQuery(field string, vals ...interface{}) *elastic.TermsQuery {
	termsQuery := elastic.NewTermsQuery(field, vals...)
	return termsQuery
}

//wildcard查询
func (c *ESClient) WildcardQuery(field string, wildcard string) *elastic.WildcardQuery {
	wildcardQuery := elastic.NewWildcardQuery(field, wildcard)
	return wildcardQuery
}

//new search source
func (c *ESClient) SearchSource(query elastic.Query, opts ...interface{}) *elastic.SearchSource {
	from, size := 0, 50
	var sortBy string
	var after interface{}
	if len(opts) > 0 {
		i, _ := common.Interface2Int(opts[0])
		if i > 0 {
			from = i
		}
	}
	if len(opts) > 1 {
		i, _ := common.Interface2Int(opts[1])
		if i > 0 {
			size = i
		}
	}
	if len(opts) > 2 {
		sortBy = opts[2].(string)
	}
	if len(opts) > 3 {
		after = opts[3]
	}
	var ss *elastic.SearchSource
	if len(sortBy) != 0 {
		ss = elastic.NewSearchSource().Query(query).Sort(sortBy, !sortOrder).From(from).Size(size)
	} else {
		ss = elastic.NewSearchSource().Query(query).From(from).Size(size)
	}
	if after != nil {
		ss.SearchAfter(after)
	}
	return ss.TrackTotalHits(true)
}

// query with scroll service, to deal with deep paging problem
func (c *ESClient) SearchByScroll(query elastic.Query, opts ...interface{}) (total int64,
	docs []*json.RawMessage, scrollId string, err error) {
	size := 50
	var sortBy, rscrollId string
	if len(opts) > 0 {
		i, _ := common.Interface2Int(opts[0])
		if i > 0 {
			size = i
		}
	}
	if len(opts) > 1 {
		sortBy = opts[1].(string)
	}
	if len(opts) > 2 {
		rscrollId = opts[2].(string)
	}
	c.scrollService.Query(query).Size(size).KeepAlive("5m")
	if len(sortBy) != 0 {
		c.scrollService.Sort(sortBy, !sortOrder)
	}
	if len(rscrollId) != 0 {
		c.scrollService.ScrollId(rscrollId)
	}
	res, err := c.scrollService.Do(c.ctx)
	if err != nil {
		return 0, nil, "", err
	}
	if res == nil || res.Hits == nil || res.Hits.TotalHits == nil {
		err = fmt.Errorf("search result is nil")
		return 0, nil, "", err
	}
	total = res.Hits.TotalHits.Value
	docs = make([]*json.RawMessage, 0, len(res.Hits.Hits))
	for _, item := range res.Hits.Hits {
		docs = append(docs, &item.Source)
	}
	if len(res.Hits.Hits) <= 0 {
		c.ClearScrollService(res.ScrollId)
	}
	return total, docs, res.ScrollId, nil
}

//ID搜索
func (c *ESClient) SearchById(index, _type, id string) (int64, []*json.RawMessage, error) {
	c.checkType(&_type)
	if c == nil {
		return 0, nil, fmt.Errorf("invalid es client")
	}
	if len(id) <= 0 {
		return 0, nil, fmt.Errorf("invalid doc id")
	}
	//if mapping set store fields, can specify store fields by use StoredFields for getService
	res, err := c.getService.Index(index).Type(_type).Id(id).ErrorTrace(true).Human(true).Do(c.ctx)
	if err != nil {
		logger.Entry().Errorf("search by id[%v] error: %v", id, err)
		return 0, nil, err
	}
	if res == nil {
		return 0, nil, fmt.Errorf("invalid response is nil.")
	}
	if !res.Found {
		return 0, nil, fmt.Errorf("doc not found by id: %v", id)
	}
	return 1, []*json.RawMessage{&res.Source}, nil
}

//ID数组搜索
func (c *ESClient) SearchByIds(index, _type string, ids []string) (int64, []*json.RawMessage, error) {
	c.checkType(&_type)
	if c == nil {
		return 0, nil, fmt.Errorf("invalid es client")
	}
	if len(ids) <= 0 {
		return 0, nil, fmt.Errorf("invalid doc id")
	}
	mgetService := elastic.NewMgetService(c.client)
	items := make([]*elastic.MultiGetItem, 0, len(ids))
	for _, id := range ids {
		item := elastic.NewMultiGetItem()
		item.Index(index).Type(_type).Id(id)
		items = append(items, item)
	}
	res, err := mgetService.Add(items...).ErrorTrace(true).Human(true).Do(c.ctx)
	if err != nil {
		logger.Entry().Errorf("search by ids[%v] error: %v", ids, err)
		return 0, nil, err
	}
	if res == nil {
		return 0, nil, fmt.Errorf("invalid response is nil.")
	}
	docs := make([]*json.RawMessage, 0, len(res.Docs))
	for _, doc := range res.Docs {
		if doc.Found {
			docs = append(docs, &doc.Source)
		}
	}
	if len(docs) == 0 {
		return 0, nil, fmt.Errorf("doc not found by ids: %v", ids)
	}
	return int64(len(docs)), docs, nil
}

//string query builder
func (c *ESClient) SearchByQueryString(query string, opts ...interface{}) *elastic.SearchSource {
	if len(query) <= 0 {
		//无搜索条件则全量匹配
		query = "*"
	}
	q, isAll := c.checkIsAllID(query)
	if query == "*" {
		isAll = true
	}
	//TODO extract field or prefix from query
	var fields []string
	if len(opts) > 3 {
		if t, ok := opts[3].([]string); ok {
			fields = t
		}
	}
	strQuery := c.StringQuery(q, isAll, fields...)
	return c.SearchSource(strQuery, opts[:3]...)
}

//term query builder
func (c *ESClient) SearchByTermQuery(field string, val interface{}) *elastic.SearchSource {
	termQuery := c.TermQuery(field, val)
	return c.SearchSource(termQuery)
}

//match query builder
func (c *ESClient) SearchByMatchQuery(field string, val interface{}, opts ...interface{}) *elastic.SearchSource {
	fuzziness := false
	if len(opts) > 0 {
		fuzziness = true
	}
	matchQuery := c.MatchQuery(field, val)
	if fuzziness {
		matchQuery = matchQuery.Fuzziness("AUTO")
	}
	return c.SearchSource(matchQuery, opts...)
}

//multi match query builder
func (c *ESClient) SearchByMultiMatchQuery(fields []string, val interface{},
	boosts []map[string]float64, opts ...interface{}) *elastic.SearchSource {
	matchQuery := c.MultiMatchQuery(val, fields, boosts...)
	fuzziness := false
	if len(opts) > 0 {
		fuzziness = true
	}
	if fuzziness {
		matchQuery = matchQuery.Fuzziness("AUTO")
	}
	return c.SearchSource(matchQuery, opts...)
}

//wildcard query builder
func (c *ESClient) SearchByWildcardQuery(field, wildcard string, opts ...interface{}) *elastic.SearchSource {
	if len(wildcard) <= 0 {
		wildcard = "*"
	}
	wildcardQuery := c.WildcardQuery(field, wildcard)
	return c.SearchSource(wildcardQuery, opts...)
}

//must条件搜索(key-value对) builder
func (c *ESClient) SearchByFilter(filter map[string]interface{}, opts ...interface{}) (*elastic.SearchSource, error) {
	if len(filter) <= 0 {
		return nil, fmt.Errorf("invalid search filter")
	}
	boolQuery := c.BoolQuery()
	query := make([]elastic.Query, 0, len(filter))
	for k, v := range filter {
		matchQuery := c.MatchQuery(k, v).Operator("AND")
		query = append(query, matchQuery)
	}
	boolQuery = boolQuery.Must(query...)
	return c.SearchSource(boolQuery, opts...), nil
}

//must搜索 builder
func (c *ESClient) SearchByMust(query ...elastic.Query) (*elastic.SearchSource, error) {
	if len(query) <= 0 {
		return nil, fmt.Errorf("invalid must query")
	}
	boolQuery := c.BoolQuery().Must(query...)
	return c.SearchSource(boolQuery), nil
}

//must not搜索 builder
func (c *ESClient) SearchByMustNot(query ...elastic.Query) (*elastic.SearchSource, error) {
	if len(query) <= 0 {
		return nil, fmt.Errorf("invalid must not query")
	}
	boolQuery := c.BoolQuery().MustNot(query...)
	return c.SearchSource(boolQuery), nil
}

//should搜索 builder, query can be match, term, range, bool, etc.
func (c *ESClient) SearchByShould(query ...elastic.Query) (*elastic.SearchSource, error) {
	if len(query) <= 0 {
		return nil, fmt.Errorf("invalid should query")
	}
	boolQuery := c.BoolQuery().Should(query...)
	return c.SearchSource(boolQuery), nil
}

//bool搜索 builder
func (c *ESClient) SearchByBool(mustQuery []elastic.Query,
	shouldQuery []elastic.Query) (*elastic.SearchSource, error) {
	boolQuery := c.BoolQuery()
	if len(mustQuery) > 0 {
		boolQuery = boolQuery.Must(mustQuery...)
	}
	if len(shouldQuery) > 0 {
		boolQuery = boolQuery.Should(shouldQuery...).MinimumNumberShouldMatch(1)
	}
	return c.SearchSource(boolQuery), nil
}

//range条件搜索 builder
func (c *ESClient) SearchByRange(field string, lowerLimit, upperLimit interface{},
	opts ...bool) *elastic.SearchSource {
	isNumber := false
	if len(opts) > 0 {
		isNumber = opts[0]
	}
	rangeQuery := elastic.NewRangeQuery(field)
	if isNumber {
		rangeQuery = rangeQuery.Gt(lowerLimit).Lt(upperLimit)
	} else {
		rangeQuery = rangeQuery.From(lowerLimit).To(upperLimit)
	}
	return c.SearchSource(rangeQuery)
}

//search builder do search action
func (c *ESClient) Search(ss *elastic.SearchSource, index, _type string, opts ...interface{}) (total int64,
	docs []*json.RawMessage, err error) {
	ss = ss.TrackTotalHits(true)
	if len(opts) > 0 {
		var sortBy string
		from, size := 0, 50
		i, _ := common.Interface2Int(opts[0])
		if i > 0 {
			from = i
		}
		if len(opts) > 1 {
			i, _ := common.Interface2Int(opts[1])
			if i > 0 {
				size = i
			}
		}
		if len(opts) > 2 {
			sortBy = opts[2].(string)
		}
		ss = ss.From(from).Size(size)
		if len(sortBy) != 0 {
			ss = ss.Sort(sortBy, false)
		}
	}
	//for debug
	source, _ := ss.Source()
	logger.Entry().Debugf("search source: %v", source)

	c.checkType(&_type)
	if c == nil {
		return 0, nil, fmt.Errorf("invalid es client")
	}
	if ss == nil {
		err = fmt.Errorf("invalid search source")
		return 0, nil, err
	}
	res, err := c.client.Search(index).Type(_type).SearchSource(ss).ErrorTrace(true).Human(true).Do(c.ctx)
	if err != nil {
		return 0, nil, err
	}
	if res == nil || res.Hits == nil || res.Hits.TotalHits == nil {
		err = fmt.Errorf("search result is nil")
		return 0, nil, err
	}
	total = res.Hits.TotalHits.Value
	docs = make([]*json.RawMessage, 0, total)
	for _, item := range res.Hits.Hits {
		docs = append(docs, &item.Source)
	}
	return
}

//upsert builder do update or index action
func (c *ESClient) Upsert(bks *elastic.BulkService, index, _type string) (total int64, err error) {
	if c == nil {
		return 0, fmt.Errorf("invalid es client")
	}
	if bks == nil {
		err = fmt.Errorf("invalid bulk request service")
		return 0, err
	}
	brs, err := bks.Timeout("5m").ErrorTrace(true).Do(c.ctx)
	if err != nil {
		time.Sleep(time.Second)
		// retry once
		brs, err = bks.Timeout("5m").ErrorTrace(true).Do(c.ctx)
	}
	if err != nil {
		logger.Entry().Errorf("es client do bulk request error: %v", err)
		return 0, err
	}
	if brs != nil {
		total = int64(len(brs.Items))
	}
	return
}

//单个写入(更新)
func (c *ESClient) UpsertOne(index, _type, id string, doc interface{}) (err error) {
	//metrics.ElasticOpCounter.WithLabelValues(counterLabel, "write_request", "total").Inc()
	if c == nil {
		return fmt.Errorf("invalid es client")
	}
	c.checkType(&_type)
	for i := 0; i < 2; i++ {
		_, err = c.client.Update().Index(index).Type(_type).Id(id).Doc(doc).DocAsUpsert(true).ErrorTrace(true).Do(c.ctx)
		if err != nil {
			time.Sleep(time.Duration(i+1) ^ 2*time.Second)
			continue
		} else {
			break
		}
	}
	if err != nil {
		//metrics.ElasticOpCounter.WithLabelValues(counterLabel, "write_request", "failed").Inc()
		logger.Entry().Errorf("es client do write one doc error: %v", err)
		return
	}
	//metrics.ElasticOpCounter.WithLabelValues(counterLabel, "write_request", "success").Inc()
	return
}

//单个删除
func (c *ESClient) DeleteOne(index, _type, id string) error {
	if c == nil {
		return fmt.Errorf("invalid es client")
	}
	if len(id) <= 0 {
		return fmt.Errorf("invalid id specified by delete")
	}
	deleteService := elastic.NewDeleteService(c.client)
	res, err := deleteService.Index(index).Type(_type).Id(id).Do(c.ctx)
	if err != nil {
		logger.Entry().Errorf("es client do delete one doc error: %v", err)
		return err
	}
	if res.Status == 404 {
		err = fmt.Errorf("doc[%v] not found", id)
	}
	return nil
}

//添加单个更新文档到批处理
func (c *ESClient) addUpdateOneToBulk(index, _type, id string, doc interface{}) error {
	bulkUpdateReq := elastic.NewBulkUpdateRequest()
	bulkUpdateReq.Index(index).Type(_type).Id(id).Doc(doc).DocAsUpsert(true)
	c.bulkService.Add(bulkUpdateReq)
	return nil
}

//添加单个删除文档到批处理
func (c *ESClient) addDeleteOneToBulk(index, _type, id string) error {
	bulkDeleteReq := elastic.NewBulkDeleteRequest()
	bulkDeleteReq.Index(index).Type(_type).Id(id)
	c.bulkService.Add(bulkDeleteReq)
	return nil
}

//添加单个文档到批处理
func (c *ESClient) AddOneToBulk(doc *DocDecl) {
	if doc == nil {
		return
	}
	if doc.Delete {
		c.addDeleteOneToBulk(doc.Index, doc.Type, doc.Id)
	} else {
		c.addUpdateOneToBulk(doc.Index, doc.Type, doc.Id, doc.Doc)
	}
}

//将剩余文档添加到批处理
func (c *ESClient) addRestToBulk() {
	cnt := 0
	if len(c.bulkCh) > 0 {
		for doc := range c.bulkCh {
			if doc == nil {
				continue
			}
			cnt++
			c.AddOneToBulk(doc)
			if cnt >= bulkSize {
				c.doBulkOperation()
			}
			if c.closed {
				break
			}
		}
	}
}

//批量写入，由接口主导
func (c *ESClient) BulkWrite(sources []*DocDecl) error {
	if c == nil {
		return fmt.Errorf("invalid es client")
	}
	bulkService := elastic.NewBulkService(c.client)
	for _, source := range sources {
		if source.Delete {
			bulkDeleteReq := elastic.NewBulkDeleteRequest()
			bulkDeleteReq.Index(source.Index).Type(source.Type).Id(source.Id)
			bulkService.Add(bulkDeleteReq)
		} else {
			bulkUpdateReq := elastic.NewBulkUpdateRequest()
			bulkUpdateReq.Index(source.Index).Type(source.Type).Id(source.Id).Doc(source.Doc).DocAsUpsert(true)
			bulkService.Add(bulkUpdateReq)
		}
	}
	_, err := bulkService.Timeout("5m").ErrorTrace(true).Do(c.ctx)
	if err != nil {
		time.Sleep(time.Second)
		_, err = bulkService.Timeout("5m").ErrorTrace(true).Do(c.ctx)
	}
	if err != nil {
		logger.Entry().Errorf("es client do bulk write request for api error: %v", err)
		return err
	}
	return nil
}

//定时批量写入，由进程主导
func (c *ESClient) doBulkOperation() error {
	defer c.bulkService.Reset()
	docCnt := float64(c.bulkService.NumberOfActions())
	if docCnt <= 0 {
		return nil
	}
	//metrics.ElasticOpCounter.WithLabelValues(counterLabel, "write_request", "total").Add(docCnt)
	var err error
	for i := 0; i < 2; i++ {
		_, err = c.bulkService.Timeout("5m").ErrorTrace(true).Do(c.ctx)
		if err != nil {
			time.Sleep(time.Duration(i+1) ^ 2*time.Second)
			continue
		} else {
			break
		}
	}
	if err != nil {
		//metrics.ElasticOpCounter.WithLabelValues(counterLabel, "write_request", "failed").Add(docCnt)
		logger.Entry().Errorf("es client do bulk write request error: %v", err)
		return err
	}
	//metrics.ElasticOpCounter.WithLabelValues(counterLabel, "write_request", "success").Add(docCnt)
	return nil
}

//定时执行bulk操作
func (c *ESClient) BulkOpTimely() {
	tk := time.NewTicker(10 * time.Second)
	defer tk.Stop()
	count := 0
	for {
		select {
		case <-c.ctx.Done():
			//退出前处理剩余doc
			c.addRestToBulk()
			c.doBulkOperation()
			return
		case doc := <-c.bulkCh:
			//满足批次条件
			c.AddOneToBulk(doc)
			count++
			if count%bulkSize == 0 {
				count = 0
				c.doBulkOperation()
			}
		case <-tk.C:
			//达到指定时间间隔强制写入
			c.doBulkOperation()
		}
	}
}

func (c *ESClient) Run() {
	if c == nil {
		return
	}
	go c.BulkOpTimely()
}

//清除scroll service, 删除游标释放内存
func (c *ESClient) ClearScrollService(scrollIds ...string) error {
	var err error
	if c.client != nil {
		_, err = c.client.ClearScroll().ScrollId(scrollIds...).Do(c.ctx)
	}
	return err
}

func (c *ESClient) Close() {
	if c == nil {
		return
	}
	if c.client != nil {
		c.client.Stop()
	}
	if c.bulkCh != nil {
		close(c.bulkCh)
	}
	c.closed = true
}

//new es client with options
func NewClient(addrs []string, timeout int, sniff bool, proxyAddr string, args ...string) (*elastic.Client, error) {
	options := []elastic.ClientOptionFunc{
		elastic.SetURL(addrs...),
		elastic.SetSniff(sniff),
		elastic.SetMaxRetries(3),
		elastic.SetHealthcheckInterval(10 * time.Minute),
		elastic.SetHealthcheckTimeoutStartup(30 * time.Second),
	}

	if len(args) != 0 {
		userName, passwd := args[0], ""
		if len(args) == 2 {
			passwd = args[1]
		}
		options = append(options, elastic.SetBasicAuth(userName, passwd))
	}

	if timeout == 0 {
		timeout = 15000
	}

	proxy := func(_ *http.Request) (*url.URL, error) {
		if len(proxyAddr) != 0 {
			return url.Parse(fmt.Sprintf("http://%s", proxyAddr))
		}
		return nil, nil
	}
	options = append(
		options, elastic.SetHttpClient(
			&http.Client{
				Timeout: time.Duration(timeout) * time.Millisecond,
				Transport: &http.Transport{
					Proxy: proxy,
					DialContext: (&net.Dialer{
						Timeout:   60 * time.Second,
						KeepAlive: 60 * time.Second,
					}).DialContext,
					MaxIdleConns:        100,
					MaxIdleConnsPerHost: 100,
					IdleConnTimeout:     90 * time.Second,
				},
			},
		),
	)
	return elastic.NewClient(options...)
}

func NewEsClient(
	ctx context.Context, addrs []string, timeout int, sniff bool, proxyAddr string, args ...string,
) (*ESClient, error) {
	client, err := NewClient(addrs, timeout, sniff, proxyAddr, args...)
	if err != nil {
		return nil, err
	}
	bulkService := elastic.NewBulkService(client)
	getService := elastic.NewGetService(client)
	scrollService := elastic.NewScrollService(client)
	return &ESClient{
		client:        client,
		bulkService:   bulkService,
		getService:    getService,
		scrollService: scrollService,
		bulkCh:        make(chan *DocDecl, bulkSize),
		ctx:           ctx,
	}, nil
}
