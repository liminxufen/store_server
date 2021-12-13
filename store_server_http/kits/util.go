package kits

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	es "github.com/store_server/dbtools/elastic"
	"github.com/store_server/logger"
	"github.com/store_server/metrics"
	"github.com/store_server/store_server_http/g"
)

var (
	IPWhiteLst []string
	//InfluxClient *monitor.InfluxDriver
	QPS         []CountQPS //保存qps统计数据
	HTTPCounter *CounterService
)

/*func InitInfluxEnv(ctx context.Context) (err error) {
	InfluxClient, err = monitor.NewInfluxDriver(ctx, fmt.Sprintf("http://%s:%v",
		g.Config().Influx.Host, g.Config().Influx.Port), g.Config().Influx.Database,
		g.Config().Influx.Measurement)
	if err == nil && InfluxClient == nil {
		err = fmt.Errorf("InfluxClient is Nil...")
	}
	return
}*/

//请求鉴权部分
func CheckIpHasPermission(ip string) (ok bool) {
	ok = false
	if len(ip) == 0 {
		logger.Entry().Error("access ip is empty...")
		return
	}
	for _, t := range IPWhiteLst { //check access ip is in white list
		if ip == t {
			ok = true
			break
		}
	}
	return
}

//请求部分共用
type CountQPS struct {
	CountPerSecond int //每秒请求数
	Timestamp      int64
}

//counter service
type CounterService struct {
	CountQPS
	CountAll int64 //请求总数
	Lock     sync.Mutex
}

func NewCounterService(ctx context.Context) *CounterService { //qps计数服务
	counter := &CounterService{}
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				counter.Lock.Lock()
				counter.Timestamp = time.Now().Unix()
				if counter.CountPerSecond > 0 {
					QPS = append(QPS, CountQPS{counter.CountPerSecond, counter.Timestamp})
					if len(QPS) > 10000 { //仅保留10000个采样点
						QPS = QPS[:10000]
					}
				}
				counter.CountPerSecond = 0
				counter.Lock.Unlock()
			case <-ctx.Done():
				return
			}
		}
	}()
	return counter
}

func (counter *CounterService) Increase() {
	counter.Lock.Lock()
	defer counter.Lock.Unlock()
	counter.CountAll++
	counter.CountPerSecond++
}

func QpsHandler(c *gin.Context) {
	/*defer c.Next()
	url := c.Request.URL.RequestURI()
	if strings.HasPrefix(url, "/metrics") { //监控指标不计入请求计数
		return
	}
	switch c.Request.Method {
	case "get", "Get", "GET":
		metrics.RequestMethodCounter.WithLabelValues("store_server", "GET").Inc()
	case "post", "Post", "POST":
		metrics.RequestMethodCounter.WithLabelValues("store_server", "POST").Inc()
	default:
		metrics.RequestMethodCounter.WithLabelValues("store_server", "Others").Inc()
	}*/
	HTTPCounter.Increase()
	//metrics.RequestTotalCounter.WithLabelValues("store_server").Inc()
}

func HttpQPSCounter() gin.HandlerFunc { //qps计数中间件
	return QpsHandler
}

func ApiLog() gin.HandlerFunc { //日志中间件
	return func(c *gin.Context) {
		loginName := ""

		logger.Entry().Infof("url=%s, loginName=%s, remoteAddr=%s", c.Request.URL.RequestURI(),
			loginName, getRequestAddress(c.Request))
		c.Next()
	}
}

func ApiAccessAuth() gin.HandlerFunc { //鉴权中间件
	return func(c *gin.Context) {
		ip := strings.TrimSpace(c.ClientIP())
		if !CheckIpHasPermission(ip) {
			c.Abort()
			ret := map[string]string{
				"access ip":               ip,
				"permission check result": "权限验证失败,请确认访问机器是否已注册访问权限,具体可咨询xxx",
			}
			rsp := APIWrapRsp(ErrCustom, "访问未授权", ret)
			c.JSON(http.StatusUnauthorized, rsp)
			return
		}
		c.Next()
	}
}

func getRequestAddress(req *http.Request) string { //从请求中获取客户端ip
	address := ""
	forwardedfor := req.Header.Get("X-Forwarded-For")
	if forwardedfor != "" {
		parts := strings.Split(forwardedfor, ",")
		if len(parts) >= 1 {
			address = parts[0]
		}
	}
	if address == "" {
		address = req.RemoteAddr
		i := strings.LastIndex(address, ":")
		if i != -1 {
			address = address[:i]
		}
	}
	return address
}

func checkReqParams(r *http.Request, key string) (bool, string) { //校验请求参数
	value := r.FormValue(key)
	if value != "" {
		return true, value
	}
	return false, value
}

func StatsRequestCost(ch chan struct{}, method string) { //统计http请求耗时分布
	begin := time.Now()
	<-ch
	metrics.RequestSummary.WithLabelValues("store_server", method).Observe(time.Now().Sub(begin).Seconds())
	return
}

//响应部分共用
func RenderJosn(w http.ResponseWriter, v interface{}) {
	d, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//w.Header().Set("Date", time.Now().Format("2006-01-02 15:04:05"))
	//w.Header().Set("Content-Encoding", "gzip")
	w.Write(d)
}

func APIWrapRsp(code int, inf string, payload interface{}) *WrapRsp {
	var msg string = ""
	if ix, ok := ErrMap[code]; ok && len(ix) != 0 {
		msg = fmt.Sprintf("%s, %s", ix, inf)
	} else {
		msg = inf
	}
	ret := &WrapRsp{
		Code:   code,
		ErrMsg: msg,
		Data:   payload,
	}
	return ret
}

//wrap response
type WrapRsp struct {
	Code   int         `json:"code"`   //`json:"code"`
	ErrMsg string      `json:"errmsg"` //`json:"errmsg"`
	Data   interface{} `json:"data"`
}

func UnmarshalInfos(info []byte, dec interface{}) (err error) {
	return json.Unmarshal(info, &dec)
}

/* ----------------------- 类型转换 ---------------------*/
func Str2Int(str string, bit int) int64 {
	if len(str) == 0 {
		return 0
	}
	r, err := strconv.ParseInt(str, 10, bit)
	if err != nil {
		logger.Entry().Error("parse string[%s] to int err[%v]", str, err)
		r = 0
	}
	return r
}

func Str2Uint(str string, bit int) uint64 {
	if len(str) == 0 {
		return 0
	}
	r, err := strconv.ParseUint(str, 10, bit)
	if err != nil {
		logger.Entry().Error("parse string[%s] to uint err[%v]", str, err)
		r = 0
	}
	return r
}

//常用封装
func MapToQueryString(data map[string]string) io.Reader {
	var buf bytes.Buffer
	for k, v := range data {
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(v)
		buf.WriteByte('&')
	}
	return bytes.NewReader(buf.Bytes())
}

func DoGet(url string, params map[string]string) (data []byte, err error) { //GET请求
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	logger.Entry().Debugf("do get requests url: %s", req.URL.String())
	//TODO SOMETHING
	doReq(req)
	return
}

func DoPost(url string, payload []byte) (data []byte, err error) { //POST请求
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return
	}
	logger.Entry().Debugf("do post requests url: %s", req.URL.String())
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Transfer-Encoding", "identity")
	//TODO SOMETHING
	doReq(req)
	return
}

func doReq(req *http.Request) (data []byte, err error) { //执行http请求
	if req.Header.Get("Authorization") == "" {
		//req.SetBasicAuth(g.Config().MediaHttp.Auth.Username, g.Config().Media.Auth.Password)
	}
	var httpClient = &http.Client{}
	httpClient.Timeout = time.Duration(g.Config().Http.HttpTimeout) * time.Second
	rsp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer rsp.Body.Close()
	data, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	return data, err
}

func GetHttpRsp(api, method string, payload io.Reader, isJson bool) ([]byte, error) { //通用http请求封装
	r, err := http.NewRequest(method, api, payload)
	if err != nil {
		logger.Entry().Errorf("new request err: %v", err)
		return nil, err
	}
	if isJson {
		r.Header.Set("Content-Type", "application/json")
	} else {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	rinfo, err := http.DefaultClient.Do(r)
	if err != nil {
		err = fmt.Errorf("query api: %s err: %v", api, err)
		return nil, err
	}
	if rinfo.StatusCode != 200 {
		err = fmt.Errorf("request api: %s response code: %v", api, rinfo.StatusCode)
		return nil, err
	}
	defer rinfo.Body.Close()
	rdata, err := ioutil.ReadAll(rinfo.Body)
	if err != nil {
		return nil, err
	}
	return rdata, nil
}

//size class
type Size interface {
	Size() int64
}

/*monitor data, elastic search log ...*/
func gen_log_doc_es(args map[string]interface{}) (id string, doc map[string]interface{}) {
	doc = map[string]interface{}{
		"operation_str":  "store_server_operate",
		"data_type_str":  "video", //default
		"region_id_int":  0,
		"operate_time":   time.Now().UnixNano() / 1e6, //uint for milliseconds
		"from_str":       "store_server_http",
		"operator_str":   "cms",
		"data_id_str":    "",
		"data_count_int": 0,
		"detail_text":    "",
	}
	if op, ok := args["operate"]; ok {
		doc["operation_str"] = op
	}
	if data_type, ok := args["data_type"]; ok {
		doc["data_type_str"] = data_type
	}
	if operator, ok := args["operator"]; ok {
		doc["operator_str"] = operator
	}
	if id_str, ok := args["data_id"]; ok {
		doc["data_id_str"] = id_str
	}
	if count, ok := args["data_count"]; ok {
		doc["data_count_int"] = count
	}
	if detail, ok := args["detail"]; ok {
		doc["detail_text"] = detail
	}
	if region, ok := args["region"]; ok {
		doc["region_id_int"] = region
	}
	id = fmt.Sprintf("%s_%s_%s_%s", doc["operation_str"], doc["data_type_str"], doc["operate_time"], doc["region_id_int"])
	return
}

func WriteLogEs(result, innerId, detail, dataType string) (err error) { //将日志写入es
	args := map[string]interface{}{
		"operate":    fmt.Sprintf("db_operate_%s", result),
		"data_id":    innerId,
		"detail":     detail,
		"data_type":  dataType,
		"data_count": 1,
	}
	id, doc := gen_log_doc_es(args)
	err = es.EsDriver.UpsertOne(g.Config().Es.Index, g.Config().Es.Type, id, doc)
	if err != nil {
		logger.Entry().Errorf("write operate log to es error: %v|id: %v|doc: %v", err, id, doc)
	}
	return
}
