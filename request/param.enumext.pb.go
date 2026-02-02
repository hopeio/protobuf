package request

import (
	errors "errors"
	"fmt"
	io "io"

	strings "github.com/hopeio/gox/strings"
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
	w.Write(strings.ToBytes(fmt.Sprintf(`"%s"`, x.String())))
}

func (x *SortType) UnmarshalGQL(v interface{}) error {
	if i, ok := v.(int32); ok {
		*x = SortType(i)
		return nil
	}
	return errors.New("enum need integer type")
}
