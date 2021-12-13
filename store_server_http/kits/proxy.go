package kits

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"

	rpcjson "github.com/gorilla/rpc/json"
	"github.com/store_server/logger"
	"github.com/store_server/utils/common"
)

//rpc body
type RpcBody struct {
	Request json.RawMessage `json:"request"`
	Method  string          `json:"method"`
	Service string          `json:"service,omitempty"` //保持与music cms兼容
}

//proxy handle
type ProxyHandle struct {
	proxy *httputil.ReverseProxy
}

func NewProxyHandle(target string) (proxy *httputil.ReverseProxy) {
	proxy = &httputil.ReverseProxy{}
	host := common.GetLocalIP()
	if len(host) == 0 {
		host = "127.0.0.1"
	}
	director := func(req *http.Request) {
		logger.Entry().Infof("src: %v url: %v scheme: %v path: %v", req.RemoteAddr, req.URL, req.URL.Scheme, req.URL.Path)
		if strings.Contains(req.URL.Path, "rpc") {
			if len(target) == 0 {
				target = fmt.Sprintf("%s:9090", host)
			}
			req.URL.Host = fmt.Sprintf("%s", target)
			req.Host = req.URL.Host
			req.URL.Path = "/rpc"
			req.Header.Set("Content-Type", "application/json")
			if req.Body != nil { //extract body and do rpc json encode
				bodyBytes, _ := ioutil.ReadAll(req.Body)
				bc := &RpcBody{}
				e := json.Unmarshal(bodyBytes, bc)
				if e != nil {
					logger.Entry().Errorf("json decode http request body for proxy error: %v", e)
				}
				msg, err := rpcjson.EncodeClientRequest(bc.Method, bc.Request)
				if err != nil {
					logger.Entry().Errorf("rpc json encode for proxy error: %v", err)
				}
				req.Header.Set("Content-Length", strconv.Itoa(len(msg))) //update content length
				req.ContentLength = int64(len(msg))
				req.Body = ioutil.NopCloser(bytes.NewBuffer(msg))
			}
		} else {
			logger.Entry().Errorf("it's not a rpc call for media http, please check url: %v", req.URL)
		}
		if req.URL.Host == "" {
			req.URL.Host = "cms.joox.ibg.com"
		}
		if req.URL.Scheme == "" {
			req.URL.Scheme = "http"
		}
		if req.URL.Path == "" {
			req.URL.Path = "/rpc"
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "store_server_proxy")
		}
	}
	modifyResponse := func(response *http.Response) error {
		//DO NOTHING
		return nil
	}
	proxy.Director = director
	proxy.ModifyResponse = modifyResponse
	return proxy
}

func (ph *ProxyHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ph.proxy.ServeHTTP(w, r)
}
