package main

import (
	"os"

	"github.com/hopeio/gox/log"
	"github.com/spf13/cobra"
)

const (
	goOut         = "go-patch_out=plugin=go,paths=source_relative"
	grpcOut       = "go-patch_out=plugin=go-grpc,paths=source_relative"
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
		if goconfig.useGateWayPlugin {
			plugins = append(plugins, gatewayOut)
		}
		if goconfig.useValidatorOutPlugin {
			plugins = append(plugins, validatorsOut)
		}
		if config.useOpenapiPlugin {
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
	},
}

func init() {
	pflag := goCmd.Flags()
	pflag.BoolVarP(&goconfig.useGateWayPlugin, "gateway", "w", false, "是否使用grpc-gateway插件")
	pflag.BoolVarP(&goconfig.useValidatorOutPlugin, "validator", "v", false, "是否使用validators插件")
}
