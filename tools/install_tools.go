/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package main

import (
	"os"

	"github.com/hopeio/gox/log"
	execx "github.com/hopeio/gox/os/exec"
	"github.com/hopeio/gox/os/fs"
)

// 提供给使用框架的人安装所需环境
func main() {
	execx.Run("go version")
	libDir, _ := execx.RunGetOutWithLog("go list -m -f {{.Dir}}  github.com/hopeio/protobuf")
	os.Chdir(libDir)
	execx.RunGetOutWithLog("go install google.golang.org/protobuf/cmd/protoc-gen-go")
	protoccmd := "protoc -I" + libDir + "/_proto --go_out=paths=source_relative:" + libDir + " " + libDir + "/_proto/hopeio/utils/"
	//execx.RunGetOutWithLog(protoccmd + "patch/*.proto")
	execx.RunGetOutWithLog(protoccmd + "apiconfig/*.proto")
	execx.RunGetOutWithLog(protoccmd + "openapiconfig/*.proto")
	execx.RunGetOutWithLog(protoccmd + "enum/*.proto")
	fs.MoveDirByMode(libDir+"/hopeio", libDir, 0)
	execx.RunGetOutWithLog("go install " + libDir + "/tools/protoc-gen-grpc-gin")
	execx.RunGetOutWithLog("go install " + libDir + "/tools/protoc-gen-enum")
	execx.RunGetOutWithLog("go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway")
	execx.RunGetOutWithLog("go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2")
	//execx.RunGetOutWithLog("go install github.com/alta/protopatch/cmd/protoc-gen-go-patch")
	execx.RunGetOutWithLog("go install google.golang.org/grpc/cmd/protoc-gen-go-grpc")
	//execx.RunGetOutWithLog("go install github.com/envoyproxy/protoc-gen-validate")
	execx.RunGetOutWithLog("go install " + libDir + "/tools/protoc-gen-validator")
	execx.RunGetOutWithLog("go install " + libDir + "/tools/protoc-gen-go-patch")
	execx.RunGetOutWithLog("go install " + libDir + "/tools/protoc-gen-gql")
	//execx.RunGetOutWithLog("go install github.com/danielvladco/go-proto-gql/protoc-gen-gogql")
	execx.RunGetOutWithLog("go install " + libDir + "/tools/protoc-gen-gogql")
	execx.RunGetOutWithLog("go install github.com/99designs/gqlgen")
	execx.RunGetOutWithLog("go install " + libDir + "/tools/protogen")
	log.Info("安装成功")
}
