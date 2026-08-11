# protobuf

[![Go Reference](https://pkg.go.dev/badge/github.com/hopeio/protobuf.svg)](https://pkg.go.dev/github.com/hopeio/protobuf)
[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-jybl%2Fprotogen-blue?logo=docker)](https://hub.docker.com/r/jybl/protogen)

[中文文档](README.zh-CN.md)

```bash
go get github.com/hopeio/protobuf@latest
```

Toolkit for Protocol Buffers in Go: a `protogen` CLI, gateway & validator plugins, shared `_proto` includes, and generated helper packages for common messages.

Describe services once. Generate gRPC stubs, HTTP routes, OpenAPI, and validation in one pass.

## Contents of this module

| Piece | Location | Role |
|-------|----------|------|
| CLI | `tools/protogen` | Orchestrates `protoc` plugins (`go` / `dart` / `ts`) |
| Gateway plugin | `tools/protoc-gen-gateway` | HTTP handlers for `net/http`, Gin, or Fiber |
| Validator plugin | `tools/protoc-gen-validator` | Request validation from proto options |
| Includes | `_proto/` | googleapis subset, well-known types, OpenAPI options, protopatch, hopeio request/response/model |
| Go helpers | `model/`, `request/`, `response/`, `utils/` | Compiled packages for those shared messages |
| Plugin helpers | `protogen/` | Option readers used by generators |

## Install plugins

Requires [`protoc`](https://github.com/protocolbuffers/protobuf/releases) on `PATH`.

```bash
go get github.com/hopeio/protobuf@main
go run "$(go list -m -f '{{.Dir}}' github.com/hopeio/protobuf)/tools/install_tools.go"
```

Installs: `protogen`, `protoc-gen-go`, `protoc-gen-go-grpc`, `protoc-gen-go-patch`, `protoc-gen-openapiv2`, `protoc-gen-gateway`, `protoc-gen-validator`.

## Generate

```bash
protogen go -d -w -v -i ./api -o ./gen
```

| Flag | Purpose |
|------|---------|
| `-i` | Proto input directory |
| `-o` | Output directory |
| `-d` | OpenAPI (openapiv2) |
| `-w` | HTTP gateway |
| `-v` | Validators |
| `-I` | Extra includes |
| `-p` | Path to this module’s `_proto` |
| `--framework` | `nethttp` (default), `gin`, `fiber` |

Also: `protogen dart`, `protogen ts`.

### Docker

```bash
docker run --rm -v "$PWD:/work" jybl/protogen \
  protogen go -d -w -v -i /work/api -o /work/gen
```

## Contract sample

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
      summary: "Fetch a product"
    };
  }
}

message Product {
  uint64 id = 1 [(go.field) = {tags: 'gorm:"primaryKey;"'}];
  string sku = 2 [(go.field) = {tags: 'gorm:"size:64;uniqueIndex"'}];
}
```

## Maintainer commands

```bash
# regenerate this module’s hopeio packages
go run ./tools/generate.go

# refresh vendored third-party protos into _proto/
./copyproto.sh
```

## License

[MIT](LICENSE)
