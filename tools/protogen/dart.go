package main

import (
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
		run(config.hopeProto)
		run(config.proto)
	},
}
