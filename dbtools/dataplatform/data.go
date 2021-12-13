package dataplatform

/** -------------- 封装到数据平台的查询操作 ------------- **/

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/store_server/logger"
	"github.com/store_server/store_server_http/g"
	"github.com/store_server/utils/common"
	"github.com/store_server/utils/errors"
	"net/http"
	"strings"
	"sync"
)

//dataplatform driver
type DataplatformDriver struct {
	ctx  context.Context
	lock sync.RWMutex
	url  string
}

func NewDataplatformDriver(ctx context.Context) *DataplatformDriver {
	driver := &DataplatformDriver{ctx: ctx}
	driver.url = driver.getSearchURL()
	return driver
}

var (
	DpDriver *DataplatformDriver
)

type SearchReq struct {
	Sql string `json:"sql"`
}

type SearchRsp struct {
	Ret     int         `json:"ret"`
	RetInfo string      `json:"ret_info"`
	Total   interface{} `json:"total"`
	Fields  []string    `json:"fields"`
	Items   []*dataItem `json:"items"`
}

type dataItem struct {
	Id     string  `json:"id"`
	Fields []*Item `json:"fields"`
}

type Item struct {
	FieldName  string      `json:"field_name"`
	FieldValue interface{} `json:"value"`
}

func (dpd *DataplatformDriver) doSearch(url string, req *SearchReq, rsp *SearchRsp) error {
	if req == nil {
		return errors.Errorf(nil, "invalid request")
	}
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	logger.Entry().Debugf("search request: %v", bytes.NewBuffer(body).String())
	//err = common.DoPostRequest(dpd.ctx, url, body, rsp)
	data, err := common.DoRequest(dpd.ctx, url, http.MethodPost, body)
	if err != nil {
		logger.Entry().Errorf("search request to dataplatform error: %v|url: %v", err, url)
		return err
	}
	err = json.Unmarshal(data, rsp)
	if err != nil {
		logger.Entry().Errorf("decode dataplatform response error: %v|response: %v",
			err, bytes.NewBuffer(data).String())
		return err
	}
	return nil
}

func (dpd *DataplatformDriver) getSearchURL() string {
	var addrs string
	var err error
	if err != nil {
		logger.Entry().Errorf("get dataplatform server address by mod_id: %v|and cmd_id: %v|error: %v",
			g.Config().Dataplatform.ModId, g.Config().Dataplatform.CmdId, err)
		addrs = fmt.Sprintf("%s:%d", g.Config().Dataplatform.Host, g.Config().Dataplatform.Port)
	} else {
		addrs = ""
	}
	addrs = fmt.Sprintf("http://%s%s", addrs, g.Config().Dataplatform.Api)
	//addrs = g.Config().Dataplatform.Api
	logger.Entry().Debugf("get dataplatform server address: %s", addrs)
	return addrs
}

func (dpd *DataplatformDriver) withOrder(orders [][2]string) string {
	if orders == nil || len(orders) <= 0 {
		return ""
	}
	format := make([]string, 0)
	values := make([]interface{}, 0)
	for _, order := range orders {
		format = append(format, "%s %s")
		if len(order) == 2 {
			if strings.ToLower(order[1]) != "asc" && strings.ToLower(order[1]) != "desc" {
				order[1] = "desc"
			}
			values = append(values, order[0], order[1])
		}
	}
	formatStr := common.JoinString(format, ", ")
	if len(values) > 0 {
		orderSt := fmt.Sprintf(formatStr, values...)
		orderSt = fmt.Sprintf("order by %s", orderSt)
		return orderSt
	}
	return ""
}

func (dpd *DataplatformDriver) transferScopeOp(opir interface{}) string {
	op := opir.(string)
	switch op {
	case "gt", "GT":
		op = ">"
	case "ge", "GE":
		op = ">="
	case "eq", "EQ":
		op = "=="
	case "ne", "NE":
		op = "!="
	case "lt", "LT":
		op = "<"
	case "le", "LE":
		op = "<="
	default:
	}
	return op
}

func (dpd *DataplatformDriver) withScope(scopes [][3]interface{}) string {
	if scopes == nil || len(scopes) <= 0 {
		return ""
	}
	format := make([]string, 0)
	values := make([]interface{}, 0)
	for _, scope := range scopes {
		format = append(format, "%v %v \"%v\"")
		if len(scope) == 3 {
			/*op := scope[1]
			buffer := bytes.NewBuffer([]byte{})
			jsonEncoder := json.NewEncoder(buffer)
			jsonEncoder.SetEscapeHTML(false)
			jsonEncoder.Encode(op)*/
			values = append(values, scope[0], dpd.transferScopeOp(scope[1]), fmt.Sprintf("%v", scope[2]))
		}
	}
	formatStr := common.JoinString(format, " and ")
	if len(values) > 0 {
		where := fmt.Sprintf(formatStr, values...)
		return where
	}
	return ""
}

//数据平台SQL搜索
func (dpd *DataplatformDriver) SearchBySql(sql string,
	opts ...interface{}) (interface{}, int64, error) {
	if len(sql) <= 0 {
		return nil, 0, errors.Errorf(nil, "invalid sql statement")
	}
	page, pagesize := 1, 100
	if len(opts) > 0 {
		page = opts[0].(int)
	}
	if len(opts) > 1 {
		pagesize = opts[1].(int)
	}
	offset := (page - 1) * pagesize
	if !strings.Contains(strings.ToLower(sql), "limit") {
		sql = fmt.Sprintf("%s limit %d, %d", sql, offset, pagesize)
	}
	//logger.Entry().Debugf("search sql: %v", sql)
	req := &SearchReq{Sql: sql}
	rsp := &SearchRsp{}
	err := dpd.doSearch(dpd.url, req, rsp)
	if err != nil {
		return nil, 0, err
	}
	if rsp != nil && rsp.Ret != 0 {
		return nil, 0, errors.Errorf(nil, "search error: %v", rsp.RetInfo)
	}
	total, _ := common.Interface2Int64(rsp.Total)
	res := make([]*dataItem, 0, len(rsp.Items))
	res = append(res, rsp.Items...)
	return res, total, nil
}

//数据平台条件搜索
func (dpd *DataplatformDriver) SearchByCondition(datatype int, queryFields []string,
	conds map[string]interface{}, page, pagesize int64, opts ...interface{}) (interface{}, int64, error) {
	fields := "*"
	if len(queryFields) > 0 {
		fields = common.JoinString(queryFields, ",")
	}
	sql := fmt.Sprintf("select %s from data_set_%d", fields, datatype)
	//拼接where部分
	where := ""
	format := make([]string, 0)
	values := make([]interface{}, 0)
	for k, v := range conds {
		format = append(format, "%s=\"%s\"")
		values = append(values, k, fmt.Sprintf("%v", v))
	}
	formatStr := common.JoinString(format, " and ")
	if len(values) > 0 {
		where = fmt.Sprintf(formatStr, values...)
		where = common.JoinString([]string{"where", where}, " ")
	}
	if len(opts) > 1 {
		scopes, ok := opts[1].([][3]interface{})
		if ok && len(scopes) > 0 {
			extraWhere := dpd.withScope(scopes)
			if len(extraWhere) > 0 {
				if len(where) > 0 {
					where = common.JoinString([]string{where, extraWhere}, " and ")
				} else {
					where = common.JoinString([]string{where, extraWhere}, " where ")
				}
			}
		}
	}
	//拼接order部分
	orders := ""
	if len(opts) > 0 {
		orderBy, ok := opts[0].([][2]string)
		if ok && len(orderBy) > 0 {
			orders = dpd.withOrder(orderBy)
		}
	}
	if page <= 0 {
		page = 1
	}
	if pagesize <= 0 {
		pagesize = 100
	}
	//拼接limits部分
	offset := (page - 1) * pagesize
	limits := fmt.Sprintf("limit %d, %d", offset, pagesize)
	sql = common.JoinString([]string{sql, where, orders, limits}, " ")
	return dpd.SearchBySql(sql)
}
