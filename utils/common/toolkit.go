package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/getsentry/sentry-go"
	"github.com/store_server/logger"
)

func UnmarshalWithNumber(data []byte, v interface{}) error {
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	return d.Decode(v)
}

// 字符串拼接 builder方式
func JoinString(strs []string, sep string) string {
	switch len(strs) {
	case 0:
		return ""
	case 1:
		return strs[0]
	}
	n := len(sep) * (len(strs) - 1)
	for i := 0; i < len(strs); i++ {
		n += len(strs[i])
	}
	var builder strings.Builder
	builder.Grow(n)
	builder.WriteString(strs[0])
	for _, str := range strs[1:] {
		builder.WriteString(sep)
		builder.WriteString(str)
	}
	return builder.String()
}

//TimeCostTrack 追踪耗时
func TimeCostTrack(start time.Time, structName, funcName string, err error) {
	elapsed := float64(time.Since(start)) / (1000 * 1000)
	hostName, _ := os.Hostname()
	msgTpl := "Hostname: %s, data struct: %s in call method: %s, cost: %f ms"
	msg := fmt.Sprintf(msgTpl, hostName, structName, funcName, elapsed)
	logger.Entry().Infof("[TimeCostTrack] msg: %v", msg)
	if err != nil {
		sentry.CaptureException(err)
	}
}

//获取环境变量
func GetEnvOrDefault(key string, dft string) string {
	val := os.Getenv(key)
	if val == "" {
		return dft
	}
	return val
}

//获取当前work path
func GetWorkPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return dir
}

//GetLocalIP 获取本机ip
func GetLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		//logger.Entry().Error(err)
		return ""
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx]
}

//GetIPs 获取真实ip
func GetIPs() (ips []string) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		//logger.Entry().Errorf("fail to get net intehurface addrs: %v", err)
		return ips
	}

	for _, address := range interfaceAddr {
		ipNet, isValidIPNet := address.(*net.IPNet)
		if isValidIPNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips
}

func GetIPByMultiAddr(faces []string) (ip string, err error) {
	for _, face := range faces {
		if ip, err = GetIPByInter(face); err == nil {
			return
		}
	}
	return
}

//GetIPByInter 通过网卡获取ip
func GetIPByInter(eth string) (string, error) {
	interfaceNet, err := net.InterfaceByName(eth)
	if err != nil {
		//logger.Entry().Errorf("fail to get net interface addrs: %v", err)
		return "", err
	}

	addressList, err := interfaceNet.Addrs()
	if err != nil {
		//logger.Entry().Errorf("fail to get net interface addrs: %v", err)
		return "", err
	}
	for _, address := range addressList {
		ipNet, isValidIPNet := address.(*net.IPNet)
		if isValidIPNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", errors.New("may network interface is wrong!")
}

func Contain(sl []string, s string) bool {
	for _, v := range sl {
		if v == s {
			return true
		}
	}
	return false
}

func ByteSlice2String(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

func String2ByteSlice(bs string) []byte {
	return *(*[]byte)(unsafe.Pointer(&bs))
}

func GenRandom(start, end, count int) []int { //产生随机数
	if end < start || (end-start) < count {
		return nil
	}
	nums := make([]int, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano())) //随机数生成器，加入时间戳保证每次生成的随机数不一样
	for len(nums) < count {
		num := r.Intn(end-start) + start //生成随机数
		exist := false
		for _, v := range nums { //查重
			if v == num {
				exist = true
				break
			}
		}
		if !exist {
			nums = append(nums, num)
		}
	}
	return nums
}

func BytesBufferPool() *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
}

func Interface2Int(s interface{}) (int, error) {
	t := fmt.Sprintf("%v", s)
	if len(t) == 0 {
		return 0, nil
	}
	d, err := strconv.Atoi(t)
	if err != nil {
		return 0, err
	}
	return d, nil
}

func Interface2Int64(s interface{}) (int64, error) {
	t := fmt.Sprintf("%v", s)
	if len(t) == 0 {
		return 0, nil
	}
	d, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		return 0, err
	}
	return d, nil
}

func String2Int64(s string) (int64, error) {
	if len(s) == 0 {
		return 0, nil
	}
	d, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return d, nil
}

/*func RequestPost(url string, data []byte, resp interface{}) error {
	opts := make([]requests.Option, 0)
	opts = append(opts, requests.WithBody(rest.ContentTypeJSON, data))
	//opts = append(opts, requests.WithConfig(rest.ClientConfig{Proxy: proxy, Timeout: 10}))
	opts = append(opts, requests.WithConfig(rest.ClientConfig{Timeout: 10}))
	return requests.PostThird(context.TODO(), resp, url, opts...)
}*/

//map to query string
func MapToQueryString(data map[string]interface{}) []byte {
	var buf bytes.Buffer
	for k, v := range data {
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(fmt.Sprintf("%v", v))
		buf.WriteByte('&')
	}
	n := len(buf.Bytes())
	return buf.Bytes()[0 : n-1]
}
