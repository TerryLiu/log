package log

import (
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type zapAdapter struct {
	Path        string // 文件绝对地址，如：/home/homework/neso/file.log
	Level       string // 日志输出的级别
	LogType     string // 日志格式类型.支持:json;csv
	MaxFileSize int    // 日志文件大小的最大值，单位(M)
	MaxBackups  int    // 最多保留备份数
	MaxAge      int    // 日志文件保存的时间，单位(天)
	Compress    bool   // 是否压缩
	Caller      bool   // 日志是否需要显示调用位置
	CallerDeep  int    // 调用文件回显的深度

	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

func (z *zapAdapter) setLogType(logType string) {
	z.LogType = logType
}

func (z *zapAdapter) setMaxFileSize(size int) {
	z.MaxFileSize = size
}

func (z *zapAdapter) setMaxBackups(n int) {
	z.MaxBackups = n
}

func (z *zapAdapter) setMaxAge(age int) {
	z.MaxAge = age
}

func (z *zapAdapter) setCompress(compress bool) {
	z.Compress = compress
}

func (z *zapAdapter) setCaller(caller bool) {
	z.Caller = caller
}
func (z *zapAdapter) setCallerDeep(callerDeep int) {
	z.CallerDeep = callerDeep
}
func EnsureCSVSuffix(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))

	if ext != ".csv" {
		// filePath = strings.TrimSuffix(filePath, ext)
		filePath = filePath + ".csv"
	}

	return filePath
}
func NewZapAdapter(path, level, contentType string) *zapAdapter {
	return &zapAdapter{
		Path:        path,
		Level:       level,
		LogType:     contentType,
		MaxFileSize: 500,
		MaxBackups:  5,
		MaxAge:      7,
		Compress:    true,
		Caller:      false,
		CallerDeep:  2,
	}
}

// createLumberjackHook 创建LumberjackHook，其作用是为了将日志文件切割，压缩
func (zapAdapter *zapAdapter) createLumberjackHook() *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   zapAdapter.Path,
		MaxSize:    zapAdapter.MaxFileSize,
		MaxBackups: zapAdapter.MaxBackups,
		MaxAge:     zapAdapter.MaxAge,
		Compress:   zapAdapter.Compress,
		LocalTime:  true,
	}
}

// 根据配置的参数来初始化日志对象
func (zapAdapter *zapAdapter) Init() {
	if zapAdapter.LogType == "csv" {
		zapAdapter.Path = EnsureCSVSuffix(zapAdapter.Path)
	}
	w := zapcore.AddSync(zapAdapter.createLumberjackHook())

	var level zapcore.Level
	var cnf zapcore.Encoder

	switch zapAdapter.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "panic":
		level = zap.PanicLevel
	default:
		level = zap.InfoLevel
	}

	conf := zap.NewProductionEncoderConfig()
	conf.EncodeTime = zapcore.ISO8601TimeEncoder
	// 除非指定了csv, 否则默认使用json内容格式
	switch zapAdapter.LogType {
	case "csv":
		cnf = NewCSVEncoder(conf)
	default:
		cnf = zapcore.NewJSONEncoder(conf)
	}

	core := zapcore.NewCore(cnf, w, level)
	zapAdapter.logger = zap.New(core)
	if zapAdapter.Caller {
		zapAdapter.logger = zapAdapter.logger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(zapAdapter.CallerDeep))
	}
	zapAdapter.sugar = zapAdapter.logger.Sugar()
}

func (zapAdapter *zapAdapter) Debug(args ...interface{}) {
	zapAdapter.sugar.Debug(args...)
}

func (zapAdapter *zapAdapter) Info(args ...interface{}) {
	zapAdapter.sugar.Info(args...)
}

func (zapAdapter *zapAdapter) Warn(args ...interface{}) {
	zapAdapter.sugar.Warn(args...)
}

func (zapAdapter *zapAdapter) Error(args ...interface{}) {
	zapAdapter.sugar.Error(args...)
}

func (zapAdapter *zapAdapter) DPanic(args ...interface{}) {
	zapAdapter.sugar.DPanic(args...)
}

func (zapAdapter *zapAdapter) Panic(args ...interface{}) {
	zapAdapter.sugar.Panic(args...)
}

func (zapAdapter *zapAdapter) Fatal(args ...interface{}) {
	zapAdapter.sugar.Fatal(args...)
}

func (zapAdapter *zapAdapter) Debugf(template string, args ...interface{}) {
	zapAdapter.sugar.Debugf(template, args...)
}

func (zapAdapter *zapAdapter) Infof(template string, args ...interface{}) {
	zapAdapter.sugar.Infof(template, args...)
}

func (zapAdapter *zapAdapter) Warnf(template string, args ...interface{}) {
	zapAdapter.sugar.Warnf(template, args...)
}

func (zapAdapter *zapAdapter) Errorf(template string, args ...interface{}) {
	zapAdapter.sugar.Errorf(template, args...)
}

func (zapAdapter *zapAdapter) DPanicf(template string, args ...interface{}) {
	zapAdapter.sugar.DPanicf(template, args...)
}

func (zapAdapter *zapAdapter) Panicf(template string, args ...interface{}) {
	zapAdapter.sugar.Panicf(template, args...)
}

func (zapAdapter *zapAdapter) Fatalf(template string, args ...interface{}) {
	zapAdapter.sugar.Fatalf(template, args...)
}

func (zapAdapter *zapAdapter) Debugw(msg string, keysAndValues ...interface{}) {
	zapAdapter.sugar.Debugw(msg, keysAndValues...)
}

func (zapAdapter *zapAdapter) Infow(msg string, keysAndValues ...interface{}) {
	zapAdapter.sugar.Infow(msg, keysAndValues...)
}

func (zapAdapter *zapAdapter) Warnw(msg string, keysAndValues ...interface{}) {
	zapAdapter.sugar.Warnw(msg, keysAndValues...)
}

func (zapAdapter *zapAdapter) Errorw(msg string, keysAndValues ...interface{}) {
	zapAdapter.sugar.Errorw(msg, keysAndValues...)
}

func (zapAdapter *zapAdapter) DPanicw(msg string, keysAndValues ...interface{}) {
	zapAdapter.sugar.DPanicw(msg, keysAndValues...)
}

func (zapAdapter *zapAdapter) Panicw(msg string, keysAndValues ...interface{}) {
	zapAdapter.sugar.Panicw(msg, keysAndValues...)
}

func (zapAdapter *zapAdapter) Fatalw(msg string, keysAndValues ...interface{}) {
	zapAdapter.sugar.Fatalw(msg, keysAndValues...)
}
