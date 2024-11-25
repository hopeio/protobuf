/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package main

import (
	"os"
	"testing"
)

func TestMV(t *testing.T) {
	t.Error(os.Rename("D:\\code\\hopeio\\hoper\\thirdparty\\protobuf\\hopeio", "D:\\code\\hopeio\\hoper\\thirdparty\\protobuf"))
}
