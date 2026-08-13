/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package main

type Config struct {
	proto, genpath, hopeProto string
	thirdIncludes             []string
	include                   string
	useOpenapiPlugin          bool
	apidocDir                 string
	args                      []string
	// jobs 控制并行 protoc 进程数，默认 NumCPU
	jobs int
}

type Goconfig struct {
	useGateWayPlugin, useValidatorOutPlugin bool
	gatewayFramework                        string
}

var config = Config{}
var goconfig = Goconfig{}
