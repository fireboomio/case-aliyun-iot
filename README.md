## 项目简介

本项目是基于飞布(Fireboom) Golang SDK 开发的阿里云 IoT 平台演示项目，通过自动化代码生成和丰富的示例，大幅提升阿里云 IoT
应用的开发效率。

## 主要特性

- 🔧 **自动代码生成**: 根据阿里云 IoT 的 `tsl.json` 配置文件自动生成业务代码
- ⚡ **高效开发**: 减少重复代码编写，专注于业务逻辑实现
- 📚 **丰富示例**: 包含多种常见 IoT 场景的实现示例
- 🖥️ **后台管理**: 集成基础设备管理功能
- 🔒 **安全可靠**: 基于阿里云 IoT 安全认证机制
- 🚀 **高性能**: 利用 Golang 的高并发特性处理海量设备连接

## 快速开始

### 前置条件

- Go 1.21 或更高版本
- 阿里云 IOT 平台账号
- 获取阿里云 IOT 平台的 `tsl.json` 配置文件
- 将您的 `tsl.json` 文件放置在 `custom-go/iot/{iot类型}` 目录下
- 复制.env.example` 为 `.env`，并填写正确的配置信息

## 项目结构

```
case-aliyun-iot/custom-go/
├── authentication/     # 认证及权限
├── crontab/            # 定时任务
├── customize/          # 自定义模块       
│   └── deviceTest      # SSE推送IOT事件
├── dict/               # 字典管理
├── function/           # 暴露iot接口
├── proxy/              # 文件导入示例
├── global/             # 全局钩子，记录请求日志
├── pkg/                # 可公开使用的包
├── iot/                # 物联网模块
│   └── base/           # IOT 基础模块
│   └── config/         # tls 代码示例
│   └── defaulted/      # tls 代码示例
│   └── hbs/            # tls 模板生成
└── README.md           # 项目说明
```

## 代码生成器

本项目核心特性是根据阿里云 IoT 的 `tsl.json` 自动生成业务代码：

生成器将根据配置文件自动创建：

- 设备连接管理代码
- MQTT 主题订阅处理
- RRPC 业务调用
- 物模型数据处理

## 配置说明

### 环境变量配置

```bash
# 阿里云IOT配置
ALIBABA_CLOUD_ACCESS_KEY=""
ALIBABA_CLOUD_ACCESS_SECRET=""
ALIBABA_CLOUD_CLIENT_ID=""
ALIBABA_CLOUD_CONSUMER_GROUP_ID=""
ALIBABA_CLOUD_REGION_ID=""
ALIBABA_CLOUD_HOST=""
ALIBABA_CLOUD_IOT_INSTANCE_ID=""
ALIBABA_CLOUD_PRODUCT_KEY=""
```

### tsl.json 配置示例

[tsl.json](custom-go%2Fiot%2Fdefaulted%2Ftsl.json)

## 开发指南

### 生成自定义IOT代码

1. 在 `custom-go/iot` 目录下按iot类型创建目录并防止tsl.json
2. 在 `.env` 开启ENABLE_IOT_CODE_GENERATE=1

### 后台功能开发

1. 集成登录认证和权限管理
2. 集成了日志记录
3. 集成了字典管理（`custom-go/dict` 订阅并设计业务逻辑）

## 常见问题

### Q: 代码生成器不工作怎么办？

A: 请检查 `tsl.json` 文件格式是否正确，确保所有必填字段已填写

### Q: 设备无法连接怎么办？

A: 检查设备三元组（ProductKey、DeviceName、DeviceSecret）是否正确

### Q: 如何添加自定义主题订阅？

A: 修改代码生成器中的模板文件，添加新的主题订阅逻辑

**注意**: 本项目为演示用途，生产环境使用前请进行充分测试和安全评估。