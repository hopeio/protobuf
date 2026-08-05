package gengateway

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/hopeio/gox/log"
	descriptor2 "github.com/hopeio/protobuf/tools/protoc-gen-gateway/descriptor"

	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	stringsx "github.com/hopeio/gox/strings"
)

type param struct {
	*descriptor2.File
	Imports            []descriptor2.GoPackage
	UseRequestContext  bool
	RegisterFuncSuffix string
	AllowPatchFeature  bool
	OmitPackageDoc     bool
	Framework          Framework
}

type binding struct {
	*descriptor2.Binding
	Registry          *descriptor2.Registry
	AllowPatchFeature bool
}

// GetBodyFieldPath returns the body field path, or "*" when the entire message is the body.
func (b binding) GetBodyFieldPath() string {
	if b.Body != nil && len(b.Body.FieldPath) != 0 {
		return b.Body.FieldPath.String()
	}
	return "*"
}

// GetBodyFieldStructName returns the CamelCase Go struct field name for the body field path.
func (b binding) GetBodyFieldStructName() (string, error) {
	if b.Body != nil && len(b.Body.FieldPath) != 0 {
		return stringsx.SnakeToCamel(b.Body.FieldPath.String()), nil
	}
	return "", errors.New("no body field found")
}

// HasQueryParam reports whether the binding has any request fields that are passed as query parameters.
func (b binding) HasQueryParam() bool {
	if b.Body != nil && len(b.Body.FieldPath) == 0 {
		return false
	}
	fields := make(map[string]bool)
	for _, f := range b.Method.RequestType.Fields {
		fields[f.GetName()] = true
	}
	if b.Body != nil {
		delete(fields, b.Body.FieldPath.String())
	}
	for _, p := range b.PathParams {
		delete(fields, p.FieldPath.String())
	}
	return len(fields) > 0
}

// QueryParamFilter builds a DoubleArray filter excluding body and path parameters from query parsing.
func (b binding) QueryParamFilter() queryParamFilter {
	var seqs [][]string
	if b.Body != nil {
		seqs = append(seqs, strings.Split(b.Body.FieldPath.String(), "."))
	}
	for _, p := range b.PathParams {
		seqs = append(seqs, strings.Split(p.FieldPath.String(), "."))
	}
	return queryParamFilter{utilities.NewDoubleArray(seqs)}
}

// HasEnumPathParam reports whether any non-repeated path parameter is an enum type.
func (b binding) HasEnumPathParam() bool {
	return b.hasEnumPathParam(false)
}

// HasRepeatedEnumPathParam reports whether any repeated path parameter is an enum type.
func (b binding) HasRepeatedEnumPathParam() bool {
	return b.hasEnumPathParam(true)
}

// hasEnumPathParam checks whether any path parameter matches the enum+repeated criteria.
func (b binding) hasEnumPathParam(repeated bool) bool {
	for _, p := range b.PathParams {
		if p.IsEnum() && p.IsRepeated() == repeated {
			return true
		}
	}
	return false
}

// LookupEnum looks up the Enum descriptor for the parameter's type name, returning nil on failure.
func (b binding) LookupEnum(p descriptor2.Parameter) *descriptor2.Enum {
	e, err := b.Registry.LookupEnum("", p.Target.GetTypeName())
	if err != nil {
		return nil
	}
	return e
}

// FieldMaskField returns the CamelCase name of the single FieldMask field in the request type, or "" if ambiguous.
func (b binding) FieldMaskField() string {
	var fieldMaskField *descriptor2.Field
	for _, f := range b.Method.RequestType.Fields {
		if f.GetTypeName() == ".google.protobuf.FieldMask" {
			if fieldMaskField != nil {
				return ""
			}
			fieldMaskField = f
		}
	}
	if fieldMaskField != nil {
		return stringsx.SnakeToCamel(fieldMaskField.GetName())
	}
	return ""
}

type queryParamFilter struct {
	*utilities.DoubleArray
}

// String serializes the DoubleArray filter as a Go literal suitable for code generation.
func (f queryParamFilter) String() string {
	encodings := make([]string, len(f.Encoding))
	for str, enc := range f.Encoding {
		encodings[enc] = fmt.Sprintf("%q: %d", str, enc)
	}
	e := strings.Join(encodings, ", ")
	return fmt.Sprintf("&utilities.DoubleArray{Encoding: map[string]int{%s}, Base: %#v, Check: %#v}", e, f.Base, f.Check)
}

type trailerParams struct {
	Services           []*descriptor2.Service
	UseRequestContext  bool
	RegisterFuncSuffix string
	Framework          Framework
}

// toColonParams converts Go 1.22 style {param} to :param for gin/fiber routers.
func toColonParams(p string) string {
	p = strings.ReplaceAll(p, "{", ":")
	p = strings.ReplaceAll(p, "}", "")
	return p
}

// httpMethod converts an HTTP verb to title-case (e.g. "GET" → "Get") for use in template calls.
func httpMethod(m string) string {
	return stringsx.UpperCaseFirst(strings.ToLower(m))
}

// applyTemplate executes the header and trailer templates to produce complete gateway Go source.
func applyTemplate(p param, reg *descriptor2.Registry) (string, error) {
	w := bytes.NewBuffer(nil)
	if err := headerTemplate.Execute(w, p); err != nil {
		return "", err
	}
	var targetServices []*descriptor2.Service

	for _, msg := range p.Messages {
		msgName := stringsx.SnakeToCamel(*msg.Name)
		msg.Name = &msgName
	}

	for _, svc := range p.Services {
		var methodWithBindingsSeen bool
		svcName := stringsx.SnakeToCamel(*svc.Name)
		svc.Name = &svcName

		for _, meth := range svc.Methods {
			log.Infof("Processing %s.%s", svc.GetName(), meth.GetName())
			methName := stringsx.SnakeToCamel(*meth.Name)
			meth.Name = &methName
			for range meth.Bindings {
				methodWithBindingsSeen = true
			}
		}
		if methodWithBindingsSeen {
			targetServices = append(targetServices, svc)
		}
	}
	if len(targetServices) == 0 {
		return "", errNoTargetService
	}

	tp := trailerParams{
		Services:           targetServices,
		UseRequestContext:  p.UseRequestContext,
		RegisterFuncSuffix: p.RegisterFuncSuffix,
		Framework:          p.Framework,
	}
	if err := localTrailerTemplate.Execute(w, tp); err != nil {
		return "", err
	}
	return w.String(), nil
}

var (
	headerTemplate = template.Must(template.New("header").Parse(`
// Code generated by protoc-gen-gateway. DO NOT EDIT.
// source: {{.GetName}}

{{if not .OmitPackageDoc}}/*
Package {{.GoPkg.Name}} is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/{{end}}
package {{.GoPkg.Name}}
import (
{{- if eq .Framework "fiber"}}
	"github.com/gofiber/fiber/v3"
	gateway "github.com/hopeio/mix/contrib/fiber"
{{- else if eq .Framework "nethttp"}}
	"net/http"
	gateway "github.com/hopeio/mix/gateway"
{{- else}}
	"github.com/gin-gonic/gin"
	gateway "github.com/hopeio/mix/contrib/gin"
{{- end}}
)
`))

	templateFuncs        = template.FuncMap{"toColonParams": toColonParams, "httpMethod": httpMethod}
	localTrailerTemplate = template.Must(template.New("local-trailer").Funcs(templateFuncs).Parse(`
{{range $svc := .Services}}
// Register{{$svc.GetName}}{{$.RegisterFuncSuffix}}Server registers the http handlers for service {{$svc.GetName}}.
// UnaryRPC        : call {{$svc.GetName}}Server directly.
// ServerStreaming : HTTP/2 chunked with gRPC length-prefixed JSON frames.
// ClientStreaming : read gRPC length-prefixed JSON frames from request body.
// BidiStreaming   : bidirectional gRPC length-prefixed frames.
{{- if eq $.Framework "fiber"}}
func Register{{$svc.GetName}}{{$.RegisterFuncSuffix}}Server(app *fiber.App, server {{$svc.InstanceName}}Server) {
{{- else if eq $.Framework "nethttp"}}
func Register{{$svc.GetName}}{{$.RegisterFuncSuffix}}Server(mux *http.ServeMux, server {{$svc.InstanceName}}Server) {
{{- else}}
func Register{{$svc.GetName}}{{$.RegisterFuncSuffix}}Server(mux *gin.Engine, server {{$svc.InstanceName}}Server) {
{{- end}}
	{{range $m := $svc.Methods}}
	{{range $b := $m.Bindings}}
	{{if and $m.GetClientStreaming $m.GetServerStreaming}}
	{{- if eq $.Framework "fiber"}}
	app.{{httpMethod $b.HTTPMethod}}({{toColonParams $b.PathTmpl.Template | printf "%q"}}, gateway.BidiStreamCall(server.{{$m.GetName}}))
	{{- else if eq $.Framework "nethttp"}}
	mux.Handle({{printf "%s %s" $b.HTTPMethod $b.PathTmpl.Template | printf "%q"}}, gateway.BidiStreamCall(server.{{$m.GetName}}))
	{{- else}}
	mux.Handle({{$b.HTTPMethod | printf "%q"}}, {{toColonParams $b.PathTmpl.Template | printf "%q"}}, gateway.BidiStreamCall(server.{{$m.GetName}}))
	{{- end}}
	{{else if $m.GetClientStreaming}}
	{{- if eq $.Framework "fiber"}}
	app.{{httpMethod $b.HTTPMethod}}({{toColonParams $b.PathTmpl.Template | printf "%q"}}, gateway.ClientSideStreamCall(server.{{$m.GetName}}))
	{{- else if eq $.Framework "nethttp"}}
	mux.Handle({{printf "%s %s" $b.HTTPMethod $b.PathTmpl.Template | printf "%q"}}, gateway.ClientSideStreamCall(server.{{$m.GetName}}))
	{{- else}}
	mux.Handle({{$b.HTTPMethod | printf "%q"}}, {{toColonParams $b.PathTmpl.Template | printf "%q"}}, gateway.ClientSideStreamCall(server.{{$m.GetName}}))
	{{- end}}
	{{else if $m.GetServerStreaming}}
	{{- if eq $.Framework "fiber"}}
	app.{{httpMethod $b.HTTPMethod}}({{toColonParams $b.PathTmpl.Template | printf "%q"}}, gateway.ServerSideStreamCall(server.{{$m.GetName}}))
	{{- else if eq $.Framework "nethttp"}}
	mux.Handle({{printf "%s %s" $b.HTTPMethod $b.PathTmpl.Template | printf "%q"}}, gateway.ServerSideStreamCall(server.{{$m.GetName}}))
	{{- else}}
	mux.Handle({{$b.HTTPMethod | printf "%q"}}, {{toColonParams $b.PathTmpl.Template | printf "%q"}}, gateway.ServerSideStreamCall(server.{{$m.GetName}}))
	{{- end}}
	{{else}}
	{{- if eq $.Framework "fiber"}}
	app.{{httpMethod $b.HTTPMethod}}({{toColonParams $b.PathTmpl.Template | printf "%q"}}, gateway.UnaryCall(server.{{$m.GetName}}))
	{{- else if eq $.Framework "nethttp"}}
	mux.Handle({{printf "%s %s" $b.HTTPMethod $b.PathTmpl.Template | printf "%q"}}, gateway.UnaryCall(server.{{$m.GetName}}))
	{{- else}}
	mux.Handle({{$b.HTTPMethod | printf "%q"}}, {{toColonParams $b.PathTmpl.Template | printf "%q"}}, gateway.UnaryCall(server.{{$m.GetName}}))
	{{- end}}
	{{end}}
	{{end}}
	{{end}}
}
{{end}}`))
)
