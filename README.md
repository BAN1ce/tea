# Tea


* 一个简单的MQTT Broker，目前只支持Qos=0等级。


## 项目目录说明

``` 

├── README.md
├── example
│   ├── main.go
│   └── tree.go
├── go.mod
├── go.sum
└── src
    ├── manage
    │   ├── client.go 客户端
    │   └── manage.go 客户端管理
    ├── mqtt
    │   ├── handle 处理客户端发送的数据
    │   │   ├── connect.go
    │   │   ├── disconnect.go
    │   │   ├── handle.go
    │   │   ├── hb.go
    │   │   ├── publish.go
    │   │   ├── subscribe.go
    │   │   └── unsubscribe.go
    │   ├── protocol mqtt协议的解析
    │   │   └── protocol.go
    │   ├── response 返回给客户端的包封装
    │   │   ├── connack.go
    │   │   ├── hback.go
    │   │   ├── publish.go
    │   │   ├── response.go
    │   │   ├── suback.go
    │   │   └── unsuback.go
    │   ├── route.go mqtt命令路由
    │   └── sub
    │       ├── hash.go 绝对订阅
    │       ├── tree.go 订阅树订阅
    │       └── tree_test.go
    ├── server
    │   └── server.go
    ├── unpack
    │   └── protocol.go
    └── utils 工具包
        ├── check.go
        └── convert.go

```


## TODO

|功能|状态|
|:---:|:---:|
|设备离线遗嘱消息的发布|未开始|
|自定义设备连接后的认证|未开始|

## BUG

|描述|状态|
|:---:|:---:|
|客户端名相同的连接应关闭前一个连接|未开始|




# License

MIT
