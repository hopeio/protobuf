/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package main

type Config struct {
	proto, genpath, currentDir                                                     string
	thirdIncludes                                                                  []string
	include                                                                        string
	useEnumPlugin, useGateWayPlugin, useValidatorOutPlugin, useGqlPlugin, stdPatch bool
	apidocDir                                                                      string
}

var config = Config{}
