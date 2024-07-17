package main

import (
	"os"
	"testing"
)

func TestMV(t *testing.T) {
	t.Error(os.Rename("D:\\code\\hopeio\\hoper\\thirdparty\\protobuf\\hopeio", "D:\\code\\hopeio\\hoper\\thirdparty\\protobuf"))
}
