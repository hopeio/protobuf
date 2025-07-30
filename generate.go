/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package main

import (
	"fmt"
	execi "github.com/hopeio/gox/os/exec"
	"github.com/hopeio/gox/os/fs"
	_go "github.com/hopeio/gox/sdk/go"
	"log"
	"os"
	"strings"
)

//go:generate mockgen -destination ../protobuf/user/user.mock.go -package user -source ../protobuf/user/user.service_grpc.pb.go UserServiceServer

var (
	libProtobufDir, proto string
	pwd, gopath, include  string
)

func init() {
	gopath = os.Getenv("GOPATH")
	if strings.HasSuffix(gopath, "/") {
		gopath = gopath[:len(gopath)-1]
	}

	pwd, _ = os.Getwd()
	libProtobufDir = _go.GetDepDir(DepProtobuf)
	proto = libProtobufDir + "/_proto"
	//libGatewayDir := _go.GetDepDir(DepGrpcGateway)
	//libGoogleDir := _go.GetDepDir(DepGoogleapis)

	include = "-I" + proto
}

func main() {
	//single("/content/moment.model.proto")
	generate(proto + "/hopeio")
	fmt.Println(fs.MoveDirByMode(libProtobufDir+"/hopeio", libProtobufDir, 0))

}

const goOut = "go-patch_out=plugin=go,paths=source_relative"
const grpcOut = "go-patch_out=plugin=go-grpc,paths=source_relative"
const enumOut = "enum_out=paths=source_relative"

const (
	goListDir      = `go list -m -f {{.Dir}} `
	goListDep      = `go list -m -f {{.Path}}@{{.Version}} `
	DepGoogleapis  = "github.com/googleapis/googleapis@v0.0.0-20220520010701-4c6f5836a32f"
	DepProtobuf    = "github.com/hopeio/protobuf"
	DepGrpcGateway = "github.com/grpc-ecosystem/grpc-gateway/v2"
	DepProtopatch  = "github.com/alta/protopatch"
)

var model = []string{goOut, grpcOut, enumOut}

func generate(dir string) {
	fileInfos, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalln(err)
	}
	var gen bool
	for i := range fileInfos {
		if fileInfos[i].IsDir() {
			generate(dir + "/" + fileInfos[i].Name())
		} else if !gen && strings.HasSuffix(fileInfos[i].Name(), ".proto") {
			protoc(dir)
		}
	}
}

func protoc(dir string) {
	cmd := "protoc " + include + " " + dir + "/*.proto"
	var args string
	for _, plugin := range model {
		args += " --" + plugin + ":" + libProtobufDir
	}
	execi.RunWithLog(cmd + args)
}
