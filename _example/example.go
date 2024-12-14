package main

import "github.com/TerryLiu/log/v2"

func main() {
	// Init用来初始化设置
	// 参数1:设置日志路径
	// 参数2:日志输出的级别
	// 参数3:是否需要另存请求类型的日志
	// 参数4:从这个参数开始,后面的参数都是LogOption结构的回调函数,比如下面的例子为设置是否显示Caller,还可以使用log.SetCallerDeep来设置深度
	// 参数5:用log.SetMaxFileSize(2)举例,设置日志大小为2MB,默认是500MB.顺便提一句,最大文件个数默认是5个,最大文件保留天数是7天.
	// 参数6:参考4....
	log.Init("./test.log", log.DebugLevel, true, log.SetCaller(true),log.SetMaxFileSize(2))

	// 输出info级别的日志
	log.Info("哈喽,艾瑞巴蒂!我是带体积限制功能的滚动异步日志套装.基于Zap和lumberjack.")
	// 请求数据的日志文件
	log.RequestLogInfow([]string{"key0000000000", "val00000000", "k111111", "v1111111111"})

	// 如果有未输出的日志缓存,就全部输出.不调用只会异步输出,不能保证输出的完整性
	log.Sync()
	/*
	输出:
	{"level":"info","ts":"2020-04-10T22:28:02.904+0800","caller":"example/example.go:15","msg":"哈喽,艾瑞巴蒂!我是带体积限制功能的滚动异步日志套装.基于Zap和lumberjack."}
	{"level":"info","ts":"2020-04-10T22:28:02.955+0800","caller":"example/example.go:17","msg":"[key0000000000 val00000000 k111111 v1111111111]"}
	*/
}
