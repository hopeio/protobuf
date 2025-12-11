package main

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/hopeio/gox/log"
	execx "github.com/hopeio/gox/os/exec"
	"github.com/spf13/cobra"
)

const (
	gogqlOut      = "gogql_out=svc=true,merge=true,paths=source_relative"
	goOut         = "go-patch_out=plugin=go,paths=source_relative"
	grpcOut       = "go-patch_out=plugin=go-grpc,paths=source_relative"
	enumOut       = "enum_out=paths=source_relative"
	gatewayOut    = "grpc-gin_out=paths=source_relative"
	validatorsOut = "validator_out=paths=source_relative"
)

const (
	DepGrpcGateway = "github.com/grpc-ecosystem/grpc-gateway/v2"
	DepProtopatch  = "github.com/alta/protopatch"
)

var goCmd = &cobra.Command{
	Use: "go",
	PreRun: func(cmd *cobra.Command, args []string) {
		plugins = append(plugins, goOut, grpcOut)
		if goconfig.useEnumPlugin {
			plugins = append(plugins, enumOut)
		}
		if goconfig.useGateWayPlugin {
			plugins = append(plugins, gatewayOut)
		}
		if goconfig.useValidatorOutPlugin {
			plugins = append(plugins, validatorsOut)
		}
		if config.useGqlPlugin {
			plugins = append(plugins, gogqlOut)
		}
		if config.useGqlPlugin || config.useOpenapiPlugin {
			_, err := os.Stat(config.apidocDir)
			if os.IsNotExist(err) {
				err = os.MkdirAll(config.apidocDir, os.ModePerm)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		run(config.proto)
		if config.useGqlPlugin {
			gengql()
		}
	},
}

func init() {
	pflag := goCmd.Flags()
	pflag.BoolVarP(&goconfig.useEnumPlugin, "enum", "e", false, "是否使用enum扩展插件")
	pflag.BoolVarP(&goconfig.useGateWayPlugin, "gateway", "w", false, "是否使用grpc-gateway插件")
	pflag.BoolVarP(&goconfig.useValidatorOutPlugin, "validator", "v", false, "是否使用validators插件")
	pflag.BoolVar(&goconfig.stdPatch, "patch", false, "是否使用原生protopatch")
}

func gengql() {
	// 完整路径
	compath, err := filepath.Abs(config.genpath)
	if err != nil {
		log.Panic(err)
	}
	// mod名
	err = os.Chdir(config.genpath)
	if err != nil {
		log.Panic(err)
	}
	out, err := execx.RunGetOut("go list -m")
	if err != nil {
		log.Panic(err)
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
		log.Panic(err)
	}
	for i := range fileInfos {
		if fileInfos[i].IsDir() {
			files, err := os.ReadDir(gqldir + fileInfos[i].Name())
			if err != nil {
				log.Panic(err)
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
						log.Panic(err)
					}
					renderValue := map[string]any{"GoMod": gomod, "SubDir": fileInfos[i].Name(), "Packages": packages}
					err = t.Execute(file, renderValue)
					if err != nil {
						log.Panic(err)
					}
					file.Close()
					execx.RunWithLog(`gqlgen --verbose --config ` + config)
					break
				}
			}
		}
	}
}
