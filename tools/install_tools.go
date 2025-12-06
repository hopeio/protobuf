/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/hopeio/gox"
	"github.com/hopeio/gox/log"
	"github.com/hopeio/gox/os/exec"
	"github.com/hopeio/gox/os/fs"
)

// 提供给使用框架的人安装所需环境
func main() {
	exec.Run("go version")
	libDir, _ := exec.RunGetOutWithLog("go list -m -f {{.Dir}}  github.com/hopeio/protobuf")
	os.Chdir(libDir)
	exec.RunGetOutWithLog("go install google.golang.org/protobuf/cmd/protoc-gen-go")
	protoccmd := "protoc -I" + libDir + "/_proto --go_out=paths=source_relative:" + libDir + " " + libDir + "/_proto/hopeio/utils/"
	//execx.RunGetOutWithLog(protoccmd + "patch/*.proto")
	exec.RunGetOutWithLog(gox.TernaryOperator(runtime.GOOS != "windows", fmt.Sprintf(`bash -c "%s"`, protoccmd+"apiconfig/*.proto"), protoccmd+"apiconfig/*.proto"))
	exec.RunGetOutWithLog(gox.TernaryOperator(runtime.GOOS != "windows", fmt.Sprintf(`bash -c "%s"`, protoccmd+"openapiconfig/*.proto"), protoccmd+"openapiconfig/*.proto"))
	exec.RunGetOutWithLog(gox.TernaryOperator(runtime.GOOS != "windows", fmt.Sprintf(`bash -c "%s"`, protoccmd+"enum/*.proto"), protoccmd+"enum/*.proto"))
	fs.MoveDirByMode(libDir+"/hopeio", libDir, 0)
	exec.RunGetOutWithLog("go install " + libDir + "/tools/protoc-gen-grpc-gin")
	exec.RunGetOutWithLog("go install " + libDir + "/tools/protoc-gen-enum")
	exec.RunGetOutWithLog("go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway")
	exec.RunGetOutWithLog("go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2")
	//execx.RunGetOutWithLog("go install github.com/alta/protopatch/cmd/protoc-gen-go-patch")
	exec.RunGetOutWithLog("go install google.golang.org/grpc/cmd/protoc-gen-go-grpc")
	//execx.RunGetOutWithLog("go install github.com/envoyproxy/protoc-gen-validate")
	exec.RunGetOutWithLog("go install " + libDir + "/tools/protoc-gen-validator")
	exec.RunGetOutWithLog("go install " + libDir + "/tools/protoc-gen-go-patch")
	exec.RunGetOutWithLog("go install " + libDir + "/tools/protoc-gen-gql")
	//execx.RunGetOutWithLog("go install github.com/danielvladco/go-proto-gql/protoc-gen-gogql")
	exec.RunGetOutWithLog("go install " + libDir + "/tools/protoc-gen-gogql")
	exec.RunGetOutWithLog("go install github.com/99designs/gqlgen")
	exec.RunGetOutWithLog("go install " + libDir + "/tools/protogen")
	log.Info("install success")
}
