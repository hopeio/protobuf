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
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

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
		// 前置检查 protoc，缺失时给出可执行的安装指引，而不是每个目录报一遍晦涩错误
		if _, err := exec.LookPath("protoc"); err != nil {
			log.Fatal("protoc not found in PATH; install it first (https://protobuf.dev/installation/) " +
				"and run `go run github.com/hopeio/protobuf/tools` to install code generators")
		}
		if err := os.MkdirAll(config.genpath, os.ModePerm); err != nil {
			log.Fatal(err)
		}
		if config.useOpenapiPlugin {
			plugins = append(plugins, openapiv2Out)
			if config.apidocDir == "" {
				var err error
				config.apidocDir, err = filepath.Abs(config.genpath + "/../apidoc")
				if err != nil {
					config.apidocDir = config.genpath + "/apidoc"
				}
			} else if config.apidocDir[len(config.apidocDir)-1] != '/' {
				config.apidocDir += "/"
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
	pflag.StringVarP(&config.genpath, "output", "o", pwd, "generate dir")
	pflag.StringVarP(&config.hopeProto, "protobuf", "p", "/proto", "手动指定hopeio/proto项目目录")
	pflag.StringArrayVarP(&config.thirdIncludes, "include", "I", nil, "外部proto依赖")
	pflag.BoolVarP(&config.useOpenapiPlugin, "openapi", "d", false, "是否使用openapi插件")
	pflag.StringVar(&config.apidocDir, "apiDocDir", "", "api doc目录")
	pflag.IntVarP(&config.jobs, "jobs", "j", runtime.NumCPU(), "并行执行的 protoc 进程数")

	rootCmd.AddCommand(goCmd)
	rootCmd.AddCommand(dartCmd)
	rootCmd.AddCommand(tsCmd)

}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// run walks dir, collecting one protoc invocation per directory that contains
// .proto files, then executes them in parallel. Any failure exits non-zero so
// CI cannot go green on partially generated code.
func run(dir string) {
	var cmds []string
	collectProtocCmds(dir, &cmds)
	if len(cmds) == 0 {
		return
	}
	jobs := config.jobs
	if jobs < 1 {
		jobs = 1
	}
	sem := make(chan struct{}, jobs)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error
	for _, cmd := range cmds {
		wg.Add(1)
		sem <- struct{}{}
		go func(cmd string) {
			defer wg.Done()
			defer func() { <-sem }()
			if err := execx.RunWithLog(cmd); err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("%s: %w", cmd, err))
				mu.Unlock()
			}
		}(cmd)
	}
	wg.Wait()
	if len(errs) > 0 {
		for _, err := range errs {
			log.Println("protoc failed:", err)
		}
		os.Exit(1)
	}
}

// collectProtocCmds recursively gathers protoc command lines for every directory containing .proto files.
func collectProtocCmds(dir string, cmds *[]string) {
	fileInfos, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalln(err)
	}
	var collected bool
	for i := range fileInfos {
		if !collected && strings.HasSuffix(fileInfos[i].Name(), ".proto") {
			collected = true
			if dir == config.proto {
				*cmds = append(*cmds, protocCmd(dir+"/*.proto", "root"))
			} else {
				*cmds = append(*cmds, protocCmd(dir+"/*.proto", filepath.Base(dir)))
			}
		}
		if fileInfos[i].IsDir() {
			collectProtocCmds(dir+"/"+fileInfos[i].Name(), cmds)
		}
	}
}
// getInclude resolves the hopeio/protobuf import path and builds the -I include flags for protoc.
func getInclude() {
	config.proto, _ = filepath.Abs(config.proto)
	config.genpath, _ = filepath.Abs(config.genpath)
	log.Println("proto:", config.proto)
	log.Println("genpath:", config.genpath)

	libDir, err := execx.RunGetOut(gox.GoListDir + DepProtobuf)
	if err == nil {
		config.hopeProto = libDir + "/_proto"
	}
	config.include = "-I" + config.hopeProto + " -I" + config.proto
	for _, include := range config.thirdIncludes {
		config.include += " -I" + include
	}

	log.Println("include:", config.include)
}

// protocCmd assembles the protoc command line for the given proto file glob and module name.
func protocCmd(file, mod string) string {
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
		args += " --" + plugin + ":" + genpath

	}
	cmd += args + " " + strings.Join(config.args, " ")
	if runtime.GOOS != "windows" {
		cmd = fmt.Sprintf(`bash -c "%s"`, cmd)
	}
	return cmd
}
