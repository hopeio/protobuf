package any

import (
	errors "errors"
	strings "github.com/hopeio/gox/strings"
	io "io"
)

func (x Encoding) Comment() string {
	switch x {
	case Encoding_JSON:
		return "Encoding_JSON"
	}
	return ""
}

func (x Encoding) MarshalGQL(w io.Writer) {
	w.Write(strings.SimpleQuoteToBytes(x.String()))
}

func (x *Encoding) UnmarshalGQL(v interface{}) error {
	if i, ok := v.(int32); ok {
		*x = Encoding(i)
		return nil
	}
	return errors.New("enum need integer type")
}
