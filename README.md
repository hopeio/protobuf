# protobuf

[![Go Reference](https://pkg.go.dev/badge/github.com/hopeio/protobuf.svg)](https://pkg.go.dev/github.com/hopeio/protobuf)
[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-jybl%2Fprotogen-blue?logo=docker)](https://hub.docker.com/r/jybl/protogen)

**One `.proto` → gRPC, HTTP gateway, OpenAPI, and validation.**

**一份 `.proto` → gRPC、HTTP 网关、OpenAPI 与参数校验。**

![protobuf](_assets/protobuf.webp)

```bash
go get github.com/hopeio/protobuf@latest
```

---

## English

### What is this repo?

**protobuf** is a Go module that packages:

1. **`protogen`** — a CLI that orchestrates `protoc` plugins in one command  
2. **Custom plugins** — `protoc-gen-gateway` (net/http · Gin · Fiber) and `protoc-gen-validator`  
3. **Shared `.proto` includes** under `_proto/` (Google APIs subset, well-known types, OpenAPI options, protopatch, common request/response/model messages)  
4. **Generated Go helpers** — `model/`, `request/`, `response/`, `utils/` for those common messages  

Write service contracts once; generate server stubs, HTTP routes, docs, and validators together.

### Install

Requires [protoc](https://github.com/protocolbuffers/protobuf/releases) on your `PATH`.

```bash
go get github.com/hopeio/protobuf@main
go run "$(go list -m -f '{{.Dir}}' github.com/hopeio/protobuf)/tools/install_tools.go"
```

This installs `protogen`, `protoc-gen-go`, `protoc-gen-go-grpc`, `protoc-gen-go-patch`, `protoc-gen-openapiv2`, plus this repo’s `protoc-gen-gateway` and `protoc-gen-validator`.

### Generate

```bash
protogen go -d -w -v -i ./proto -o ./gen
```

| Flag | Meaning |
|------|---------|
| `-i` / `--proto` | Input proto directory |
| `-o` / `--output` | Output directory |
| `-d` / `--openapi` | OpenAPI (openapiv2) |
| `-w` / `--gateway` | HTTP gateway code |
| `-v` / `--validator` | Validation code |
| `-I` / `--include` | Extra include paths |
| `-p` / `--protobuf` | Path to this module’s `_proto` |
| `--framework` | `nethttp` (default), `gin`, or `fiber` |

Subcommands: `protogen go` · `protogen dart` · `protogen ts`.

### Docker

No local plugin install required:

```bash
docker run --rm -v "$PWD:/work" jybl/protogen \
  protogen go -d -w -v -i /work/proto -o /work/gen
```

### Example contract

```protobuf
syntax = "proto3";
package user;

import "google/api/annotations.proto";
import "protoc-gen-openapiv2/options/annotations.proto";
import "hopeio/request/param.proto";
import "patch/go.proto";

option go_package = "example/user";

service UserService {
  rpc Info (hopeio.request.Id) returns (User) {
    option (google.api.http) = { get: "/api/user/{id}" };
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
      summary: "Get user by id"
    };
  }
}

message User {
  uint64 id = 1 [(go.field) = {tags: 'gorm:"primaryKey;"'}];
  string name = 2 [(go.field) = {tags: 'gorm:"size:64;not null"'}];
}
```

### Layout

```
_proto/                 # includes shipped with the module
model/ request/ response/ utils/   # Go packages for common messages
protogen/               # helpers used by plugins
tools/
  protogen/             # CLI
  protoc-gen-gateway/
  protoc-gen-validator/
  install_tools.go
  docker/
copyproto.sh            # refresh vendored protos from module cache (maintainers)
```

Regenerate this module’s own packages:

```bash
go run ./tools/generate.go
```

### License

[MIT](LICENSE)

---

## 中文

### 这个仓库是什么？

**protobuf** 是一个 Go module，打包了：

1. **`protogen`** — 一条命令编排多个 `protoc` 插件  
2. **自研插件** — `protoc-gen-gateway`（net/http · Gin · Fiber）与 `protoc-gen-validator`  
3. **公共 `_proto/`** — Google APIs 子集、well-known types、OpenAPI options、protopatch，以及通用 request/response/model  
4. **已生成的 Go 辅助包** — `model/`、`request/`、`response/`、`utils/`  

服务契约写一次，同时生成 gRPC 桩、HTTP 路由、文档与校验代码。

### 安装

需本机安装 [protoc](https://github.com/protocolbuffers/protobuf/releases)。

```bash
go get github.com/hopeio/protobuf@main
go run "$(go list -m -f '{{.Dir}}' github.com/hopeio/protobuf)/tools/install_tools.go"
```

将安装 `protogen`、`protoc-gen-go`、`protoc-gen-go-grpc`、`protoc-gen-go-patch`、`protoc-gen-openapiv2`，以及本仓的 `protoc-gen-gateway`、`protoc-gen-validator`。

### 生成

```bash
protogen go -d -w -v -i ./proto -o ./gen
```

| 标志 | 含义 |
|------|------|
| `-i` / `--proto` | proto 输入目录 |
| `-o` / `--output` | 输出目录 |
| `-d` / `--openapi` | OpenAPI（openapiv2） |
| `-w` / `--gateway` | HTTP 网关代码 |
| `-v` / `--validator` | 校验代码 |
| `-I` / `--include` | 额外 include |
| `-p` / `--protobuf` | 本模块 `_proto` 路径 |
| `--framework` | `nethttp`（默认）、`gin` 或 `fiber` |

子命令：`protogen go` · `protogen dart` · `protogen ts`。

### Docker

无需本机安装插件：

```bash
docker run --rm -v "$PWD:/work" jybl/protogen \
  protogen go -d -w -v -i /work/proto -o /work/gen
```

### 契约示例

```protobuf
syntax = "proto3";
package user;

import "google/api/annotations.proto";
import "protoc-gen-openapiv2/options/annotations.proto";
import "hopeio/request/param.proto";
import "patch/go.proto";

option go_package = "example/user";

service UserService {
  rpc Info (hopeio.request.Id) returns (User) {
    option (google.api.http) = { get: "/api/user/{id}" };
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
      summary: "按 id 获取用户"
    };
  }
}

message User {
  uint64 id = 1 [(go.field) = {tags: 'gorm:"primaryKey;"'}];
  string name = 2 [(go.field) = {tags: 'gorm:"size:64;not null"'}];
}
```

### 目录

```
_proto/                 # 随模块分发的 include
model/ request/ response/ utils/   # 公共消息的 Go 包
protogen/               # 插件辅助
tools/
  protogen/             # CLI
  protoc-gen-gateway/
  protoc-gen-validator/
  install_tools.go
  docker/
copyproto.sh            # 维护者：从 module cache 同步上游 proto
```

重新生成本模块自身包：

```bash
go run ./tools/generate.go
```

### 许可证

[MIT](LICENSE)
