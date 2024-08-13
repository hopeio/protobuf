package main

import (
	"github.com/hopeio/utils/fs"
	"github.com/hopeio/utils/log"
	execi "github.com/hopeio/utils/os/exec"
	"os"
)

// 提供给使用框架的人安装所需环境
func main() {
	execi.StdOutCmd("go version")
	libDir, _ := execi.CmdLog("go list -m -f {{.Dir}}  github.com/hopeio/protobuf")
	os.Chdir(libDir)
	execi.CmdLog("go install google.golang.org/protobuf/cmd/protoc-gen-go")
	protoccmd := "protoc -I" + libDir + "/_proto --go_out=paths=source_relative:" + libDir + " " + libDir + "/_proto/hopeio/utils/"
	//execi.CmdLog(protoccmd + "patch/*.proto")
	execi.CmdLog(protoccmd + "apiconfig/*.proto")
	execi.CmdLog(protoccmd + "openapiconfig/*.proto")
	execi.CmdLog(protoccmd + "enum/*.proto")
	fs.MoveDirByMode(libDir+"/hopeio", libDir, 0)
	execi.CmdLog("go install " + libDir + "/tools/protoc-gen-grpc-gin")
	execi.CmdLog("go install " + libDir + "/tools/protoc-gen-enum")
	execi.CmdLog("go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway")
	execi.CmdLog("go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2")
	//execi.CmdLog("go install github.com/alta/protopatch/cmd/protoc-gen-go-patch")
	execi.CmdLog("go install google.golang.org/grpc/cmd/protoc-gen-go-grpc")
	//execi.CmdLog("go install github.com/envoyproxy/protoc-gen-validate")
	execi.CmdLog("go install " + libDir + "/tools/protoc-gen-validator")
	execi.CmdLog("go install " + libDir + "/tools/protoc-gen-go-patch")
	execi.CmdLog("go install " + libDir + "/tools/protoc-gen-gql")
	//execi.CmdLog("go install github.com/danielvladco/go-proto-gql/protoc-gen-gogql")
	execi.CmdLog("go install " + libDir + "/tools/protoc-gen-gogql")
	execi.CmdLog("go install github.com/99designs/gqlgen")
	execi.CmdLog("go install " + libDir + "/tools/protogen")
	log.Info("安装成功")
}
