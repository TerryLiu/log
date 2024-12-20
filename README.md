# log

base on zap&lumberjack, an easy using logger module

```go

func main() {
    // Init用来初始化设置
    // 参数1:设置日志路径
    // 参数2:日志输出的级别
    // 参数3:是否需要请求类型的日志；默认情况下普通日志格式为json，请求类型日志格式为csv；
    // 参数4:从这个参数开始,后面的参数都是LogOption结构的回调函数,比如下面的例子为设置是否显示Caller,还可以使用log.SetCallerDeep来设置深度；用log.SetLogType("csv")可将普通日志格式设置为csv, 用log.SetRequestType("json")可将请求日志设置为json格式；
    // 参数5:用log.SetMaxFileSize(2)举例,设置日志大小为2MB,默认是500MB.顺便提一句,最大文件个数默认是5个,最大文件保留天数是7天.
    // 参数6:参考4....
    log.Init("./test.log", log.DebugLevel, true, log.SetCaller(true), log.SetMaxFileSize(2), log.SetCompress(true))
    
    // 输出info级别的日志
    log.Info("哈喽,艾瑞巴蒂!我是带体积限制功能的滚动异步日志套装.基于Zap和lumberjack.")
    log.Infof("哈喽,艾瑞巴蒂%s!我是带体积限制功能的滚动异步日志套装.基于Zap和lumberjack.", "-待插入的变量值-")
    log.Infow("输出消息名称，后面的键值对代表要打印的数据", zap.String("key00000000", "val11111"), zap.String("k111111", "v1111111111"))
    
    // 请求数据的日志文件
    log.RequestLogInfo([]string{"key,0000000000", "val,00000000", "k111111", "v1111111111"})
    log.RequestLogInfof("哈喽,艾瑞巴蒂%s,我是带体积限制功能的滚动异步日志套装.基于Zap和lumberjack.", "-待插入的变量值-")
    log.RequestLogInfow("哈喽,这里演示csv的内容中包含逗号的效果", zap.String("key,00000000", "val11111"), zap.String("k111111", "v11111,11111"))
    
    // 如果有未输出的日志缓存,就全部输出.不调用只会异步输出,不能保证输出的完整性
    log.Sync()

}

```