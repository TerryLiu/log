package main

import (
	"github.com/terryliu/log/v2"
	"go.uber.org/zap"
)

func main() {
	// Init用来初始化设置
	// 参数1:设置日志路径
	// 参数2:日志输出的级别
	// 参数3:是否需要请求类型的日志(它为info级), 需要的话,所有的请求都会打在.Request扩展名的文件中.
	// 注意: 日志分为两种,1是普通日志,格式默认为json，2是请求类型日志,格式默认为csv.
	// 参数4:是否按等级分开打印普通日志; 说明: 各等级分开打印时,日志扩展名被分成4类文件：.DEBUG，.INFO, .WARN, .ERROR; 其中error,panic,Fatal都会打印在.ERROR中.
	// 参数5:从这个参数开始,后面的参数都是LogOption结构的回调函数,比如下面的例子为设置是否显示Caller,还可以使用log.SetCallerDeep来设置深度；用log.SetLogType("csv")可将普通日志格式设置为csv, 用log.SetRequestType("json")可将请求日志设置为json格式；
	// 参数6:用log.SetMaxFileSize(2)举例,设置日志大小为2MB,默认是500MB.顺便提一句,最大文件个数默认是5个,最大文件保留天数是7天.
	// 参数7:参考5....
	log.Init("./test.log", log.DebugLevel, true, false, log.SetCaller(true), log.SetLogType("csv"), log.SetMaxFileSize(2), log.SetCompress(true))

	// 输出info级别的日志
	log.Debug("哈喽,艾瑞巴蒂!我是带体积限制功能的滚动异步日志套装.基于Zap和lumberjack.")
	log.Infof("哈喽,艾瑞巴蒂%s!我是带体积限制功能的滚动异步日志套装.基于Zap和lumberjack.", "-待插入的变量值-")
	log.Warnf("哈喽,艾瑞巴蒂%s!我是带体积限制功能的滚动异步日志套装.基于Zap和lumberjack.", "-待插入的变量值-")
	log.Errorw("输出消息名称，后面的键值对代表要打印的数据", zap.String("key00000000", "val11111"), zap.String("k111111", "v1111111111"))

	// 请求数据的日志文件
	log.RequestLogInfo([]string{"key,0000000000", "val,00000000", "k111111", "v1111111111"})
	log.RequestLogInfo("{\"column\": \"002004003\", \"industryCode\": \"\",\n    \"industryName\": \"\",\n    \"emIndustryCode\": \"\",\n    \"indvInduCode\": \"481\"}")
	log.RequestLogInfof("哈喽,艾瑞巴蒂%s,我是带体积限制功能的滚动异步日志套装.基于Zap和lumberjack.", "-待插入的变量值-")
	log.RequestLogInfow("哈喽,这里演示csv的内容中包含逗号的效果", zap.String("key,00000000", "val11111"), zap.String("k111111", "v11111,11111"))

	// 最后将未输出的日志缓存,就全部输出.不调用只会异步输出,不能保证输出的完整性
	log.Sync()

}
