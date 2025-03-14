/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package main

import (
	execi "github.com/hopeio/utils/os/exec"
)

func protoc(cmd string) {
	execi.RunWithLog(cmd)
}
