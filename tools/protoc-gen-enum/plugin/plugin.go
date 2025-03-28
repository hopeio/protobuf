/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package plugin

import (
	"github.com/hopeio/protobuf/tools/protoc-gen-enum/options"
	"google.golang.org/protobuf/compiler/protogen"
	"strconv"
	"strings"
)

type pathType int

const (
	pathTypeImport pathType = iota
	pathTypeSourceRelative
)

type Builder struct {
	plugin        *protogen.Plugin
	importStatus  protogen.GoImportPath
	importCodes   protogen.GoImportPath
	importStrings protogen.GoImportPath
	importStrconv protogen.GoImportPath
	importErrcode protogen.GoImportPath
	importErrors  protogen.GoImportPath
	importIo      protogen.GoImportPath
}

func NewBuilder(gen *protogen.Plugin) *Builder {
	return &Builder{
		plugin:        gen,
		importStatus:  "google.golang.org/grpc/status",
		importCodes:   "google.golang.org/grpc/codes",
		importStrings: "github.com/hopeio/utils/strings",
		importStrconv: "strconv",
		importErrcode: "github.com/hopeio/utils/errors/errcode",
		importErrors:  "errors",
		importIo:      "io",
	}
}

func parseParameter(param string) map[string]string {
	paramMap := make(map[string]string)

	for _, p := range strings.Split(param, ",") {
		if i := strings.Index(p, "="); i < 0 {
			paramMap[p] = ""
		} else {
			paramMap[p[0:i]] = p[i+1:]
		}
	}

	return paramMap
}

func (b *Builder) Generate() error {

	genFileMap := make(map[string]*protogen.GeneratedFile)

	for _, protoFile := range b.plugin.Files {
		if !protoFile.Generate {
			continue
		}

		if options.NoExtGenAll(protoFile) || len(protoFile.Enums) == 0 {
			continue
		}

		fileName := protoFile.GeneratedFilenamePrefix
		g := b.plugin.NewGeneratedFile(fileName+".enumext.pb.go", protoFile.GoImportPath)
		genFileMap[fileName] = g
		// third traverse: build associations
		for _, enum := range protoFile.Enums {
			if options.NoExtGen(enum) {
				genFileMap[fileName] = g
				break
			}
		}

	}

	for _, protoFile := range b.plugin.Files {
		fileName := protoFile.GeneratedFilenamePrefix
		g, ok := genFileMap[fileName]
		if !ok || len(protoFile.Enums) == 0 {
			continue
		}

		g.P("package ", protoFile.GoPackageName)

		for _, enum := range protoFile.Enums {
			b.generate(protoFile, enum, g)
		}

	}

	return nil
}

func (b *Builder) generate(f *protogen.File, e *protogen.Enum, g *protogen.GeneratedFile) {

	b.generateText(f, e, g)
	if options.EnabledEnumJsonMarshal(f, e) {
		b.generateJsonMarshal(e, g)
	}
	if options.EnabledEnumErrCode(e) {
		b.generateErrCode(e, g)
	}
	if options.EnabledEnumGqlGen(f, e) {
		b.generateGQLMarshal(e, g)
	}
}

func (b *Builder) generateText(f *protogen.File, e *protogen.Enum, g *protogen.GeneratedFile) {
	noEnumPrefix := options.FileOptions(f).GetNoEnumPrefix()
	ccTypeName := e.GoIdent
	if len(e.Values) >= 64 {
		g.P("var (")
		g.P(ccTypeName, "_text = map[", ccTypeName, "]string{")
		for _, ev := range e.Values {
			opts := options.EnumValueOptions(ev)
			name := opts.GetName()
			if name == "" {
				name = ev.GoIdent.GoName
				if noEnumPrefix {
					name = replacePrefix(ev.GoIdent.GoName, e.GoIdent.GoName+"_", "")
				}
			}
			if text := options.GetEnumText(ev); text != "" {
				g.P(name, " :", strconv.Quote(text), ",")
			} else {
				g.P(name, " :", strconv.Quote(name), ",")
			}
		}
		g.P("}")
		g.P(")")
		g.P()
	}
	g.P("func (x ", ccTypeName, ") Text() string {")
	if len(e.Values) >= 64 {
		g.P("return ", ccTypeName, "_text[x]")
	} else {
		g.P("switch x {")
		for _, ev := range e.Values {
			opts := options.EnumValueOptions(ev)
			name := opts.GetName()
			if name == "" {
				name = ev.GoIdent.GoName
				if noEnumPrefix {
					name = replacePrefix(ev.GoIdent.GoName, e.GoIdent.GoName+"_", "")
				}
			}

			//PrintComments(e.Comments, g)

			g.P("case ", name, " :")
			if text := options.GetEnumText(ev); text != "" {
				g.P("return ", strconv.Quote(text))
			} else {
				g.P("return ", strconv.Quote(name))
			}

		}
		g.P("}")
		g.P("return ", strconv.Quote(""))
	}
	g.P("}")
	g.P()
}

func (b *Builder) generateGQLMarshal(e *protogen.Enum, g *protogen.GeneratedFile) {
	ccTypeName := e.GoIdent

	typ := "uint32"
	if typ1 := options.GetEnumType(e); typ1 != "" {
		typ = typ1
	}
	g.P("func (x ", ccTypeName, ") MarshalGQL(w ", b.importIo.Ident("Writer"), ") {")
	g.P(`w.Write(`, b.importStrings.Ident("QuoteToBytes"), `(x.String()))`)
	g.P("}")
	g.P()
	g.P("func (x *", ccTypeName, ") UnmarshalGQL(v interface{}) error {")
	g.P(`if i, ok := v.(`, typ, "); ok {")
	g.P(`*x = `, ccTypeName, `(i)`)
	g.P("return nil")
	g.P("}")
	g.P(`return `, b.importErrors.Ident("New"), `("enum need integer type")`)
	g.P("}")
	g.P()
}

func (b *Builder) generateJsonMarshal(e *protogen.Enum, g *protogen.GeneratedFile) {
	ccTypeName := e.GoIdent

	g.P("func (x ", ccTypeName, ") MarshalJSON() ([]byte, error) {")
	g.P("return ", b.importStrings.Ident("QuoteToBytes"), "(x.String())", ", nil")
	g.P("}")
	g.P()
	g.P("func (x *", ccTypeName, ") UnmarshalJSON(data []byte) error {")

	g.P("if len(data) > 0 && data[0] == '\"' {")

	g.P("value, ok := ", ccTypeName, `_value[string(data[1:len(data)-1])]`)
	g.P("if ok {")

	g.P("*x = ", ccTypeName, "(value)")
	g.P("return nil")
	g.P("}")
	g.P("} else {")
	g.P("value, err := ", b.importStrconv.Ident("ParseInt"), `(string(data), 10, 32)`)
	g.P("if err == nil {")
	g.P("_, ok := ", ccTypeName, `_name[int32(value)]`)
	g.P("if ok {")
	g.P("*x = ", ccTypeName, "(value)")
	g.P("return nil")
	g.P("}")
	g.P("}")
	g.P("}")
	g.P(`return `, b.importErrors.Ident("New"), `("invalid enum value: `, ccTypeName, `")`)

	g.P("}")
	g.P()
}

func (b *Builder) generateErrCode(e *protogen.Enum, g *protogen.GeneratedFile) {
	ccTypeName := e.GoIdent

	g.P("func (x ", ccTypeName, ") Error() string {")
	g.P(`return x.Text()`)
	g.P("}")
	g.P()
	g.P("func (x ", ccTypeName, ") ErrRep() *", b.importErrcode.Ident("ErrRep"), " {")
	g.P(`return &errcode.ErrRep{Code: errcode.ErrCode(x), Msg: x.Text()}`)
	g.P("}")
	g.P()
	g.P("func (x ", ccTypeName, ") Msg(msg string) *", b.importErrcode.Ident("ErrRep"), " {")
	g.P(`return &errcode.ErrRep{Code: errcode.ErrCode(x), Msg: msg}`)
	g.P("}")
	g.P()
	g.P("func (x ", ccTypeName, ") Wrap(err error) *", b.importErrcode.Ident("ErrRep"), " {")
	g.P(`return &errcode.ErrRep{Code: errcode.ErrCode(x), Msg: err.Error()}`)
	g.P("}")
	g.P()

	g.P("func (x ", ccTypeName, ") GRPCStatus() *", b.importStatus.Ident("Status"), " {")
	g.P(`return `, `status.New(`, b.importCodes.Ident("Code"), `(x), x.Text())`)
	g.P("}")
	g.P()

	g.P("func (x ", ccTypeName, ") ErrCode() errcode.ErrCode {")
	g.P(`return errcode.ErrCode(x)`)
	g.P("}")
	g.P()

	g.P("func init() {")
	g.P("for code := range ", ccTypeName, "_name {")
	g.P("errcode.Register(errcode.ErrCode(code), ", ccTypeName, "(code).Text())")
	g.P(`}`)
	g.P("}")
	g.P()
}

func replacePrefix(s, prefix, with string) string {
	return with + strings.TrimPrefix(s, prefix)
}
