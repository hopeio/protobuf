package request

import (
	errors "errors"
	strings "github.com/hopeio/gox/strings"
	io "io"
)

func (x SortType) Comment() string {
	switch x {
	case SortType_ASC:
		return "SortType_ASC"
	case SortType_DESC:
		return "SortType_DESC"
	}
	return ""
}

func (x SortType) MarshalGQL(w io.Writer) {
	w.Write(strings.SimpleQuoteToBytes(x.String()))
}

func (x *SortType) UnmarshalGQL(v interface{}) error {
	if i, ok := v.(int32); ok {
		*x = SortType(i)
		return nil
	}
	return errors.New("enum need integer type")
}
