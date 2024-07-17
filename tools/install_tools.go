package main

import (
	"github.com/hopeio/utils/io/fs"
	"github.com/hopeio/utils/log"
	osi "github.com/hopeio/utils/os"
	"os"
)

// 提供给使用框架的人安装所需环境
func main() {
	osi.StdOutCmd("go version")
	libDir, _ := osi.CmdLog("go list -m -f {{.Dir}}  github.com/hopeio/protobuf")
	os.Chdir(libDir)
	osi.CmdLog("go install google.golang.org/protobuf/cmd/protoc-gen-go")
	protoccmd := "protoc -I" + libDir + "/_proto --go_out=paths=source_relative:" + libDir + " " + libDir + "/_proto/hopeio/utils/"
	//osi.CmdLog(protoccmd + "patch/*.proto")
	osi.CmdLog(protoccmd + "apiconfig/*.proto")
	osi.CmdLog(protoccmd + "openapiconfig/*.proto")
	osi.CmdLog(protoccmd + "enum/*.proto")
	fs.MoveDirByMode(libDir+"/hopeio", libDir, 0)
	osi.CmdLog("go install " + libDir + "/tools/protoc-gen-grpc-gin")
	osi.CmdLog("go install " + libDir + "/tools/protoc-gen-enum")
	osi.CmdLog("go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway")
	osi.CmdLog("go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2")
	//osi.CmdLog("go install github.com/alta/protopatch/cmd/protoc-gen-go-patch")
	osi.CmdLog("go install google.golang.org/grpc/cmd/protoc-gen-go-grpc")
	//osi.CmdLog("go install github.com/envoyproxy/protoc-gen-validate")
	osi.CmdLog("go install " + libDir + "/tools/protoc-gen-validator")
	osi.CmdLog("go install " + libDir + "/tools/protoc-gen-go-patch")
	osi.CmdLog("go install " + libDir + "/tools/protoc-gen-gql")
	//osi.CmdLog("go install github.com/danielvladco/go-proto-gql/protoc-gen-gogql")
	osi.CmdLog("go install " + libDir + "/tools/protoc-gen-gogql")
	osi.CmdLog("go install github.com/99designs/gqlgen")
	osi.CmdLog("go install " + libDir + "/tools/protogen")
	log.Info("安装成功")
}
