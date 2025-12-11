package main

import (
	"runtime"

	"github.com/spf13/cobra"
)

const (
	tsOut = `ts_proto_out=globalThisPolyfill=true,esModuleInterop=true,importSuffix=.js`
)

const (
	tsPluginPath = `./node_modules/.bin/protoc-gen-ts_proto`
)

var tsCmd = &cobra.Command{
	Use: "ts",
	Run: func(cmd *cobra.Command, args []string) {
		plugins = append(plugins, tsOut)
		if runtime.GOOS == "windows" {
			pluginPaths = append(pluginPaths, `protoc-gen-ts_proto=".\node_modules\.bin\protoc-gen-ts_proto.cmd"`)
		} else {
			pluginPaths = append(pluginPaths, tsPluginPath)
		}
		run(config.proto)
	},
}
