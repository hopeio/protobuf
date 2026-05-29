package gengateway

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/hopeio/gox/log"
	descriptor2 "github.com/hopeio/protobuf/tools/protoc-gen-grpc-gin/descriptor"

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
}

type binding struct {
	*descriptor2.Binding
	Registry          *descriptor2.Registry
	AllowPatchFeature bool
}

// GetBodyFieldPath returns the binding body's fieldpath.
func (b binding) GetBodyFieldPath() string {
	if b.Body != nil && len(b.Body.FieldPath) != 0 {
		return b.Body.FieldPath.String()
	}
	return "*"
}

// GetBodyFieldStructName returns the binding body's struct field name.
func (b binding) GetBodyFieldStructName() (string, error) {
	if b.Body != nil && len(b.Body.FieldPath) != 0 {
		return stringsx.SnakeToCamel(b.Body.FieldPath.String()), nil
	}
	return "", errors.New("no body field found")
}

// HasQueryParam determines if the binding needs parameters in query string.
//
// It sometimes returns true even though actually the binding does not need.
// But it is not serious because it just results in a small amount of extra codes generated.
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

// HasEnumPathParam returns true if the path parameter slice contains a parameter
// that maps to an enum proto field that is not repeated, if not false is returned.
func (b binding) HasEnumPathParam() bool {
	return b.hasEnumPathParam(false)
}

// HasRepeatedEnumPathParam returns true if the path parameter slice contains a parameter
// that maps to a repeated enum proto field, if not false is returned.
func (b binding) HasRepeatedEnumPathParam() bool {
	return b.hasEnumPathParam(true)
}

// hasEnumPathParam returns true if the path parameter slice contains a parameter
// that maps to a enum proto field and that the enum proto field is or isn't repeated
// based on the provided 'repeated' parameter.
func (b binding) hasEnumPathParam(repeated bool) bool {
	for _, p := range b.PathParams {
		if p.IsEnum() && p.IsRepeated() == repeated {
			return true
		}
	}
	return false
}

// LookupEnum looks up a enum type by path parameter.
func (b binding) LookupEnum(p descriptor2.Parameter) *descriptor2.Enum {
	e, err := b.Registry.LookupEnum("", p.Target.GetTypeName())
	if err != nil {
		return nil
	}
	return e
}

// FieldMaskField returns the golang-style name of the variable for a FieldMask, if there is exactly one of that type in
// the message. Otherwise, it returns an empty string.
func (b binding) FieldMaskField() string {
	var fieldMaskField *descriptor2.Field
	for _, f := range b.Method.RequestType.Fields {
		if f.GetTypeName() == ".google.protobuf.FieldMask" {
			// if there is more than 1 FieldMask for this request, then return none
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

// queryParamFilter is a wrapper of utilities.DoubleArray which provides String() to output DoubleArray.Encoding in a stable and predictable format.
type queryParamFilter struct {
	*utilities.DoubleArray
}

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
}

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
			for _, b := range meth.Bindings {
				methodWithBindingsSeen = true
				// Local
				if err := localHandlerTemplate.Execute(w, binding{
					Binding:           b,
					Registry:          reg,
					AllowPatchFeature: p.AllowPatchFeature,
				}); err != nil {
					return "", err
				}
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
	}
	// Local
	if err := localTrailerTemplate.Execute(w, tp); err != nil {
		return "", err
	}
	return w.String(), nil
}

var (
	headerTemplate = template.Must(template.New("header").Parse(`
// Code generated by protoc-gen-grpc-gin. DO NOT EDIT.
// source: {{.GetName}}

{{if not .OmitPackageDoc}}/*
Package {{.GoPkg.Name}} is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/{{end}}
package {{.GoPkg.Name}}
import (
	{{range $i := .Imports}}{{if $i.Standard}}{{$i | printf "%s\n"}}{{end}}{{end}}

	{{range $i := .Imports}}{{if not $i.Standard}}{{$i | printf "%s\n"}}{{end}}{{end}}
)

// Suppress "imported and not used" errors
var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = metadata.Join
var _ = gox.Pointer[bool]
var _ = strings.Bool
`))

	localHandlerTemplate = template.Must(template.New("local-handler").Parse(`
{{if and .Method.GetClientStreaming .Method.GetServerStreaming}}
{{else if .Method.GetClientStreaming}}
{{else if .Method.GetServerStreaming}}
{{else}}
{{template "local-client-rpc-request-func" .}}
{{end}}
`))

	_ = template.Must(localHandlerTemplate.New("local-request-func-signature").Parse(strings.Replace(`
{{if .Method.GetServerStreaming}}
{{else}}
func local_request_{{.Method.Service.GetName}}_{{.Method.GetName}}_{{.Index}}(server {{.Method.Service.InstanceName}}Server, ctx *gin.Context) (proto.Message, grpc_0.ServerMetadata, error)
{{end}}`, "\n", "", -1)))

	_ = template.Must(localHandlerTemplate.New("local-client-rpc-request-func").Parse(`
{{$AllowPatchFeature := .AllowPatchFeature}}
{{template "local-request-func-signature" .}} {
	var stream grpc_0.ServerTransportStream
	var protoReq {{.Method.RequestType.GoType .Method.Service.File.GoPkg.Path}}
{{if or (or .Body .HasQueryParam) (and (ne .HTTPMethod "GET") (ne .HTTPMethod "DELETE"))}}

	if err := gateway.Bind(ctx, &protoReq); err != nil {
		return nil, stream.ServerMetadata(), err
	}
{{end}}
	ctx.Request = ctx.Request.WithContext(grpc.NewContextWithServerTransportStream(metadata.NewIncomingContext(ctx.Request.Context(), metadata.MD(ctx.Request.Header)), &stream))
{{if .PathParams}}
	var (
		err error
	{{- if .HasEnumPathParam}}
		e int32
	{{- end}}
	{{- if .HasRepeatedEnumPathParam}}
		es []int32
	{{- end}}
	)

	{{$binding := .}}
	{{range $param := .PathParams}}
	{{$enum := $binding.LookupEnum $param}}


{{if $param.IsNestedProto3}}
	err = gin_1.PopulateFieldFromPath(&protoReq, {{$param | printf "%q"}}, ctx.Param({{$param | printf "%q"}}))
	if err != nil {
		return nil, stream.ServerMetadata(), status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", {{$param | printf "%q"}}, err)
	}
	{{if $enum}}
		e{{if $param.IsRepeated}}s{{end}}, err = {{$param.ConvertFuncExpr}}(ctx.Param({{$param | printf "%q"}}){{if $param.IsRepeated}}, {{$binding.Registry.GetRepeatedPathParamSeparator | printf "%c" | printf "%q"}}{{end}}, {{$enum.GoType $param.Method.Service.File.GoPkg.Path}}_value)
		if err != nil {
			return nil, stream.ServerMetadata(), status.Errorf(codes.InvalidArgument, "could not parse path as enum value, parameter: %s, error: %v", {{$param | printf "%q"}}, err)
		}
	{{end}}
{{else if $enum}}
	e{{if $param.IsRepeated}}s{{end}}, err = {{$param.ConvertFuncExpr}}(ctx.Param({{$param | printf "%q"}}){{if $param.IsRepeated}}, {{$binding.Registry.GetRepeatedPathParamSeparator | printf "%c" | printf "%q"}}{{end}}, {{$enum.GoType $param.Method.Service.File.GoPkg.Path}}_value)
	if err != nil {
		return nil, stream.ServerMetadata(), status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", {{$param | printf "%q"}}, err)
	}
{{else}}
	{{$param.AssignableExpr "protoReq"}}, err = {{$param.ConvertFuncExpr}}(ctx.Param({{$param | printf "%q"}}){{if $param.IsRepeated}}, {{$binding.Registry.GetRepeatedPathParamSeparator | printf "%c" | printf "%q"}}{{end}})
	if err != nil {
		return nil, stream.ServerMetadata(), status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", {{$param | printf "%q"}}, err)
	}
{{end}}
{{if and $enum $param.IsRepeated}}
	s := make([]{{$enum.GoType $param.Method.Service.File.GoPkg.Path}}, len(es))
	for i, v := range es {
		s[i] = {{$enum.GoType $param.Method.Service.File.GoPkg.Path}}(v)
	}
	{{$param.AssignableExpr "protoReq"}} = s
{{else if $enum}}
	{{$param.AssignableExpr "protoReq"}} = {{$enum.GoType $param.Method.Service.File.GoPkg.Path}}(e)
{{end}}
	{{end}}
{{end}}

{{if .Method.GetServerStreaming}}
	// TODO
{{else}}
	resp, err := server.{{.Method.GetName}}(ctx.Request.Context(), &protoReq)
	return resp, stream.ServerMetadata(), err
{{end}}
}`))

	localTrailerTemplate = template.Must(template.New("local-trailer").Parse(`
{{$UseRequestContext := .UseRequestContext}}
{{range $svc := .Services}}
// Register{{$svc.GetName}}{{$.RegisterFuncSuffix}}Server registers the http handlers for service {{$svc.GetName}} to "mux".
// UnaryRPC     :call {{$svc.GetName}}Server directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
// Note that using this registration option will cause many gRPC library features to stop working. Consider using Register{{$svc.GetName}}{{$.RegisterFuncSuffix}}FromEndpoint instead.
func Register{{$svc.GetName}}{{$.RegisterFuncSuffix}}Server(mux *gin.Engine, server {{$svc.InstanceName}}Server) {
	{{range $m := $svc.Methods}}
	{{range $b := $m.Bindings}}
	{{if or $m.GetClientStreaming $m.GetServerStreaming}}
	mux.Handle({{$b.HTTPMethod | printf "%q"}}, {{$b.PathTmpl.Template | printf "%q"}}, func(ctx *gin.Context) {
		err := status.Error(codes.Unimplemented, "streaming calls are not yet supported in the in-process transport")
		gateway.HttpError(ctx, err)
		return
	})
	{{else}}
	mux.Handle({{$b.HTTPMethod | printf "%q"}}, {{$b.PathTmpl.Template | printf "%q"}}, func(ctx *gin.Context) {
		resp, md, err := local_request_{{$svc.GetName}}_{{$m.GetName}}_{{$b.Index}}(server, ctx)
		if err !=nil {
			gateway.HttpError(ctx, err)
			return
		}
		{{ if $b.ResponseBody }}
		gateway.ForwardResponseMessage(ctx, md, response_{{$svc.GetName}}_{{$m.GetName}}_{{$b.Index}}{resp})
		{{ else }}
		gateway.ForwardResponseMessage(ctx, md, resp)
		{{end}}
	})
	{{end}}
	{{end}}
	{{end}}
}
{{end}}`))
)
