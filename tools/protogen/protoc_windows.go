package main

import (
	execi "github.com/hopeio/utils/os/exec"
)

func protoc(cmd string) {
	execi.Run(cmd)
}
