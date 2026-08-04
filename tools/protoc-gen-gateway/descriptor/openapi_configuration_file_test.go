package descriptor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hopeio/protobuf/utils/openapiconfig"
)

func TestRegisterOpenAPIOptions_NilAndSuccess(t *testing.T) {
	reg := NewRegistry()
	if err := registerOpenAPIOptions(reg, &openapiconfig.OpenAPIConfig{}, "empty.yaml"); err != nil {
		t.Fatal(err)
	}
	if err := registerOpenAPIOptions(reg, &openapiconfig.OpenAPIConfig{OpenapiOptions: nil}, "nil-options.yaml"); err != nil {
		t.Fatal(err)
	}

	loadFile(t, reg, `
		name: 'a.proto'
		package: 'example.foo'
		options < go_package: 'github.com/hopeio/protobuf/tools/protoc-gen-gateway/testdata/example' >
		message_type < name: 'ExampleMessage' >
		service <
			name: "AService"
			method <
				name: "Meth"
				input_type: "ExampleMessage"
				output_type: "ExampleMessage"
				options <
					[google.api.http] < post: "/v1/a" body: "*" >
				>
			>
		>
	`)
	cfg := &openapiconfig.OpenAPIConfig{
		OpenapiOptions: &openapiconfig.OpenAPIOptions{
			File: []*openapiconfig.OpenAPIFileOption{{File: "a.proto"}},
		},
	}
	if err := registerOpenAPIOptions(reg, cfg, "ok.yaml"); err != nil {
		t.Fatal(err)
	}
	if _, ok := reg.GetOpenAPIFileOption("a.proto"); !ok {
		t.Fatal("file option not registered")
	}
}

func TestRegisterOpenAPIOptions_RegisterError(t *testing.T) {
	reg := NewRegistry()
	cfg := &openapiconfig.OpenAPIConfig{
		OpenapiOptions: &openapiconfig.OpenAPIOptions{
			File: []*openapiconfig.OpenAPIFileOption{{File: "missing.proto"}},
		},
	}
	err := registerOpenAPIOptions(reg, cfg, "bad.yaml")
	if err == nil || !strings.Contains(err.Error(), "failed to register option in bad.yaml") {
		t.Fatalf("want wrapped register error, got %v", err)
	}
}

func TestRegistry_LoadOpenAPIConfigFromYAML_File(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "openapi.yaml")
	content := []byte(`
openapiOptions:
  file:
  - file: a.proto
    option:
      schemes:
      - HTTP
      - HTTPS
`)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}

	reg := NewRegistry()
	loadFile(t, reg, `
		name: 'a.proto'
		package: 'example.foo'
		options < go_package: 'github.com/hopeio/protobuf/tools/protoc-gen-gateway/testdata/example' >
		message_type < name: 'ExampleMessage' >
		service <
			name: "AService"
			method <
				name: "Meth"
				input_type: "ExampleMessage"
				output_type: "ExampleMessage"
				options <
					[google.api.http] < post: "/v1/a" body: "*" >
				>
			>
		>
	`)
	if err := reg.LoadOpenAPIConfigFromYAML(path); err != nil {
		t.Fatal(err)
	}
	opt, ok := reg.GetOpenAPIFileOption("a.proto")
	if !ok || opt == nil || len(opt.Schemes) != 2 {
		t.Fatalf("file option=%v ok=%v", opt, ok)
	}
}

func TestRegistry_LoadOpenAPIConfigFromYAML_MissingFile(t *testing.T) {
	reg := NewRegistry()
	err := reg.LoadOpenAPIConfigFromYAML(filepath.Join(t.TempDir(), "missing.yaml"))
	if err == nil || !strings.Contains(err.Error(), "failed to read gRPC API Configuration description") {
		t.Fatalf("want read error, got %v", err)
	}
}

func TestRegistry_LoadOpenAPIConfigFromYAML_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	if err := os.WriteFile(path, []byte("openapiOptions:\n  file:\n  - - broken"), 0o644); err != nil {
		t.Fatal(err)
	}
	reg := NewRegistry()
	err := reg.LoadOpenAPIConfigFromYAML(path)
	if err == nil || !strings.Contains(err.Error(), "failed to parse gRPC API Configuration from YAML") {
		t.Fatalf("want parse error, got %v", err)
	}
}
