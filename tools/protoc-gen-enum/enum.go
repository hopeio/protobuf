package main

import (
	"github.com/hopeio/protobuf/tools/protoc-gen-enum/plugin"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	protogen.Options{}.Run(func(p *protogen.Plugin) error {
		b := plugin.NewBuilder(p)
		return b.Generate()
	})

}
