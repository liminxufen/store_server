package common

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/store_server/utils/errors"
	"github.com/store_server/utils/requests"
	"github.com/store_server/utils/rest"
)

func DoPostRequest(ctx context.Context, url string, data []byte, rsp interface{},
	extras ...interface{}) error {
	opts := make([]requests.Option, 0)
	opts = append(opts, requests.WithBody(rest.ContentTypeJSON, data))
	if len(extras) > 0 {
		proxy := extras[0].(string)
		opts = append(opts, requests.WithConfig(rest.ClientConfig{Proxy: proxy, Timeout: 600}))
	} else {
		opts = append(opts, requests.WithConfig(rest.ClientConfig{Timeout: 600}))
	}
	return requests.PostThird(ctx, rsp, url, opts...)
}

func DoRequest(ctx context.Context, url, method string, body []byte) (data []byte, err error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	if strings.ToUpper(method) == "GET" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req.Header.Set("Content-Type", "application/json")
	}
	var httpClient = &http.Client{}
	//httpClient.Timeout = time.Duration(g.Config().Http.HttpTimeout) * time.Second
	rsp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	data, err = ioutil.ReadAll(rsp.Body)
	if rsp.StatusCode != 200 {
		return nil, errors.Errorf(nil, "http response code: %v|response body: %v",
			rsp.StatusCode, bytes.NewBuffer(data).String())
	}
	if err != nil {
		return nil, err
	}
	return data, nil
}
