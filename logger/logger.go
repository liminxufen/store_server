package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type Level uint8

const (
	LDebug Level = iota
	LInfo
	LWarn
	LError
	LFatal
)

var (
	DebugLogger  = log.New(os.Stdout, "[DEBUG]", log.LstdFlags|log.Lshortfile)
	InfoLogger   = log.New(os.Stdout, "[INFO]", log.LstdFlags)
	WarnLogger   = log.New(os.Stdout, "[WARN]", log.LstdFlags|log.Lshortfile)
	ErrorLogger  = log.New(os.Stdout, "[ERROR]", log.LstdFlags|log.Lshortfile)
	FatalLogger  = log.New(os.Stdout, "[FATAL]", log.LstdFlags|log.Lshortfile)
	CurrentLevel = LInfo
	Log          *Logger
)

func InitFatal(v ...interface{}) {
	s := fmt.Sprint(v...)
	FatalLogger.Output(2, s)
	panic(s)
}

func InitFatalf(f string, v ...interface{}) {
	s := fmt.Sprintf(f, v)
	FatalLogger.Output(2, s)
	panic(s)
}

/*打印各级别日志调用方法*/
func Debug(v ...interface{}) {
	if CurrentLevel <= LDebug {
		DebugLogger.Output(2, fmt.Sprint(v...))
	}
}

func Debugf(format string, v ...interface{}) {
	if CurrentLevel <= LDebug {
		DebugLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

func Info(v ...interface{}) {
	if CurrentLevel <= LInfo {
		InfoLogger.Output(2, fmt.Sprint(v...))
	}
}

func Infof(format string, v ...interface{}) {
	if CurrentLevel <= LInfo {
		InfoLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

func Warn(v ...interface{}) {
	if CurrentLevel <= LWarn {
		WarnLogger.Output(2, fmt.Sprint(v...))
	}
}

func Warnf(format string, v ...interface{}) {
	if CurrentLevel <= LWarn {
		WarnLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

func Error(v ...interface{}) {
	if CurrentLevel <= LError {
		ErrorLogger.Output(2, fmt.Sprint(v...))
	}
}

func Errorf(format string, v ...interface{}) {
	if CurrentLevel <= LError {
		ErrorLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

func Fatal(v ...interface{}) {
	s := fmt.Sprint(v...)
	FatalLogger.Output(2, s)
	panic(s)
}

func Fatalf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	FatalLogger.Output(2, s)
	panic(s)
}

/*定义日志实现*/
type Logger struct {
	Name  string
	Path  string
	debug bool
	*os.File
}

func NewLogger(name, path string, debug bool) *Logger {
	return &Logger{Name: name, Path: path, debug: debug}
}

func (lg *Logger) Debug(v ...interface{}) {
	if lg.debug {
		lg.commonLog("DEBUG", v...)
	}
}

func (lg *Logger) Debugf(format string, v ...interface{}) {
	if lg.debug {
		lg.commonLogf("DEBUG", format, v...)
	}
}

func (lg *Logger) Info(v ...interface{}) {
	lg.commonLog("INFO", v...)
}

func (lg *Logger) Infof(format string, v ...interface{}) {
	lg.commonLogf("INFO", format, v...)
}

func (lg *Logger) Warn(v ...interface{}) {
	lg.commonLog("WARN", v...)
}

func (lg *Logger) Warnf(format string, v ...interface{}) {
	lg.commonLogf("WARN", format, v...)
}

func (lg *Logger) Error(v ...interface{}) {
	lg.commonLog("ERROR", v...)
}

func (lg *Logger) Errorf(format string, v ...interface{}) {
	lg.commonLogf("ERROR", format, v...)
}

func (lg *Logger) Fatal(v ...interface{}) {
	s := fmt.Sprint(v...)
	lg.Error(v...)
	panic(s)
}

func (lg *Logger) Fatalf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	lg.Errorf(format, v...)
	panic(s)
}

func (lg *Logger) commonLog(logType string, v ...interface{}) {
	logStr := fmt.Sprintln(v...)
	logStr = fmt.Sprintf("[%s] %s", logType, logStr)
	writeLog(logStr)
}

func (lg *Logger) commonLogf(logType string, format string, v ...interface{}) {
	logStr := fmt.Sprintf(format, v...)
	logStr = fmt.Sprintf("[%s] %s", logType, logStr)
	writeLog(logStr)
}

func writeLog(logStr string) (err error) {
	err = resetOutputIfNeed()
	if err != nil {
		fmt.Println("writeLog ERROR, resetOutput: ", err)
		return
	}
	var (
		fn  string
		lno int
	)
	_, fn, lno, ok := runtime.Caller(3)
	if ok {
		fntmp := fn
		for i := len(fn) - 1; i > 0; i-- {
			if fn[i] == '/' {
				fntmp = fntmp[i+1:]
				break
			}
		}
		fn = fntmp
	} else {
		fn = "???"
		lno = 0
	}
	logStr = fmt.Sprintf("%s:%d %s", fn, lno, logStr)
	log.Print(logStr)
	return
}

func resetOutputIfNeed() (err error) { //手动分割日志文件，切换输出到新文件
	needReset := false
	logFilePath := getLogFilePath()
	if Log.File == nil {
		needReset = true
	} else {
		fileInfo, statErr := Log.File.Stat()
		if statErr != nil {
			needReset = true
		} else if !strings.HasSuffix(logFilePath, fileInfo.Name()) {
			needReset = true
		}
	}
	if needReset {
		oldFile := Log.File
		newFile, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		defer closeFile(oldFile)
		Log.File = newFile
		var output io.Writer
		if Log.debug {
			output = io.MultiWriter(os.Stderr, Log.File)
		} else {
			output = Log.File
		}
		log.SetOutput(output)
	}
	return
}

func getLogFilePath() (path string) {
	if Log == nil {
		return
	}
	if len(Log.Path) == 0 {
		return
	}
	tmps := strings.Split(Log.Path, ".")
	if tmps[len(tmps)-1] == "log" {
		tmps = tmps[:len(tmps)-1]
	}
	dateStr := time.Now().Format("20060102")
	tmps = append(tmps, []string{dateStr, "log"}...)
	return strings.Join(tmps, ".")
}

func closeFile(file *os.File) {
	if file == nil {
		return
	}
	err := file.Close()
	if err != nil {
		log.Println("ERROR AT CLOSING LOG FILE: ", err.Error())
	}
	return
}

func Init(name, logFilePath string, debug bool) {
	log.SetFlags(log.Ldate | log.Ltime)
	Log = NewLogger(name, logFilePath, debug)
}
