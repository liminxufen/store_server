package kits

//"errors define"
import (
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/store_server/logger"
)

//定义常见错误码及消息
const (
	ErrInnerServer = 500
	ErrCustom      = 1
)

const (
	ErrParams = iota + 400
)

const (
	ErrOther = iota + 402
)

var ErrMap = map[int]string{
	ErrInnerServer: "Inner Server Error",
	ErrParams:      "Parameters Error",
	ErrOther:       "",
}

//定义错误捕获处理
//func CatchErr(module_name string, eP *error, logger *lg.Logger, vs ...interface{}) {
func CatchErr(module_name string, eP *error, logger *logger.JLoggerEntry, vs ...interface{}) {
	if rE := recover(); rE != nil {
		*eP = fmt.Errorf("panic|%v", rE.(error))
	}
	if *eP != nil {
		s := []interface{}{module_name, *eP}
		segs := []string{}
		s = append(s, vs...)
		for _, _ = range s {
			segs = append(segs, "%v")
		}
		*eP = fmt.Errorf(strings.Join(segs, "|"), s...)
		sentry.CaptureException(*eP)
		sentry.Flush(time.Second * 1)
	}
	if logger != nil && *eP != nil {
		//logger.Error((*eP).Error())
	}
	return
}

func PanicMsg(msg string) { //抛出指定Panic异常信息
	if len(msg) == 0 {
		panic("")
	}
	panic(msg)
	return
}

func HandlePanicMsg(err *error) { //捕获Panic异常
	if rerr := recover(); rerr != nil {
		logger.Entry().Infof("%v", rerr)
		if str, ok := rerr.(string); ok {
			if len(str) == 0 {
				*err = nil
				return
			}
			*err = errors.New(str)
		} else {
			*err = errors.New("Some Error happend inside package!!!" + PadStack())
		}
		sentry.CaptureException(*err)
	}
	return
}

func PadStack() string { //格式化调用栈信息
	return "\n-------------------------------------------\n" +
		string(debug.Stack()) + "\n-------------------------------------------\n"
}
