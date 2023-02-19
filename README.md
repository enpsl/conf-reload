# conf-reload


Conf-Reload是一个Go库，用于常驻进程模式下配置文件信息获取，并支持配置热加载

工作流程概述：

- 读取配置文件信息并基于`Broker.unmarshaller`进行解析

- `Broker`会单独启动一个`Goroutine`和`fsnotify`进行通信

- 当`fsnotify`有文件创建或写事件进来时,`Broker`通知`Engine`重载配置文件

# Features
- 支持`toml` `yaml` `json`配置文件获取
- 支持key多层级获取
- 支持热更新，配置文件更新后，配置会重载
- 基于`mapstructure`库实现，配置信息获取支持弱类型转化
- 基于`fsnotify`库实现了IO多路复用的事件通知机制，性能较好
- 配置结构获取提供多种类型转化api

# Stability and Compatibility

**Status**: 目前处于功能拓展中，不过基础功能已经完善，不会在改变api结构

> ☝️ **Important Note**: v1.0.0之前为beta版本. v1.0.0版本为当前stable版本.


# Quick Start

```bash
go get -u github.com/enpsl/conf-reload
```

## _example
_example目录是测试用例，可`copy`到当前项目下测试运行

配置引擎加载,加载一次即可
```go
f = "_example/example.toml"
conf_reload.LoadEngine(f, conf_relod.WithLevelSplit("."), conf_relod.WithLogLevel(0))
```
LoadEngine的一些[option](https://pkg.go.dev/github.com/enpsl/conf-reload#Option)选项说明:
- `WithLevelSplit(string)`配置信息分隔符设置，默认是`.`

- `WithWeaklyTypedInput(bool)` 调用`DecodeToStruct`时,会启用弱类型转化

- `WithLogger(Logger)` 外部日志接入，需实现[Logger](https://pkg.go.dev/github.com/enpsl/conf-reload@v1.0.0#Logger)，不传入会默认使用项目自带终端输出方式记录日志

- `WithCapacity(int)` `LRU`缓存容量设置，低于当前设置级别的日志记录不会在终端输出，可按照下面展示的级别进行设置

- `WithWatched(int)` 是否开启`Broker Watch`检测，某些场景如命令行模式，不需要热加载，可关闭此选项即可停止文件监听

- `WithLogLevel(int)`日志[级别](https://pkg.go.dev/github.com/enpsl/conf-reload@v1.0.0/internal/log#Level)设置，低于当前设置级别的日志记录不会在终端输出，可按照下面展示的级别进行设置

配置信息读取,可更改文件内容观察文件变化情况
```go
var http = &Http{}
for {
    err := conf_reload.DecodeToStruct("server.http", http)
    if err != nil {
        panic(err)
    }
    fmt.Println(http)
    time.Sleep(2 * time.Second)
}
```
输出示例:
```bash
conf-reload-test (main) ✗ go run _example/example.go
conf-reload@v1.0.0: pid=16749 2023/02/19 12:37:42.571884 DEBUG: map[server:map[config:map[connection:false depends:[tcp ip] publish:2023-02-19 timeout:10s] http:map[host:0.0.0.0 port:8080]]]
conf-reload@v1.0.0: pid=16749 2023/02/19 12:37:42.571917 DEBUG: map[config:map[connection:false depends:[tcp ip] publish:2023-02-19 timeout:10s] http:map[host:0.0.0.0 port:8080]]
&{0.0.0.0 8080}
&{0.0.0.0 8080}
&{0.0.0.0 8080}
conf-reload@v1.0.0: pid=16749 2023/02/19 12:37:46.693127 DEBUG: modified file:/Users/kuailexingqiu/go/src/conf-reload-test/_example/example.toml, /Users/kuailexingqiu/go/src/conf-reload-test/_example/example.toml
conf-reload@v1.0.0: pid=16749 2023/02/19 12:37:46.693437 DEBUG: map[server:map[config:map[connection:false depends:[tcp ip] publish:2023-02-19 timeout:10s] http:map[host:0.0.0.0 port:8081]]]
conf-reload@v1.0.0: pid=16749 2023/02/19 12:37:48.572885 DEBUG: map[config:map[connection:false depends:[tcp ip] publish:2023-02-19 timeout:10s] http:map[host:0.0.0.0 port:8081]]
&{0.0.0.0 8081}
&{0.0.0.0 8081}
```

如果想了解更多api，See [godoc](https://pkg.go.dev/github.com/enpsl/conf-reload@v1.0.0#pkg-functions)

# License
Copyright (c) 2023-present enpsl. conf-reload is free and open-source software licensed under the MIT License. 