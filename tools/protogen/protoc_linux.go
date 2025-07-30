/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package main

import (
	execi "github.com/hopeio/gox/os/exec"
)

// 直接执行protoc linux下会报 Could not make proto path relative: /xxx/*.proto: No such file or directory,找不到原因，无解
func protoc(cmd string) {
	cmd = "bash -c \"" + cmd + "\""
	execi.RunWithLog(cmd)
}
