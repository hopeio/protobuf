# protobuf

[![Go Reference](https://pkg.go.dev/badge/github.com/hopeio/protobuf.svg)](https://pkg.go.dev/github.com/hopeio/protobuf)
[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-jybl%2Fprotogen-blue?logo=docker)](https://hub.docker.com/r/jybl/protogen)

[English](README.md)

```bash
go get github.com/hopeio/protobuf@latest
```

面向 Go 的 Protocol Buffers 工具集：`protogen` CLI、gateway / validator 插件、公共 `_proto`，以及通用消息的生成辅助包。

服务契约写一次，一次生成 gRPC 桩、HTTP 路由、OpenAPI 与校验代码。

## 模块里有什么

| 部分 | 位置 | 作用 |
|------|------|------|
| CLI | `tools/protogen` | 编排 `protoc` 插件（`go` / `dart` / `ts`） |
| Gateway 插件 | `tools/protoc-gen-gateway` | 生成 `net/http`、Gin 或 Fiber 的 HTTP 处理 |
| Validator 插件 | `tools/protoc-gen-validator` | 从 proto 选项生成校验 |
| Include | `_proto/` | googleapis 子集、well-known types、OpenAPI options、protopatch、hopeio request/response/model |
| Go 辅助包 | `model/`、`request/`、`response/`、`utils/` | 上述公共消息的编译结果 |
| 插件辅助 | `protogen/` | 生成器用的 option 读取 |

## 安装插件

需要本机 [`protoc`](https://github.com/protocolbuffers/protobuf/releases)。

```bash
go get github.com/hopeio/protobuf@main
go run "$(go list -m -f '{{.Dir}}' github.com/hopeio/protobuf)/tools/install_tools.go"
```

将安装：`protogen`、`protoc-gen-go`、`protoc-gen-go-grpc`、`protoc-gen-go-patch`、`protoc-gen-openapiv2`、`protoc-gen-gateway`、`protoc-gen-validator`。

## 生成

```bash
protogen go -d -w -v -i ./api -o ./gen
```

| 标志 | 作用 |
|------|------|
| `-i` | proto 输入目录 |
| `-o` | 输出目录 |
| `-d` | OpenAPI（openapiv2） |
| `-w` | HTTP gateway |
| `-v` | 校验代码 |
| `-I` | 额外 include |
| `-p` | 本模块 `_proto` 路径 |
| `--framework` | `nethttp`（默认）、`gin`、`fiber` |

另有：`protogen dart`、`protogen ts`。

### Docker

```bash
docker run --rm -v "$PWD:/work" jybl/protogen \
  protogen go -d -w -v -i /work/api -o /work/gen
```

## 契约示例

```protobuf
syntax = "proto3";
package catalog;

import "google/api/annotations.proto";
import "protoc-gen-openapiv2/options/annotations.proto";
import "hopeio/request/param.proto";
import "patch/go.proto";

option go_package = "example.com/catalog";

service ProductService {
  rpc Get (hopeio.request.Id) returns (Product) {
    option (google.api.http) = { get: "/api/products/{id}" };
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
      summary: "获取商品"
    };
  }
}

message Product {
  uint64 id = 1 [(go.field) = {tags: 'gorm:"primaryKey;"'}];
  string sku = 2 [(go.field) = {tags: 'gorm:"size:64;uniqueIndex"'}];
}
```

## 维护命令

```bash
# 重新生成本模块 hopeio 相关包
go run ./tools/generate.go

# 把上游 proto 同步进 _proto/
./copyproto.sh
```

## 许可证

[MIT](LICENSE)
