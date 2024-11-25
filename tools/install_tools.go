/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package main

import (
	"github.com/hopeio/utils/fs"
	"github.com/hopeio/utils/log"
	execi "github.com/hopeio/utils/os/exec"
	"os"
)

// 提供给使用框架的人安装所需环境
func main() {
	execi.Run("go version")
	libDir, _ := execi.RunGetOutWithLog("go list -m -f {{.Dir}}  github.com/hopeio/protobuf")
	os.Chdir(libDir)
	execi.RunGetOutWithLog("go install google.golang.org/protobuf/cmd/protoc-gen-go")
	protoccmd := "protoc -I" + libDir + "/_proto --go_out=paths=source_relative:" + libDir + " " + libDir + "/_proto/hopeio/utils/"
	//execi.RunGetOutWithLog(protoccmd + "patch/*.proto")
	execi.RunGetOutWithLog(protoccmd + "apiconfig/*.proto")
	execi.RunGetOutWithLog(protoccmd + "openapiconfig/*.proto")
	execi.RunGetOutWithLog(protoccmd + "enum/*.proto")
	fs.MoveDirByMode(libDir+"/hopeio", libDir, 0)
	execi.RunGetOutWithLog("go install " + libDir + "/tools/protoc-gen-grpc-gin")
	execi.RunGetOutWithLog("go install " + libDir + "/tools/protoc-gen-enum")
	execi.RunGetOutWithLog("go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway")
	execi.RunGetOutWithLog("go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2")
	//execi.RunGetOutWithLog("go install github.com/alta/protopatch/cmd/protoc-gen-go-patch")
	execi.RunGetOutWithLog("go install google.golang.org/grpc/cmd/protoc-gen-go-grpc")
	//execi.RunGetOutWithLog("go install github.com/envoyproxy/protoc-gen-validate")
	execi.RunGetOutWithLog("go install " + libDir + "/tools/protoc-gen-validator")
	execi.RunGetOutWithLog("go install " + libDir + "/tools/protoc-gen-go-patch")
	execi.RunGetOutWithLog("go install " + libDir + "/tools/protoc-gen-gql")
	//execi.RunGetOutWithLog("go install github.com/danielvladco/go-proto-gql/protoc-gen-gogql")
	execi.RunGetOutWithLog("go install " + libDir + "/tools/protoc-gen-gogql")
	execi.RunGetOutWithLog("go install github.com/99designs/gqlgen")
	execi.RunGetOutWithLog("go install " + libDir + "/tools/protogen")
	log.Info("安装成功")
}
