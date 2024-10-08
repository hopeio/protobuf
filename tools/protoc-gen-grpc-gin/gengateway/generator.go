package gengateway

import (
	"errors"
	"fmt"
	descriptor2 "github.com/hopeio/protobuf/tools/protoc-gen-grpc-gin/descriptor"
	"github.com/hopeio/utils/log"
	"go/format"
	"path"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	errNoTargetService = errors.New("no target service defined in the file")
)

type Generator interface {
	// Generate generates output files from input .proto files.
	Generate(targets []*descriptor2.File) ([]*descriptor2.ResponseFile, error)
}

type generator struct {
	reg                *descriptor2.Registry
	baseImports        []descriptor2.GoPackage
	useRequestContext  bool
	registerFuncSuffix string
	allowPatchFeature  bool
	standalone         bool
}

// New returns a new generator which generates grpc gateway files.
func New(reg *descriptor2.Registry, useRequestContext bool, registerFuncSuffix string,
	allowPatchFeature, standalone bool) Generator {
	var imports []descriptor2.GoPackage
	for _, pkgpath := range []string{
		"context",
		"io",

		"google.golang.org/protobuf/proto",
		"google.golang.org/grpc",
		"google.golang.org/grpc/codes",
		"google.golang.org/grpc/grpclog",
		"google.golang.org/grpc/metadata",
		"google.golang.org/grpc/status",
		"github.com/gin-gonic/gin",
		"github.com/hopeio/utils/net/http/gin/binding",
		"github.com/hopeio/utils/net/http/grpc",
		"github.com/hopeio/utils/net/http/grpc/gateway/gin",
	} {
		pkg := descriptor2.GoPackage{
			Path: pkgpath,
			Name: path.Base(pkgpath),
		}
		if err := reg.ReserveGoPackageAlias(pkg.Name, pkg.Path); err != nil {
			for i := 0; ; i++ {
				alias := fmt.Sprintf("%s_%d", pkg.Name, i)
				if err := reg.ReserveGoPackageAlias(alias, pkg.Path); err != nil {
					continue
				}
				pkg.Alias = alias
				break
			}
		}
		imports = append(imports, pkg)
	}

	return &generator{
		reg:                reg,
		baseImports:        imports,
		useRequestContext:  useRequestContext,
		registerFuncSuffix: registerFuncSuffix,
		allowPatchFeature:  allowPatchFeature,
		standalone:         standalone,
	}
}

func (g *generator) Generate(targets []*descriptor2.File) ([]*descriptor2.ResponseFile, error) {
	var files []*descriptor2.ResponseFile
	for _, file := range targets {
		log.Infof("Processing %s", file.GetName())

		code, err := g.generate(file)
		if err == errNoTargetService {
			log.Infof("%s: %v", file.GetName(), err)
			continue
		}
		if err != nil {
			return nil, err
		}
		formatted, err := format.Source([]byte(code))
		if err != nil {
			log.Errorf("%v: %s", err, code)
			return nil, err
		}
		files = append(files, &descriptor2.ResponseFile{
			GoPkg: file.GoPkg,
			CodeGeneratorResponse_File: &pluginpb.CodeGeneratorResponse_File{
				Name:    proto.String(file.GeneratedFilenamePrefix + ".pb.gw.go"),
				Content: proto.String(string(formatted)),
			},
		})
	}
	return files, nil
}

func (g *generator) generate(file *descriptor2.File) (string, error) {
	pkgSeen := make(map[string]bool)
	var imports []descriptor2.GoPackage
	for _, pkg := range g.baseImports {
		pkgSeen[pkg.Path] = true
		imports = append(imports, pkg)
	}

	if g.standalone {
		imports = append(imports, file.GoPkg)
	}

	for _, svc := range file.Services {
		for _, m := range svc.Methods {
			imports = append(imports, g.addEnumPathParamImports(file, m, pkgSeen)...)
			pkg := m.RequestType.File.GoPkg
			if len(m.Bindings) == 0 ||
				pkg == file.GoPkg || pkgSeen[pkg.Path] {
				continue
			}
			pkgSeen[pkg.Path] = true
			imports = append(imports, pkg)
		}
	}
	params := param{
		File:               file,
		Imports:            imports,
		UseRequestContext:  g.useRequestContext,
		RegisterFuncSuffix: g.registerFuncSuffix,
		AllowPatchFeature:  g.allowPatchFeature,
	}
	if g.reg != nil {
		params.OmitPackageDoc = g.reg.GetOmitPackageDoc()
	}
	return applyTemplate(params, g.reg)
}

// addEnumPathParamImports handles adding import of enum path parameter go packages
func (g *generator) addEnumPathParamImports(file *descriptor2.File, m *descriptor2.Method, pkgSeen map[string]bool) []descriptor2.GoPackage {
	var imports []descriptor2.GoPackage
	for _, b := range m.Bindings {
		for _, p := range b.PathParams {
			e, err := g.reg.LookupEnum("", p.Target.GetTypeName())
			if err != nil {
				continue
			}
			pkg := e.File.GoPkg
			if pkg == file.GoPkg || pkgSeen[pkg.Path] {
				continue
			}
			pkgSeen[pkg.Path] = true
			imports = append(imports, pkg)
		}
	}
	return imports
}
