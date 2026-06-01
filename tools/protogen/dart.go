package main

import (
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

const (
	dartOut = "dart_out=grpc"
)

var dartCmd = &cobra.Command{
	Use:   "dart",
	Short: "generate dart code",
	Run: func(cmd *cobra.Command, args []string) {
		plugins = append(plugins, dartOut)
		if runtime.GOOS == "windows" {
			local := os.Getenv("LOCALAPPDATA")
			if local != "" {
				pluginPaths = append(pluginPaths, `protoc-gen-dart="`+local+`\Pub\Cache\bin\protoc-gen-dart.bat"`)
			}
		}
		run(config.hopeProto)
		run(config.proto)
	},
}
