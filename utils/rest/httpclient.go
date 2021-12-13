package rest

import (
	"context"
	"fmt"
	"github.com/store_server/utils/ctxutil"
	"github.com/store_server/utils/errors"
	"net/http"
	"net/url"
	"time"
)

type TraceOption struct {
	TraceRequestHeader bool
	TraceRequestBody   bool
	TraceRespBody      bool
}

type ClientConfig struct {
	RequestID            string
	EndPointURLToken     string        // URL Token
	Timeout              int           // 单位为秒
	NumMaxRetries        int           // 最大重试次数
	RetryMaxTimeDuration time.Duration // 总计重试时间 ExpireTime
	TraceOption          *TraceOption  // 追踪打印属性
	IsExternalURL        bool          // 默认为FALSE,是否为调用对外api,如是，则不在header的UserAgent中带上内网ip.数据脱敏
	IsRESTStatusCode     bool          // RESTFUL API风格，使用状态码来表示资源态.
	Proxy                string        //客户端代理
}

type Client struct {
	ctx          context.Context
	HTTPClient   *http.Client
	CustomHeader http.Header
	CustomCookie *http.Cookie
	EndPointURL  string
	ClientConfig
}

const (
	_MinTime        = 5  // 最大随机重试休眠秒数
	_DefaultTimeout = 60 // 默认超时时间(s)
	_DefaultRetries = 5  // 默认重试次数

)

// 最大重试间隔时间 = (默认超时时间 + 最大随机重试Sleep秒数) * 默认重试次数 = (60 + 5) * 5s = 325s
const _MaxRetryTimeDuration = time.Duration((_DefaultTimeout+_MinTime)*_DefaultRetries) * time.Second

func GetHTTPClient(t int, p ...ClientConfig) *http.Client {
	var proxy string
	if len(p) > 0 {
		proxy = p[0].Proxy
	}
	proxyUrl, _ := url.Parse(fmt.Sprintf("http://%s", proxy))
	transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	client := &http.Client{
		Timeout: time.Duration(getRealTimeout(t)) * time.Second,
	}
	if len(proxy) > 0 {
		client.Transport = transport
	}
	return client
}

func getRealTimeout(t ...int) int {
	var timeout int
	if len(t) > 0 {
		for _, v := range t {
			timeout = v
		}
	}
	if timeout <= 0 {
		timeout = _DefaultTimeout
	}
	return timeout
}

func (c *Client) initCustomHeader() {
	if c.CustomHeader == nil {
		c.CustomHeader = http.Header{}
	}

	// 默认JSON
	c.CustomHeader.Set(HeaderContentType, ContentTypeJSON)
}

func (c *Client) addHeaderRequestID() {
	if c.CustomHeader.Get(HeaderRequestId) != "" {
		return
	}
	if c.RequestID == "" {
		c.RequestID = ctxutil.RequestIDFromContext(c.ctx)
	}
	c.CustomHeader.Add(HeaderRequestId, c.RequestID)
}

func (c *Client) GetHTTPClient() *http.Client {
	if c.HTTPClient == nil {
		c.HTTPClient = GetHTTPClient(_DefaultTimeout)
	}
	return c.HTTPClient
}

func (c *Client) CustomHeaders(keyValues ...string) error {
	if len(keyValues)%2 != 0 {
		return errors.Errorf(nil, "wrong key value,len:%d", len(keyValues))
	}
	var key string
	for i := 0; i*2 < len(keyValues); i++ {
		key = keyValues[i*2]
		if key != "" {
			_, ok := _standHeaders[key]
			if ok {
				c.CustomHeader.Set(key, keyValues[i*2+1])
			} else {
				c.CustomHeader.Add(key, keyValues[i*2+1])
			}
		}
	}
	return nil
}

var _standHeaders = map[string]struct{}{
	HeaderContentType:    {},
	HeaderUserAgent:      {},
	HeaderAcceptLanguage: {},
	HeaderAuthorization:  {},
	HeaderCookie:         {},
	HeaderRequestId:      {},
}

func copyHTTPHeader(req *http.Request, head http.Header) {
	if req == nil {
		return
	}
	for k, v := range head {
		for _, vv := range v {
			_, ok := _standHeaders[k]
			if ok {
				req.Header.Set(k, vv)
			} else {
				req.Header.Add(k, vv)
			}
		}
	}
}
