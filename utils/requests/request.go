package requests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/store_server/utils/errors"
	"github.com/store_server/utils/rest"
)

type option struct {
	Configs     []rest.ClientConfig
	Headers     []string
	ContentType string
	Cookie      *http.Cookie
	Query       url.Values
	Body        []byte
	Error       error
}

type Option func(*option)

func WithBody(contentType string, body []byte) Option {
	return func(o *option) {
		o.ContentType = contentType
		o.Body = body
	}
}

var defaultConfig = rest.ClientConfig{
	TraceOption: &rest.TraceOption{
		TraceRequestBody:   true,
		TraceRequestHeader: true,
		TraceRespBody:      false,
	},
}

func WithConfig(configs ...rest.ClientConfig) Option {
	return func(o *option) {
		o.Configs = append([]rest.ClientConfig{defaultConfig}, configs...)
	}
}

func WithQuery(params url.Values) Option {
	return func(o *option) {
		o.Query = params
	}
}

func WithHeaders(kv ...string) Option {
	return func(o *option) {
		o.Headers = kv
	}
}

func WithBodyJSON(v interface{}) Option {
	return func(o *option) {
		o.ContentType = rest.ContentTypeJSON

		body, err := json.Marshal(v)
		if err != nil {
			o.Body, o.Error = nil, errors.Errorf(err, "json.Marshal failed")
		} else {
			o.Body, o.Error = body, nil
		}
	}
}

func WithCookie(cookie *http.Cookie) Option {
	return func(o *option) {
		o.Cookie = cookie
	}
}

func newReqOption(options ...Option) *option {
	var opt = new(option)
	for _, f := range options {
		if f == nil {
			continue
		}
		f(opt)
	}
	return opt
}

func (o *option) CreateClient(ctx context.Context, uri string) (*rest.BaseClient, error) {
	if o.Error != nil {
		return nil, o.Error
	}

	var target string
	if encoded := o.Query.Encode(); encoded != "" {
		target = fmt.Sprintf("%s?%s", uri, o.Query.Encode())
	} else {
		target = uri
	}

	client := rest.NewBaseHttpClient(ctx, target, o.Configs...)

	if o.ContentType != "" {
		err := client.CustomHeaders(rest.HeaderContentType, o.ContentType)
		if err != nil {
			return nil, err
		}
	}

	if o.Cookie != nil {
		client.CustomCookie = o.Cookie
	}

	if len(o.Headers) > 0 {
		if err := client.CustomHeaders(o.Headers...); err != nil {
			return nil, err
		}
	}

	return client, nil
}

func Get(ctx context.Context, dest interface{}, url string, options ...Option) error {
	if dest == nil {
		return nil
	}

	client, err := newReqOption(options...).CreateClient(ctx, url)
	if err != nil {
		return errors.Errorf(err, "failed to create http client")
	}

	result, err := client.GetAndParse()
	if err != nil {
		return err
	}

	if !result.Success {
		return errors.Errorf(nil, "response is not success")
	}

	if err := json.Unmarshal(result.Result, dest); err != nil {
		return err
	}
	return nil
}

func Post(ctx context.Context, dest interface{}, url string, options ...Option) error {
	if dest == nil {
		return nil
	}

	opt := newReqOption(options...)
	if opt.Body == nil {
		opt.Body = make([]byte, 0)
	}

	client, err := opt.CreateClient(ctx, url)
	if err != nil {
		return errors.Errorf(err, "failed to create http client")
	}

	result, err := client.PostAndParse(opt.Body)
	if err != nil {
		return err
	}

	if !result.Success {
		return errors.Errorf(nil, "response is not success")
	}

	if err := json.Unmarshal(result.Result, dest); err != nil {
		return err
	}
	return nil
}

// PostThird 向其他系统发起post请求(非rest.APIJSONResult返回结构)
func PostThird(ctx context.Context, dest interface{}, url string, options ...Option) error {
	if dest == nil {
		return nil
	}

	opt := newReqOption(options...)
	if opt.Body == nil {
		opt.Body = make([]byte, 0)
	}

	client, err := opt.CreateClient(ctx, url)
	if err != nil {
		return errors.Errorf(err, "failed to create http client")
	}

	result, err := client.Post(opt.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(result, dest); err != nil {
		return err
	}
	return nil
}

func PostThirdWithBody(ctx context.Context, url string, options ...Option) ([]byte, error) {

	opt := newReqOption(options...)
	if opt.Body == nil {
		opt.Body = make([]byte, 0)
	}

	client, err := opt.CreateClient(ctx, url)
	if err != nil {
		return nil, errors.Errorf(err, "failed to create http client")
	}

	result, err := client.Post(opt.Body)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func Put(ctx context.Context, dest interface{}, url string, options ...Option) error {
	if dest == nil {
		return nil
	}

	opt := newReqOption(options...)
	if opt.Body == nil {
		opt.Body = make([]byte, 0)
	}

	client, err := opt.CreateClient(ctx, url)
	if err != nil {
		return errors.Errorf(err, "failed to create http client")
	}

	result, err := client.PutAndParse(opt.Body)
	if err != nil {
		return err
	}

	if !result.Success {
		return errors.Errorf(nil, "response is not success")
	}

	if err := json.Unmarshal(result.Result, dest); err != nil {
		return err
	}
	return nil
}
