# toktik
#### 目录结构

```shell
.
├── README.md
├── docs                           // 文档
├── go.mod                         // Go Modules
├── idl                            // thrift/protobuf文件
├── internal                       // 私有库/服务实现
│   └── gateway
│       ├── biz
│       │   ├── handler            // 逻辑
│       │   │   └── ping.go
│       │   └── router             // 路由
│       │       └── register.go
│       ├── etc                    // 配置文件
│       │   └── config.yaml    
│       ├── main.go
│       ├── router.go
│       └── router_gen.go
├── pkg                            // 公共包
└── sh                             // Shell脚本
```

