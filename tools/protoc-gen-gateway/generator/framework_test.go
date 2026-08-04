package gengateway

import "testing"

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
