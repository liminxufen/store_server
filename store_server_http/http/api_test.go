package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/store_server/dbtools/driver"
	"github.com/store_server/logger"
	"github.com/store_server/store_server_http/g"
	"github.com/store_server/store_server_http/kits"
	"github.com/stretchr/testify/assert"
)

var (
	storeServerhttp *StoreServerHttp
)

func apiSetup() {
	g.Config().MusicMysql = driver.DefaultTestScheme.Mysql
	g.Config().IpWhiteList = "127.0.0.1"
	g.Config().Http.Listen = ":9080"
	InitIpWhiteList(g.Config().IpWhiteList)
	storeServerhttp, _ = NewDefaultStoreServerHttp()
	/*go func() {
		storeServerhttp.Start()
	}()*/
}

func apiCleanup() {
	storeServerhttp.closeDoneChan <- struct{}{}
	storeServerhttp.Stop()
	g.Config().IpWhiteList = ""
	kits.IPWhiteLst = []string{}
}

func readFile(writer *io.Writer, filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(*writer, f)
	if err != nil {
		return err
	}
	return nil
}

func genUploadVideoReq() (io.Reader, error) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	formFile, err := writer.CreateFormFile("video", "video_example.mp4")
	if err != nil {
		return nil, err
	}
	err = readFile(&formFile, "./video_example.mp4")
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func genUploadAudioReq() (io.Reader, error) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	formFile, err := writer.CreateFormFile("audio", "audio_example.mp3")
	if err != nil {
		return nil, err
	}
	err = readFile(&formFile, "./audio_example.mp3")
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func mockRequest(r *gin.Engine, method, path string, data interface{}) *httptest.ResponseRecorder {
	var body io.Reader
	var req *http.Request
	if data != nil {
		rd, err := json.Marshal(data)
		if err != nil {
			logger.Entry().Error(err)
			return nil
		}
		body = bytes.NewReader(rd)
	} else {
		body = nil
	}
	req, _ = http.NewRequest(method, path, body)
	if method == "POST" {
		//req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Content-Type", "multipart/form-data")
	} else if method == "GET" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Set("X-Forwarded-For", "127.0.0.1")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestRequestUploadVideo(t *testing.T) {
	rsp := gin.H{
		"code":   0,
		"errmsg": "ok",
		"data":   nil,
	}
	apiSetup()
	body, err := genUploadVideoReq()
	assert.NoError(t, err)

	fmt.Println("start mock request to upload video...")
	w := mockRequest(router, "POST", "/static_resource/videos/upload", body)
	assert.Equal(t, http.StatusOK, w.Code)

	var response kits.WrapRsp
	e := kits.UnmarshalInfos([]byte(w.Body.String()), &response)
	assert.NoError(t, e)
	assert.Equal(t, response.ErrMsg, rsp["errmsg"])
}

func TestRequestUploadAudio(t *testing.T) {
	body, err := genUploadAudioReq()
	assert.NoError(t, err)

	fmt.Println("start mock request to upload audio...")
	w := mockRequest(router, "POST", "/static_resource/audios/upload", body)
	assert.Equal(t, http.StatusOK, w.Code)

	var response kits.WrapRsp
	e := kits.UnmarshalInfos([]byte(w.Body.String()), &response)
	assert.NoError(t, e)
	assert.Equal(t, response.ErrMsg, "ok")
}

func TestTranscodeCallbackVideo(t *testing.T) {

}

func TestTranscodeCallbackAudio(t *testing.T) {

}

func TestGetExternalVideoInfo(t *testing.T) {

}

func TestGetExternalAudioInfo(t *testing.T) {

}

func TestCleanup(t *testing.T) {
	apiCleanup()
}
