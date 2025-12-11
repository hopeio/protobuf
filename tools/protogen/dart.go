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

	},
}
