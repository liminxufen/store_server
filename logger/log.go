package logger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"golang.org/x/crypto/ssh/terminal"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

// Define log level string
const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warning"
	ErrorLevel = "error"

	fieldTraceback   = "traceback"
	fieldUser        = "user"
	fieldMessageJSON = "message_json"
	fieldFile        = "file"
	fieldLineNo      = "lineno"
	fieldEvent       = "func"
	fieldService     = "service"
	fieldIP          = "host"
	fieldPort        = "port"
	fieldProcess     = "process"
	fieldTraceID     = "traceId"
)

const (
	maximumCallerDepth int = 25
	knownLogFrames     int = 4
)

var (
	minimumCallerDepth = 1
	// Used for caller information initialisation
	callerInitOnce sync.Once
	// qualified package name, cached at first use
	logPackage string

	logger      *JLogger
	LoggerEntry *JLoggerEntry
)

const (
	red    = 31
	green  = 32
	yellow = 33
	blue   = 34
	cyan   = 36
)

func InitStructLog(logLevel, logPath, service string, hooks ...logrus.Hook) (err error) {
	logger, err = NewJLogger(logLevel, logPath, service, hooks...)
	if err != nil {
		return err
	}
	return nil
}

func Entry() (je *JLoggerEntry) {
	return logger.Entry()
}

func colored(level logrus.Level, line string) string {
	var levelColor int
	switch level {
	case logrus.DebugLevel:
		levelColor = cyan
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel:
		levelColor = red
	default:
		levelColor = blue
	}
	return fmt.Sprintf("\033[1;%dm%s\033[0m", levelColor, line)
}

type jsonTextField string

func (f jsonTextField) MarshalJSON() (data []byte, err error) {
	s := string(f)
	if s == "" {
		data = []byte("null")
	} else {
		data = []byte(s)
	}
	return
}

type line struct {
	Time        string        `json:"time"`
	Level       string        `json:"level"`
	IP          string        `json:"host"`
	Port        string        `json:"port"`
	Service     string        `json:"service"`
	Process     int           `json:"process"`
	TraceID     string        `json:"traceId"`
	File        string        `json:"file"`
	LineNo      int           `json:"lineno"`
	Event       string        `json:"func"`
	MessageText string        `json:"message"`
	MessageJSON jsonTextField `json:"message_json,omitempty"`
	User        string        `json:"user,omitempty"`
	Traceback   string        `json:"traceback,omitempty"`
}

func (l *line) reset() {
	l.Traceback = ""
	l.TraceID = ""
	l.MessageJSON = ""
	l.User = ""
}

type stdFormatter struct {
	colored  bool
	linePool sync.Pool
}

func (f *stdFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	line := f.linePool.Get().(*line)
	defer func() {
		line.reset()
		f.linePool.Put(line)
	}()

	line.Level = entry.Level.String()
	line.Time = entry.Time.Format(time.RFC3339)

	line.Service = entry.Data[fieldService].(string)
	line.IP = entry.Data[fieldIP].(string)

	line.File = entry.Data[fieldFile].(string)
	line.LineNo = entry.Data[fieldLineNo].(int)
	line.Event = entry.Data[fieldEvent].(string)

	if t, ok := entry.Data[fieldTraceID]; ok {
		line.TraceID = t.(string)
	}

	if t, ok := entry.Data[fieldProcess]; ok {
		line.Process = t.(int)
	}

	if t, ok := entry.Data[fieldTraceback]; ok {
		line.Traceback = t.(string)
	}

	if t, ok := entry.Data[fieldMessageJSON]; ok {
		v := t.(map[string]interface{})
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		line.MessageJSON = jsonTextField(b)
	}
	if t, ok := entry.Data[fieldUser]; ok {
		line.User = t.(string)
	}

	line.MessageText = entry.Message
	b, err := json.Marshal(line)
	b = append(b, '\n')

	if f.colored {
		bs := string(b)
		b = []byte(colored(entry.Level, bs))
	}

	return b, err
}

type multiWriter struct {
	writers []io.Writer
}

//创建
func newMultiWriter(writers ...io.Writer) *multiWriter {
	wrs := make([]io.Writer, 0, len(writers))
	wrs = append(wrs, writers...)
	return &multiWriter{
		writers: wrs,
	}
}

func (t *multiWriter) Write(p []byte) (n int, err error) {
	for _, w := range t.writers {
		n, err = w.Write(p)
		if err != nil {
			continue
		}
		if n != len(p) {
			err = errors.New("error in short write")
		}
	}
	return len(p), nil
}

type JLogger struct {
	logger           *logrus.Logger
	jLoggerEntryPool sync.Pool
	service          string
	ip               string
	Port             int
	Process          int
	ctx              context.Context
}

func (jlogger *JLogger) Entry() (je *JLoggerEntry) {
	return newJLoggerEntry(jlogger)
}

func NewJLogger(logLevel string, logPath string, service string, hooks ...logrus.Hook) (jlogger *JLogger, err error) {
	jlogger = new(JLogger)
	jlogger.service = service
	jlogger.ip = localIP()
	jlogger.Process = os.Getpid()

	var level logrus.Level
	if level, err = logrus.ParseLevel(logLevel); err != nil {
		return nil, err
	}
	dir := filepath.Dir(logPath)
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}
	logger := logrus.New()
	logger.Level = level

	format := &stdFormatter{}
	format.linePool.New = func() interface{} {
		return new(line)
	}
	if isTerminal() {
		os.Stderr = os.Stdout
		logger.Out = os.Stdout

		format.colored = true
		logger.Formatter = format
	} else {
		format.colored = false
		logger.Formatter = format

		writer, err := rotatelogs.New(logPath + "-%Y%m%d")
		if err != nil {
			return nil, err
		}
		logger.Out = writer
	}
	for _, hook := range hooks {
		logger.AddHook(hook)
	}
	jlogger.logger = logger
	return jlogger, nil
}

type JLoggerEntry struct {
	entry     *logrus.Entry
	jsonField map[string]interface{}
	mu        sync.Mutex
}

func newJLoggerEntry(jLogger *JLogger) (je *JLoggerEntry) {
	je = new(JLoggerEntry)
	je.jsonField = make(map[string]interface{})
	je.entry = jLogger.logger.WithField(fieldService, jLogger.service).
		WithField(fieldIP, jLogger.ip).WithField(fieldProcess, jLogger.Process)
	return je
}

func (je *JLoggerEntry) WithError(err error) (j *JLoggerEntry) {
	je.entry = je.entry.WithField(fieldTraceback, formatErrorStack(err.Error()))
	return je
}

func (je *JLoggerEntry) WithUser(user string) (j *JLoggerEntry) {
	je.entry = je.entry.WithField(fieldUser, user)
	return je
}

func (je *JLoggerEntry) WithContext(ctx context.Context) (j *JLoggerEntry) {
	spanCtx := opentracing.SpanFromContext(ctx)
	if spanCtx != nil {
		sctx, ok := spanCtx.Context().(jaeger.SpanContext)
		if ok {
			je.entry = je.entry.WithField(fieldTraceID, sctx.TraceID().String())
		} else {
			je.entry = je.entry.WithField(fieldTraceID, genTraceId())
		}
	}
	return je
}

func (je *JLoggerEntry) WithMessageJSON(m map[string]interface{}) (j *JLoggerEntry) {
	je.mu.Lock()
	for k, v := range m {
		je.jsonField[k] = v
	}
	je.mu.Unlock()
	je.entry = je.entry.WithField(fieldMessageJSON, je.jsonField)
	return je
}

func (je *JLoggerEntry) withCaller() (j *logrus.Entry) {
	file, lineNo, function := getCaller()
	je.entry = je.entry.WithField(fieldFile, file).
		WithField(fieldLineNo, lineNo).
		WithField(fieldEvent, function)
	return je.entry
}

func (je *JLoggerEntry) Debug(args ...interface{}) {
	je.withCaller().Debug(args...)
}

func (je *JLoggerEntry) Debugf(format string, args ...interface{}) {
	je.withCaller().Debugf(format, args...)
}

func (je *JLoggerEntry) Info(args ...interface{}) {
	je.withCaller().Info(args...)
}

func (je *JLoggerEntry) Infof(format string, args ...interface{}) {
	je.withCaller().Infof(format, args...)
}

func (je *JLoggerEntry) Warn(args ...interface{}) {
	je.withCaller().Warn(args...)
}

func (je *JLoggerEntry) Warnf(format string, args ...interface{}) {
	je.withCaller().Warnf(format, args...)
}

func (je *JLoggerEntry) Error(args ...interface{}) {
	je.withCaller().Error(args...)
}

func (je *JLoggerEntry) Errorf(format string, args ...interface{}) {
	je.withCaller().Errorf(format, args...)
}

// getCaller retrieves the name of the first non-log package calling function
func getCaller() (file string, line int, function string) {
	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		logPackage = getPackageName(runtime.FuncForPC(pcs[0]).Name())

		// now that we have the cache, we can skip a minimum count of known-log functions
		// XXX this is dubious, the number of frames may vary store an entry in a logger interface
		minimumCallerDepth = knownLogFrames
	})

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if pkg != logPackage {
			return f.File, f.Line, f.Function
		}
	}

	// if we got here, we failed to find the caller's context
	return "???", -1, "???"
}

// getPackageName reduces a fully qualified function name to the packagname
// There really ought to be to be a better way...
func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}

// FormatErrorStack  Remove  project name  attached by errors.ErrorStack,
func formatErrorStack(format string) string {
	if lines := strings.Split(format, "\n"); len(lines) > 0 {
		var newLines []string
		for _, i := range lines {
			i = strings.Replace(i, "\t", "\\t", -1)
			newLines = append(newLines, i)
		}
		return strings.Join(newLines, "\\n")
	}
	return format
}

func isTerminal() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd())) || terminal.IsTerminal(int(os.Stderr.Fd()))
}

func localIP() string {
	tables, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, t := range tables {
		addrs, err := t.Addrs()
		if err != nil {
			return ""
		}
		for _, a := range addrs {
			ipnet, ok := a.(*net.IPNet)
			if !ok || ipnet.IP.IsLoopback() {
				continue
			}
			if v4 := ipnet.IP.To4(); v4 != nil {
				return v4.String()
			}
		}
	}
	return ""
}

var incr = 1000

func genTraceId() string {
	ip := localIP()
	ipslice := strings.Split(ip, ".")
	ipIs := make([]interface{}, len(ipslice))
	for _, is := range ipslice {
		in, _ := strconv.Atoi(is)
		ipIs = append(ipIs, in)
	}
	formatStr := strings.Repeat("%x", len(ipIs))
	ipX := fmt.Sprintf(formatStr, ipIs...)
	ts := int64(time.Now().Unix() * 1000)
	pid := os.Getpid()
	if incr > 9999 {
		incr = 1000
	}
	//traceId component by ip, timestamp, autoincrement, process id
	traceId := fmt.Sprintf("%s%v%v%v", ipX, ts, incr, pid)
	incr++
	return traceId
}
