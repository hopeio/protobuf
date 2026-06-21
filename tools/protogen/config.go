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
}

type Goconfig struct {
	useGateWayPlugin, useValidatorOutPlugin bool
}

var config = Config{}
var goconfig = Goconfig{}
