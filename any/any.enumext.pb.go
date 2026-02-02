package any

import (
	errors "errors"
	"fmt"
	io "io"

	strings "github.com/hopeio/gox/strings"
)

func (x Encoding) Comment() string {
	switch x {
	case Encoding_JSON:
		return "Encoding_JSON"
	}
	return ""
}

func (x Encoding) MarshalGQL(w io.Writer) {
	w.Write(strings.ToBytes(fmt.Sprintf(`"%s"`, x.String())))
}

func (x *Encoding) UnmarshalGQL(v interface{}) error {
	if i, ok := v.(int32); ok {
		*x = Encoding(i)
		return nil
	}
	return errors.New("enum need integer type")
}
