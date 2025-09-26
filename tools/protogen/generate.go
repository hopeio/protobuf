/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	execx "github.com/hopeio/gox/os/exec"
	"github.com/hopeio/gox/os/fs"
	gox "github.com/hopeio/gox/sdk/go"
	"github.com/spf13/cobra"
)

/*
*文件名正则不支持以及enum生成和model生成用的都是gogo的，所以顺序必须是gogo_out在前，enum_out在后
 */

//go:generate mockgen -destination ../protobuf/user/user.mock.go -package user -source ../protobuf/user/user.service_grpc.pb.go UserServiceServer

const (
	goOut         = "go-patch_out=plugin=go,paths=source_relative"
	grpcOut       = "go-patch_out=plugin=go-grpc,paths=source_relative"
	enumOut       = "enum_out=paths=source_relative"
	gatewayOut    = "grpc-gin_out=paths=source_relative"
	openapiv2Out  = "openapiv2_out=allow_merge=true,merge_file_name="
	validatorsOut = "validator_out=paths=source_relative"
	gqlOut        = "gql_out=svc=true,merge=true,paths=source_relative"
	gogqlOut      = "gogql_out=svc=true,merge=true,paths=source_relative"
	dartOut       = "dart_out=grpc"
)

const (
	goListDir      = `go list -m -f {{.Dir}} `
	goListDep      = `go list -m -f {{.Path}}@{{.Version}} `
	DepGoogleapis  = "github.com/googleapis/googleapis@v0.0.0-20220520010701-4c6f5836a32f"
	DepProtobuf    = "github.com/hopeio/protobuf"
	DepGrpcGateway = "github.com/grpc-ecosystem/grpc-gateway/v2"
	DepProtopatch  = "github.com/alta/protopatch"
)

var plugin = []string{goOut, grpcOut}

//"gqlgencfg_out=paths=source_relative",

var enumPlugin = enumOut
var gatewayPlugin = []string{gatewayOut, openapiv2Out}
var validatorOutPlugin = validatorsOut
var gqlPlugin = []string{gqlOut, gogqlOut}

var rootCmd = &cobra.Command{
	Use: "protogen",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		os.MkdirAll(config.genpath, os.ModePerm)
		if config.useEnumPlugin {
			plugin = append(plugin, enumPlugin)
		}
		if config.useGateWayPlugin {
			plugin = append(plugin, gatewayPlugin...)
		}
		if config.useValidatorOutPlugin {
			plugin = append(plugin, validatorOutPlugin)
		}
		if config.useGqlPlugin {
			plugin = append(plugin, gqlPlugin...)
		}
		getInclude()
	},
}

func init() {
	protodef, _ := filepath.Abs("/proto")
	pwd, _ := os.Getwd()
	pflag := rootCmd.PersistentFlags()
	pflag.StringVarP(&config.proto, "proto", "p", protodef, "proto dir")
	pflag.StringVarP(&config.genpath, "genpath", "o", pwd+"/protobuf", "generate dir")
	pflag.StringVarP(&config.currentDir, "current", "c", "/proto", "手动指定protobuf项目目录")
	pflag.StringArrayVarP(&config.thirdIncludes, "include", "I", nil, "外部proto依赖")
	pflag.BoolVarP(&config.useEnumPlugin, "enum", "e", false, "是否使用enum扩展插件")
	pflag.BoolVarP(&config.useGateWayPlugin, "gateway", "w", false, "是否使用grpc-gateway插件")
	pflag.BoolVarP(&config.useValidatorOutPlugin, "validator", "v", false, "是否使用validators插件")
	pflag.BoolVarP(&config.useGqlPlugin, "graphql", "g", false, "是否使用graphql插件")
	pflag.BoolVar(&config.stdPatch, "patch", false, "是否使用原生protopatch")
	pflag.StringVar(&config.apidocDir, "apiDocDir", "", "api doc目录")
	rootCmd.AddCommand(&cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use: "go",
		Run: func(cmd *cobra.Command, args []string) {
			run(config.proto)
			if config.useGqlPlugin {
				gengql()
			}
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use: "dart",
		Run: func(cmd *cobra.Command, args []string) {

		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use: "ts",
		Run: func(cmd *cobra.Command, args []string) {

		},
	})

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
				protocCmd(plugin, dir+"/*.proto", "root", "")
			} else {
				protocCmd(plugin, dir+"/*.proto", filepath.Base(dir), dir[len(config.proto)+1:])
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
	if config.apidocDir == "" {
		config.apidocDir = config.genpath + "/apidoc/"
	} else {
		if config.apidocDir[len(config.apidocDir)-1] != '/' {
			config.apidocDir += "/"
		}
	}
	if config.useGateWayPlugin || config.useGqlPlugin {
		_, err := os.Stat(config.apidocDir)
		if os.IsNotExist(err) {
			err = os.Mkdir(config.apidocDir, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

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
func gengql() {
	// 完整路径
	compath, err := filepath.Abs(config.genpath)
	if err != nil {
		log.Panicln(err)
	}
	// mod名
	err = os.Chdir(config.genpath)
	if err != nil {
		log.Panicln(err)
	}
	out, err := execx.RunGetOut("go list -m")
	if err != nil {
		log.Panicln(err)
	}
	mods := strings.Split(out, "\n")
	mod := mods[len(mods)-1]
	// 调用方mod路径
	out, err = execx.RunGetOut("go list -m -f {{.Dir}}")
	// 如果生成路径包含模块名
	_, after, _ := strings.Cut(compath, out)
	gomod := strings.ReplaceAll(mod+after, "\\", "/")
	packages := getPackages(compath)
	gqldir := config.apidocDir
	fileInfos, err := os.ReadDir(gqldir)
	if err != nil {
		log.Panicln(err)
	}
	for i := range fileInfos {
		if fileInfos[i].IsDir() {
			files, err := os.ReadDir(gqldir + fileInfos[i].Name())
			if err != nil {
				log.Panicln(err)
			}
			for j := range files {
				if strings.HasSuffix(files[j].Name(), ".graphql") {
					os.Chdir(gqldir + fileInfos[i].Name())
					/*			data, err := os.ReadFile(fileInfos[i].Name() + ".graphql")
								if err != nil {
									return
								}
								dataStr := stringsx.ToString(data)
								dataStr = strings.ReplaceAll(dataStr, ": Int", ": Int!")
								dataStr = strings.ReplaceAll(dataStr, ": String", ": String!")
								dataStr = strings.ReplaceAll(dataStr, ": Boolean", ": Boolean!")
								dataStr = strings.ReplaceAll(dataStr, ": Float", ": Float!")*/

					//这里用模板生成yml
					t := template.Must(template.New("yml").Parse(ymlTpl))
					config := fileInfos[i].Name() + `.gqlgen.yml`
					file, err := os.Create(config)
					if err != nil {
						log.Panicln(err)
					}
					renderValue := map[string]any{"GoMod": gomod, "SubDir": fileInfos[i].Name(), "Packages": packages}
					err = t.Execute(file, renderValue)
					if err != nil {
						log.Panicln(err)
					}
					file.Close()
					execx.RunWithLog(`gqlgen --verbose --config ` + config)
					break
				}
			}
		}
	}
}

func protocCmd(plugins []string, file, mod, modDir string) {
	cmd := "protoc " + config.include + " " + file
	var args string

	for _, plugin := range plugins {
		genpath := config.genpath
		if strings.HasPrefix(plugin, "openapiv2_out") {
			plugin += mod
			genpath = config.apidocDir + modDir
			err := fs.MkdirAll(genpath)
			if err != nil {
				log.Panicln(err)
			}
		}

		if strings.HasPrefix(plugin, "gql_out") {
			genpath = config.apidocDir
		}
		args += " --" + plugin + ":" + genpath

	}
	cmd += args
	if runtime.GOOS != "windows" {
		cmd = "bash -c \"" + cmd + "\""
	}
	execx.RunWithLog(cmd)
}
