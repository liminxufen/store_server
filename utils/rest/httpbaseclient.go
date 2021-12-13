package rest

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/store_server/logger"
	"github.com/store_server/utils/errors"
)

type BaseClient struct {
	Client
	HTTPRequest  *http.Request  // 用于返回请求，直接设值无效。引用时注意做非空检查。
	HTTPResponse *http.Response // 用于返回执行结果。引用时注意做非空检查。
}

func NewBaseHttpClient(ctx context.Context, url string, configs ...ClientConfig) *BaseClient {
	c := &BaseClient{}
	c.ctx = ctx
	c.EndPointURL = url

	if len(configs) > 0 {
		for _, cf := range configs {
			c.ClientConfig = cf
		}
		if c.Timeout <= 0 {
			c.Timeout = _DefaultTimeout
		}
	} else {
		c.NumMaxRetries = _DefaultRetries
		c.RetryMaxTimeDuration = _MaxRetryTimeDuration
		c.Timeout = _DefaultTimeout
	}

	if c.TraceOption == nil {
		c.TraceOption = &TraceOption{}
	}

	c.HTTPClient = GetHTTPClient(c.Timeout, c.ClientConfig)

	c.CustomHeader = http.Header{}
	c.initCustomHeader()

	// c.HTTPRequest = &http.Request{Header: c.CustomHeader}
	// c.HTTPRequest = c.HTTPRequest.WithContext(ctx) // 传递ctx

	return c
}

func (bc *BaseClient) Get() ([]byte, error) {
	bc.addHeaderRequestID()
	return bc.doGet()
}

func (bc *BaseClient) doGet() ([]byte, error) {
	startTime := time.Now()

	if bc.EndPointURL == "" {
		return nil, errors.Errorf(nil, "EndPointURL empty")
	}

	req, err := http.NewRequestWithContext(bc.ctx, http.MethodGet, bc.EndPointURL, nil)
	if err != nil {
		return nil, errors.Errorf(err, "耗时:%s", time.Since(startTime))
	}

	copyHTTPHeader(req, bc.CustomHeader)

	// 设置HTTP.Cookie
	if bc.CustomCookie != nil {
		req.AddCookie(bc.CustomCookie)
	}

	if bc.IsExternalURL {
		req.Header.Set(HeaderUserAgent, _DefaultUserAgent)
	} else {
		req.Header.Set(HeaderUserAgent, "")
	}

	bc.HTTPRequest = req

	client := bc.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Errorf(err, "耗时:%s", time.Since(startTime))
	}
	bc.HTTPResponse = resp

	statusCode := resp.StatusCode
	if !bc.IsRESTStatusCode {
		if statusCode != http.StatusOK {

		}
	}

	if resp.Body == nil {
		return nil, errors.Errorf(nil, "response body is empty,url:%s,StatusCode:%d", bc.EndPointURL, statusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		_ = resp.Body.Close()
		return nil, errors.Errorf(err, "耗时:%s", time.Since(startTime))
	}
	_ = resp.Body.Close()

	if bc.TraceOption.TraceRespBody {
		logger.Infof("%s url:%s,body:%s", bc.RequestID, bc.EndPointURL, string(respBody))
	}

	if bc.IsRESTStatusCode && statusCode != http.StatusOK {
		// RESTFUL风格时可获取返回的body自行解析
		return respBody, errors.Errorf(nil, "HTTP status codes != 200")
	}

	return respBody, nil
}

func (bc *BaseClient) Post(body []byte) ([]byte, error) {
	bc.addHeaderRequestID()
	return bc.doPost(body)
}

func (bc *BaseClient) doPost(body []byte) ([]byte, error) {
	return bc.doUpdateRequest(http.MethodPost, body)
}

func (bc *BaseClient) doPut(body []byte) ([]byte, error) {
	return bc.doUpdateRequest(http.MethodPut, body)
}

func (bc *BaseClient) doDelete() ([]byte, error) {
	return bc.doUpdateRequest(http.MethodDelete, nil)
}

func (bc *BaseClient) doUpdateRequest(method string, body []byte) ([]byte, error) {
	startTime := time.Now()

	if bc.EndPointURL == "" {
		return nil, errors.Errorf(nil, "EndPointURL empty")
	}

	if method != http.MethodDelete && body == nil {
		return nil, errors.Errorf(nil, "body empty")
	}

	req, err := http.NewRequestWithContext(bc.ctx, method, bc.EndPointURL, bytes.NewReader(body))
	if err != nil {
		return nil, errors.Errorf(err, "耗时:%s", time.Since(startTime))
	}

	if bc.IsExternalURL {
		req.Header.Set(HeaderUserAgent, _DefaultUserAgent)
	} else {
		req.Header.Set(HeaderUserAgent, "")
	}
	copyHTTPHeader(req, bc.CustomHeader)

	// 设置HTTP.Cookie
	if bc.CustomCookie != nil {
		req.AddCookie(bc.CustomCookie)
	}

	bc.HTTPRequest = req

	client := bc.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Errorf(err, "耗时:%s", time.Since(startTime))
	}
	bc.HTTPResponse = resp

	statusCode := resp.StatusCode
	if !bc.IsRESTStatusCode {
		if statusCode != http.StatusOK {

		}
	}

	if resp.Body == nil {
		return nil, errors.Errorf(nil, "response body is empty,url:%s,StatusCode:%d", bc.EndPointURL, statusCode)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Errorf(err, "耗时:%s", time.Since(startTime))
	}

	if bc.TraceOption.TraceRespBody {
		logger.Infof("%s url:%s,body:%s", bc.RequestID, bc.EndPointURL, string(respBody))
	}

	if bc.IsRESTStatusCode && statusCode != http.StatusOK {
		// RESTFUL风格时可获取返回的body自行解析
		return respBody, errors.Errorf(nil, "HTTP status codes != 200")
	}
	return respBody, nil
}

func (bc *BaseClient) GetAndParse() (*APIJSONResult, error) {
	bc.addHeaderRequestID()
	body, err := bc.doGet()
	if err != nil {
		return nil, errors.Errorf(err, "do get")
	}
	return ParseResultJSON(body)
}

func (bc *BaseClient) PostAndParse(body []byte) (*APIJSONResult, error) {
	bc.addHeaderRequestID()
	body, err := bc.doPost(body)
	if err != nil {
		return nil, errors.Errorf(err, "do post")
	}
	return ParseResultJSON(body)
}

func (bc *BaseClient) PutAndParse(body []byte) (*APIJSONResult, error) {
	bc.addHeaderRequestID()
	body, err := bc.doPut(body)
	if err != nil {
		return nil, errors.Errorf(err, "do post")
	}
	return ParseResultJSON(body)
}
