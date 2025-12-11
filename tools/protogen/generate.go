/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	execx "github.com/hopeio/gox/os/exec"

	gox "github.com/hopeio/gox/sdk/go"
	"github.com/spf13/cobra"
)

/*
*文件名正则不支持以及enum生成和model生成用的都是gogo的，所以顺序必须是gogo_out在前，enum_out在后
 */

//go:generate mockgen -destination ../protobuf/user/user.mock.go -package user -source ../protobuf/user/user.service_grpc.pb.go UserServiceServer

const (
	openapiv2Out = "openapiv2_out=allow_merge=true,merge_file_name="
	gqlOut       = "gql_out=svc=true,merge=true,paths=source_relative"
)

const (
	goListDir     = `go list -m -f {{.Dir}} `
	goListDep     = `go list -m -f {{.Path}}@{{.Version}} `
	DepGoogleapis = "github.com/googleapis/googleapis@v0.0.0-20220520010701-4c6f5836a32f"
	DepProtobuf   = "github.com/hopeio/protobuf"
)

var plugins = []string{}
var pluginPaths = []string{}

//"gqlgencfg_out=paths=source_relative",

var rootCmd = &cobra.Command{
	Use: "protogen",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		os.MkdirAll(config.genpath, os.ModePerm)
		if config.useGqlPlugin {
			plugins = append(plugins, gqlOut)
		}
		if config.useOpenapiPlugin {
			plugins = append(plugins, openapiv2Out)
		}
		if config.useGqlPlugin || config.useOpenapiPlugin {
			if config.apidocDir == "" {
				var err error
				config.apidocDir, err = filepath.Abs(config.genpath + "/../apidoc")
				if err != nil {
					config.apidocDir = config.genpath + "/apidoc"
				}
			} else {
				if config.apidocDir[len(config.apidocDir)-1] != '/' {
					config.apidocDir += "/"
				}
			}
		}
		getInclude()
	},
}

func init() {
	protodef, _ := filepath.Abs("/proto")
	pwd, _ := os.Getwd()
	pflag := rootCmd.PersistentFlags()
	pflag.StringVarP(&config.proto, "proto", "i", protodef, "proto dir")
	pflag.StringVarP(&config.genpath, "output", "o", pwd+"/protobuf", "generate dir")
	pflag.StringVarP(&config.currentDir, "protobuf", "p", "/proto", "手动指定protobuf项目目录")
	pflag.StringArrayVarP(&config.thirdIncludes, "include", "I", nil, "外部proto依赖")
	pflag.BoolVarP(&config.useGqlPlugin, "graphql", "g", false, "是否使用graphql插件")
	pflag.BoolVarP(&config.useOpenapiPlugin, "openapi", "d", false, "是否使用openapi插件")
	pflag.StringVar(&config.apidocDir, "apiDocDir", "", "api doc目录")

	rootCmd.AddCommand(goCmd)
	rootCmd.AddCommand(dartCmd)
	rootCmd.AddCommand(tsCmd)

}

func main() {
	//single("/content/moment.model.proto")
	rootCmd.Execute()
	//gengql()

}

func run(dir string) {
	fileInfos, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalln(err)
	}
	var exec bool
	for i := range fileInfos {
		if !exec && strings.HasSuffix(fileInfos[i].Name(), ".proto") {
			exec = true
			if dir == config.proto {
				protocCmd(dir+"/*.proto", "root", "")
			} else {
				protocCmd(dir+"/*.proto", filepath.Base(dir), dir[len(config.proto)+1:])
			}
		}
		if fileInfos[i].IsDir() {
			run(dir + "/" + fileInfos[i].Name())
		}
	}
}
func getInclude() {
	pwd, _ := os.Getwd()
	config.proto, _ = filepath.Abs(config.proto)
	config.genpath, _ = filepath.Abs(config.genpath)
	log.Println("proto:", config.proto)
	log.Println("genpath:", config.genpath)

	/*	generatePath := "generate" + time.Now().Format("150405")
		err = os.Mkdir(generatePath, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		generatePath = pwd + "/" + generatePath
		defer os.RemoveAll(generatePath)
		err = os.Chdir(generatePath)
		if err != nil {
			log.Fatal(err)
		}
		execx.RunGetOut("go mod init generate")*/

	libDir, err := execx.RunGetOut(gox.GoListDir + DepProtobuf)
	if err == nil {
		config.currentDir = libDir + "/_proto"
	}
	config.include = "-I" + config.currentDir + " -I" + config.proto
	for _, include := range config.thirdIncludes {
		config.include += " -I" + include
	}
	/*	os.Chdir(libDir)
		DepGrpcGateway, _ = execx.RunGetOut(goListDep + DepGrpcGateway)
		DepProtopatch, _ = execx.RunGetOut(goListDep + DepProtopatch)
		os.Chdir(generatePath)
		libGoogleDir := gox.GetDepDir(DepGoogleapis)
		libGatewayDir := gox.GetDepDir(DepGrpcGateway)*/

	os.Chdir(pwd)

	log.Println("include:", config.include)

}

// 找出所以包含go文件的包

func getPackages(dir string) []string {
	p := make(map[string]struct{})
	getPackagesHelper(dir, "", p)
	var packages []string
	for pkg, _ := range p {
		packages = append(packages, pkg)
	}
	return packages
}

func getPackagesHelper(dir, pre string, p map[string]struct{}) {
	fileInfos, err := os.ReadDir(dir)
	if err != nil {
		log.Panicln(err)
	}
	for i := range fileInfos {
		if fileInfos[i].IsDir() {
			files, err := os.ReadDir(dir + "/" + fileInfos[i].Name())
			if err != nil {
				log.Panicln(err)
			}
			for j := range files {
				if strings.HasSuffix(files[j].Name(), ".go") {
					if pre != "" {
						p[pre+"/"+fileInfos[i].Name()] = struct{}{}
					} else {
						p[fileInfos[i].Name()] = struct{}{}
					}

					break
				}
			}
			getPackagesHelper(dir+"/"+fileInfos[i].Name(), fileInfos[i].Name(), p)
		}
	}
}
func protocCmd(file, mod, modDir string) {
	cmd := "protoc " + config.include + " " + file
	var args string

	for _, pluginPath := range pluginPaths {
		args += " --plugin=" + pluginPath
	}

	for _, plugin := range plugins {
		genpath := config.genpath
		if strings.HasPrefix(plugin, "openapiv2_out") {
			plugin += mod
			genpath = config.apidocDir
		}

		if strings.HasPrefix(plugin, "gql_out") {
			genpath = config.apidocDir
		}
		args += " --" + plugin + ":" + genpath

	}
	cmd += args
	if runtime.GOOS != "windows" {
		cmd = fmt.Sprintf(`bash -c "%s"`, cmd)
	}
	execx.RunWithLog(cmd)
}
