package log

import (
	"fmt"
)

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	PanicLevel = "panic"
)
const (
	FileTypeLog = iota
	FileTypeRequest
	DebugLevelLog
	InfoLevelLog
	WarnLevelLog
	ErrorLevelLog
	PanicLevelLog
)

var logger *Log

// Log 默认会使用zap作为日志输出引擎. Log集成了日志切割的功能。默认文件大小1024M，自动压缩
// 最大有3个文件备份，备份保存时间7天。默认不会打印日志被调用的文文件名和位置;
// 输出:日志默认会被分成两类文件：.log, .log.Request;可以通过SetLogType和SetRequestType函数修改两类文件的存储格式
// debug,info,warn,error,panic都会打印在xxx.log. 所有的请求都会打在xxx.log.Request
// Adapter:经过比对现在流行的日志库：zap, logrus, zerolog; logrus 虽说格式化，插件化良好，但是
// 其内部实现锁竞争太过剧烈，性能不好. zap 性能好，格式一般， zerolog性能没有zap好， 相比
// 来说就没啥突出优点了

type Log struct {
	Path           string
	Level          string
	NeedRequestLog bool // 是否需要独立的Request日志
	NeedLevelsLog  bool // 是否需要各个等级的日志分开打印
	adapters       []*zapAdapter
}

type LogOption interface {
	apply(*Log)
}

type logOptionFunc func(*Log)

func (f logOptionFunc) apply(log *Log) {
	f(log)
}

// contentType=json;csv
func SetLogType(contentType string) LogOption {
	return logOptionFunc(func(log *Log) {
		for k, _ := range log.adapters {
			if k != FileTypeRequest {
				log.adapters[k].setLogType(contentType)
			}
		}
	})
}

// contentType=json;csv
func SetRequestType(contentType string) LogOption {
	return logOptionFunc(func(log *Log) {
		for k, _ := range log.adapters {
			if k == FileTypeRequest {
				log.adapters[k].setLogType(contentType)
			}
		}
	})
}

func SetMaxFileSize(size int) LogOption {
	return logOptionFunc(func(log *Log) {
		for i, _ := range log.adapters {
			log.adapters[i].setMaxFileSize(size)
		}
	})
}

func SetMaxBackups(n int) LogOption {
	return logOptionFunc(func(log *Log) {
		for i, _ := range log.adapters {
			log.adapters[i].setMaxBackups(n)
		}
	})
}

func SetMaxAge(age int) LogOption {
	return logOptionFunc(func(log *Log) {
		for i, _ := range log.adapters {
			log.adapters[i].setMaxAge(age)
		}
	})
}

func SetCompress(compress bool) LogOption {
	return logOptionFunc(func(log *Log) {
		for i, _ := range log.adapters {
			log.adapters[i].setCompress(compress)
		}
	})
}

func SetCaller(caller bool) LogOption {
	return logOptionFunc(func(log *Log) {
		for i, _ := range log.adapters {
			log.adapters[i].setCaller(caller)
		}
	})
}
func SetCallerDeep(callerDeep int) LogOption {
	return logOptionFunc(func(log *Log) {
		for i, _ := range log.adapters {
			log.adapters[i].setCallerDeep(callerDeep)
		}
	})
}

// Init init logger
func Init(path, level string, needRequestLog, needLevelsLog bool, options ...LogOption) {
	logger = &Log{Path: path, Level: level}
	logger.createFiles(level, needRequestLog, needLevelsLog, options...)
}

// Sync flushes buffer, if any
func Sync() {
	if logger == nil {
		return
	}

	for _, v := range logger.adapters {
		v.logger.Sync()
	}
}

//
// func (l *Log) maxFileSize(fileType int) int {
// 	if fileType==FileTypeLog || fileType==FileTypeRequest {
// 		return l.adapters[logType].MaxFileSize
// 	}
// 	return 0
// }
//

func (l *Log) createFiles(level string, needRequestLog, needLevelsLog bool, options ...LogOption) {
	adapters := make([]*zapAdapter, 6)
	// adapters := make(map[string]*zapAdapter, 2)
	adapters[FileTypeLog] = NewZapAdapter(fmt.Sprintf("%s", l.Path), level, "json")
	adapters[FileTypeRequest] = NewZapAdapter(fmt.Sprintf("%s.Request", l.Path), InfoLevel, "csv")
	adapters[DebugLevelLog] = NewZapAdapter(fmt.Sprintf("%s.DEBUG", l.Path), DebugLevel, "json")
	adapters[InfoLevelLog] = NewZapAdapter(fmt.Sprintf("%s.INFO", l.Path), InfoLevel, "json")
	adapters[WarnLevelLog] = NewZapAdapter(fmt.Sprintf("%s.WARN", l.Path), WarnLevel, "json")
	adapters[ErrorLevelLog] = NewZapAdapter(fmt.Sprintf("%s.ERROR", l.Path), ErrorLevel, "json")
	l.NeedRequestLog = needRequestLog
	l.NeedLevelsLog = needLevelsLog
	l.adapters = adapters

	// options为回调函数,用来作为log对象的中间件进行调用
	for _, opt := range options {
		// 将log对象作为参数传入回调函数中
		opt.apply(l)
	}

	for _, adapter := range adapters {
		adapter.Init()
	}

}

// 判断是否需要独立打印不同等级的日志
func needLevelsLog(selfLevel int) int {
	if logger == nil {
		return 0
	}
	if logger.NeedLevelsLog {
		if PanicLevelLog == selfLevel {
			return ErrorLevelLog
		}
		return selfLevel
	}
	if selfLevel >= DebugLevelLog {
		return FileTypeLog
	}
	return FileTypeRequest
}

// Debug 使用方法：log.Debug("test")
func Debug(args ...interface{}) {
	if logger == nil {
		return
	}
	logger.adapters[needLevelsLog(DebugLevelLog)].Debug(args...)
}

// Debugf 使用方法：log.Debugf("test:%s", err)
func Debugf(template string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.adapters[needLevelsLog(DebugLevelLog)].Debugf(template, args...)
}

// Debugw 使用方法：log.Debugw("test", "field1", "value1", "field2", "value2")
func Debugw(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}

	logger.adapters[needLevelsLog(DebugLevelLog)].Debugw(msg, keysAndValues...)
}

func Info(args ...interface{}) {
	if logger == nil {
		return
	}

	logger.adapters[needLevelsLog(InfoLevelLog)].Info(args...)
}

func Infof(template string, args ...interface{}) {
	if logger == nil {
		return
	}

	logger.adapters[needLevelsLog(InfoLevelLog)].Infof(template, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}

	logger.adapters[needLevelsLog(InfoLevelLog)].Infow(msg, keysAndValues...)
}

func Output(calldepth int, s string) error {
	Info(s)
	return nil
}
func Println(v ...interface{}) {
	Info(v)
}
func Printf(format string, v ...interface{}) {
	Infof(format, v)
}

func Warn(args ...interface{}) {
	if logger == nil {
		return
	}

	logger.adapters[needLevelsLog(WarnLevelLog)].Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	if logger == nil {
		return
	}

	logger.adapters[needLevelsLog(WarnLevelLog)].Warnf(template, args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}

	logger.adapters[needLevelsLog(WarnLevelLog)].Warnw(msg, keysAndValues...)
}

func Error(args ...interface{}) {
	if logger == nil {
		return
	}

	logger.adapters[needLevelsLog(ErrorLevelLog)].Error(args...)
}

func Errorf(template string, args ...interface{}) {
	if logger == nil {
		return
	}

	logger.adapters[needLevelsLog(ErrorLevelLog)].Errorf(template, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}

	logger.adapters[needLevelsLog(ErrorLevelLog)].Errorw(msg, keysAndValues...)
}

func Panic(args ...interface{}) {
	if logger == nil {
		return
	}

	logger.adapters[needLevelsLog(PanicLevelLog)].Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	if logger == nil {
		return
	}

	logger.adapters[needLevelsLog(PanicLevelLog)].Panicf(template, args...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}

	logger.adapters[needLevelsLog(PanicLevelLog)].Panicw(msg, keysAndValues...)
}

func Fatal(args ...interface{}) {
	if logger == nil {
		return
	}

	logger.adapters[needLevelsLog(PanicLevelLog)].Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	if logger == nil {
		return
	}

	logger.adapters[needLevelsLog(PanicLevelLog)].Fatalf(template, args...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}

	logger.adapters[needLevelsLog(PanicLevelLog)].Fatalw(msg, keysAndValues...)
}

// 参数keysAndValues为一个切片,元素1为key,元素2为val;以此类推.
func RequestLogInfo(keysAndValues ...interface{}) {
	if logger == nil || !logger.NeedRequestLog {
		return
	}
	logger.adapters[needLevelsLog(FileTypeRequest)].Info(keysAndValues...)
}
func RequestLogInfof(template string, args ...interface{}) {
	if logger == nil || !logger.NeedRequestLog {
		return
	}
	logger.adapters[needLevelsLog(FileTypeRequest)].Infof(template, args...)
}
func RequestLogInfow(template string, keysAndValues ...interface{}) {
	if logger == nil || !logger.NeedRequestLog {
		return
	}
	logger.adapters[needLevelsLog(FileTypeRequest)].Infow(template, keysAndValues...)
}
