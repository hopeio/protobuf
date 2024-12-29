/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package time

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/hopeio/utils/encoding/binary"
	stringsi "github.com/hopeio/utils/strings"
	timei "github.com/hopeio/utils/time"
	"io"
	"time"
)

func (ts *Date) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*ts = Date{Days: int32(nullTime.Time.Unix() / int64(timei.DaySecond))}
	return
}

func (ts Date) Value() (driver.Value, error) {
	return []byte(ts.Time().Format(time.DateOnly)), nil
}

func (ts Date) GormDataType() string {
	return "time"
}

func (ts Date) Time() time.Time {
	return time.Unix(int64(ts.Days)*int64(timei.DaySecond), 0)
}

func (ts *Date) MarshalBinary() ([]byte, error) {
	return binary.ToBinary(ts.Days), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (ts *Date) UnmarshalBinary(data []byte) error {
	ts.Days = binary.BinaryTo[int32](data)
	return nil
}

func (ts *Date) GobEncode() ([]byte, error) {
	return ts.MarshalBinary()
}

func (ts *Date) GobDecode(data []byte) error {
	return ts.UnmarshalBinary(data)
}

func (ts *Date) MarshalJSON() ([]byte, error) {
	if ts == nil {
		return []byte("null"), nil
	}
	t := ts.Time()
	if y := t.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(time.DateOnly)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, time.DateOnly)
	b = append(b, '"')
	return b, nil
}

func (ts *Date) UnmarshalJSON(data []byte) error {
	str := string(data)
	if len(str) == 0 || str == "null" {
		return nil
	}
	t, err := time.Parse(time.DateOnly, str[1:len(str)-1])
	if err != nil {
		return err
	}
	ts.Days = int32(t.Unix() / int64(timei.DaySecond))
	return nil
}
func (x *Date) MarshalGQL(w io.Writer) {
	w.Write([]byte(x.Time().Format(time.DateOnly)))
}

func (x *Date) UnmarshalGQL(v interface{}) error {
	if i, ok := v.(string); ok {
		t, err := time.Parse(time.DateOnly, i)
		if err != nil {
			return err
		}
		x.Days = int32(t.Unix() / int64(timei.DaySecond))
		return nil
	}
	return errors.New("enum need integer type")
}

func (d Date) MarshalText() ([]byte, error) {
	return stringsi.ToBytes(d.Time().Format(time.DateOnly)), nil
}

func (d *Date) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	str := stringsi.BytesToString(data)
	t, err := time.Parse(time.DateOnly, str)
	if err != nil {
		return err
	}
	d.Days = int32(t.Unix() / int64(timei.DaySecond))
	return nil
}
