package descriptor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hopeio/protobuf/utils/apiconfig"
	"google.golang.org/genproto/googleapis/api/annotations"
)

func TestRegisterHTTPRulesFromGrpcAPIService(t *testing.T) {
	reg := NewRegistry()
	svc := &apiconfig.GrpcAPIService{
		Http: &annotations.Http{
			Rules: []*annotations.HttpRule{
				{
					Selector: " svc.Method ",
					Pattern:  &annotations.HttpRule_Post{Post: "/v1/x"},
				},
			},
		},
	}
	if err := registerHTTPRulesFromGrpcAPIService(reg, svc, "test.yaml"); err != nil {
		t.Fatal(err)
	}
	rules := reg.LookupExternalHTTPRules(".svc.Method")
	if len(rules) != 1 || rules[0].GetPost() != "/v1/x" {
		t.Fatalf("rules=%v", rules)
	}
}

func TestRegisterHTTPRules_WildcardRejected(t *testing.T) {
	reg := NewRegistry()
	svc := &apiconfig.GrpcAPIService{
		Http: &annotations.Http{
			Rules: []*annotations.HttpRule{
				{
					Selector: "*.Echo",
					Pattern:  &annotations.HttpRule_Get{Get: "/v1/echo"},
				},
			},
		},
	}
	err := registerHTTPRulesFromGrpcAPIService(reg, svc, "bad.yaml")
	if err == nil || !strings.Contains(err.Error(), "without wildcards") {
		t.Fatalf("want wildcard error, got %v", err)
	}
}

func TestRegisterHTTPRules_NilHTTP(t *testing.T) {
	reg := NewRegistry()
	if err := registerHTTPRulesFromGrpcAPIService(reg, &apiconfig.GrpcAPIService{}, "empty"); err != nil {
		t.Fatal(err)
	}
}

func TestLoadGrpcAPIServiceFromYAML_FileAndRegister(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "api.yaml")
	content := []byte(`
type: google.api.Service
config_version: 3
http:
  rules:
  - selector: grpctest.YourService.Echo
    post: /v1/myecho
    body: "*"
`)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}
	reg := NewRegistry()
	if err := reg.LoadGrpcAPIServiceFromYAML(path); err != nil {
		t.Fatal(err)
	}
	rules := reg.LookupExternalHTTPRules(".grpctest.YourService.Echo")
	if len(rules) != 1 || rules[0].GetPost() != "/v1/myecho" {
		t.Fatalf("rules=%v", rules)
	}
}

func TestLoadOpenAPIConfigFromYAML_RejectUnknownField(t *testing.T) {
	_, err := loadOpenAPIConfigFromYAML([]byte(`
openapiOptions:
  file:
  - file: a.proto
unknownField: 1
`), "unknown.yaml")
	if err == nil {
		t.Fatal("want error for unknown field with DiscardUnknown=false")
	}
}
