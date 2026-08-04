package gengateway

import (
	descriptorx "github.com/hopeio/protobuf/tools/protoc-gen-gateway/descriptor"
	"strings"
	"testing"
)

func TestParseFramework(t *testing.T) {
	cases := []struct {
		in   string
		want Framework
		ok   bool
	}{
		{"gin", FrameworkGin, true},
		{"GIN", FrameworkGin, true},
		{"fiber", FrameworkFiber, true},
		{"nethttp", FrameworkNetHTTP, true},
		{"net/http", FrameworkNetHTTP, true},
		{"http", FrameworkNetHTTP, true},
		{"", FrameworkNetHTTP, true},
		{"  gin  ", FrameworkGin, true},
		{"chi", "", false},
	}
	for _, tc := range cases {
		got, err := ParseFramework(tc.in)
		if tc.ok {
			if err != nil {
				t.Fatalf("ParseFramework(%q): %v", tc.in, err)
			}
			if got != tc.want {
				t.Fatalf("ParseFramework(%q)=%q want %q", tc.in, got, tc.want)
			}
		} else if err == nil {
			t.Fatalf("ParseFramework(%q) want error", tc.in)
		}
	}
}

func TestNewFrameworkImports(t *testing.T) {
	cases := []struct {
		fw   Framework
		want []string
	}{
		{FrameworkGin, []string{"github.com/gin-gonic/gin", "github.com/hopeio/mix/contrib/gin"}},
		{FrameworkFiber, []string{"github.com/gofiber/fiber/v3", "github.com/hopeio/mix/contrib/fiber"}},
		{FrameworkNetHTTP, []string{"net/http", "github.com/hopeio/mix/gateway"}},
	}
	for _, tc := range cases {
		reg := descriptorx.NewRegistry()
		gen := New(reg, false, "Handler", true, false, tc.fw)
		if gen == nil {
			t.Fatalf("New(%q) returned nil", tc.fw)
		}
		got, err := applyTemplate(param{
			File:               crossLinkFixture(minimalServiceFile(t)),
			RegisterFuncSuffix: "Handler",
			AllowPatchFeature:  true,
			Framework:          tc.fw,
		}, reg)
		if err != nil {
			t.Fatalf("framework %s applyTemplate: %v", tc.fw, err)
		}
		for _, imp := range tc.want {
			if !strings.Contains(got, imp) {
				t.Fatalf("framework %s: want import %q in output", tc.fw, imp)
			}
		}
	}
}
