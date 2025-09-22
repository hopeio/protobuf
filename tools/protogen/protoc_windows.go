/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package main

import (
	execx "github.com/hopeio/gox/os/exec"
)

func protoc(cmd string) {
	execx.RunWithLog(cmd)
}
