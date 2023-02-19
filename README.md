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

# Quick Start

配置引擎加载,加载一次即可
```go
f = "_example/example.toml"
conf_reload.LoadEngine(f, conf_relod.WithLevelSplit("."), conf_relod.WithLogLevel(0))
```
`conf_relod.WithLevelSplit(".")`配置信息分隔符设置，默认是`.`
`conf_relod.WithLogLevel(0)`日志级别设置，低于当前设置级别的日志记录不会在终端输出，可按照下面展示的级别进行设置

```go
DebugLevel = 0

InfoLevel = 1

WarnLevel = 2

ErrorLevel = 3

FatalLevel = 4
```

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

```bash
conf-reload@v0.0.0: pid=30342 2023/02/19 09:04:04.466710 DEBUG: map[server:map[config:map[connection:false depends:[tcp ip] publish:2023-02-19 timeout:10s] http:map[host:0.0.0.0 port:8080]]]
conf-reload@v0.0.0: pid=30342 2023/02/19 09:04:04.466747 DEBUG: map[config:map[connection:false depends:[tcp ip] publish:2023-02-19 timeout:10s] http:map[host:0.0.0.0 port:8080]]
&{0.0.0.0 8080}
&{0.0.0.0 8080}
&{0.0.0.0 8080}
&{0.0.0.0 8080}
conf-reload@v0.0.0: pid=30342 2023/02/19 09:04:10.586570 DEBUG: modified file:/Users/kuailexingqiu/go/src/conf-reload/_example/example.toml, /Users/kuailexingqiu/go/src/conf-reload/_example/example.toml
conf-reload@v0.0.0: pid=30342 2023/02/19 09:04:10.586780 DEBUG: map[server:map[config:map[connection:false depends:[tcp ip] publish:2023-02-19 timeout:10s] http:map[host:0.0.0.0 port:80801]]]
conf-reload@v0.0.0: pid=30342 2023/02/19 09:04:12.470547 DEBUG: map[config:map[connection:false depends:[tcp ip] publish:2023-02-19 timeout:10s] http:map[host:0.0.0.0 port:80801]]
&{0.0.0.0 80801}
&{0.0.0.0 80801}
```

# License
Copyright (c) 2023-present enpsl. conf-reload is free and open-source software licensed under the MIT License. 