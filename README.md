# log

base on zap, an easy using logger module

```go
package main

import "github.com/terryliu/log"

func main() {
	// Init用来初始化设置
	// 参数1:设置日志路径
	// 参数2:日志输出的级别
	// 参数3:是否需要另存请求类型的日志
	// 参数4:从这个参数开始就都是回调函数,可以用来设置是否显示Caller
	// 参数5:这里用log.SetMaxFileSize(1)举例,设置日志大小为1MB,默认是1024MB
	log.Init("./test.log", log.DebugLevel, true, log.SetCaller(true), log.SetMaxFileSize(1))
	log.Info("hello 中文日志.")
	//请求数据的日志文件
	log.RequestLogInfow([]string{"key000","val000","k001","v001"})
	// 如果有未输出的日志缓存,就刷新输出,不调用只会异步输出,不保证输出的完整性
	log.Sync()
	/*
	输出:
	{"level":"info","ts":"2020-04-10T22:28:02.904+0800","caller":"example/example.go:15","msg":"hello 中文日志."}
	{"level":"info","ts":"2020-04-10T22:28:02.955+0800","caller":"example/example.go:17","msg":"[key000 val000 k001 v001]"}
	*/
}
```