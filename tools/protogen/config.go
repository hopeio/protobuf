/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package main

type Config struct {
	proto, genpath, currentDir string
	thirdIncludes              []string
	include                    string
	useGqlPlugin               bool
	useOpenapiPlugin           bool
	apidocDir                  string
}

type Goconfig struct {
	useEnumPlugin, useGateWayPlugin, useValidatorOutPlugin, stdPatch bool
}

var config = Config{}
var goconfig = Goconfig{}
