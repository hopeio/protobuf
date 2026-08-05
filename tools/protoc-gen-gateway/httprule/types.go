package httprule

import (
	"fmt"
	"strings"
)

type template struct {
	segments []segment
	verb     string
	template string
}

type segment interface {
	fmt.Stringer
	compile() (ops []op)
}

type wildcard struct{}

type deepWildcard struct{}

type literal string

type variable struct {
	path     string
	segments []segment
}

// String returns the wildcard token "*".
func (wildcard) String() string {
	return "*"
}

// String returns the deep-wildcard token "**".
func (deepWildcard) String() string {
	return "**"
}

// String returns the literal string value.
func (l literal) String() string {
	return string(l)
}

// String returns the variable in "{field=segments}" notation.
func (v variable) String() string {
	var segs []string
	for _, s := range v.segments {
		segs = append(segs, s.String())
	}
	return fmt.Sprintf("{%s=%s}", v.path, strings.Join(segs, "/"))
}

// String returns the full path template string including leading slash and optional verb suffix.
func (t template) String() string {
	var segs []string
	for _, s := range t.segments {
		segs = append(segs, s.String())
	}
	str := strings.Join(segs, "/")
	if t.verb != "" {
		str = fmt.Sprintf("%s:%s", str, t.verb)
	}
	return "/" + str
}
