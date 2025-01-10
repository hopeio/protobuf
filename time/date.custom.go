/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package time

import (
	"database/sql/driver"
	timei "github.com/hopeio/utils/time"
	"io"
	"time"
)

func (ts *Date) Scan(value interface{}) error {
	return (*timei.Date)(&ts.Days).Scan(value)
}

func (ts Date) Value() (driver.Value, error) {
	return timei.Date(ts.Days).Value()
}

func (ts Date) GormDataType() string {
	return "time"
}

func (ts Date) Time() time.Time {
	return timei.Date(ts.Days).Time()
}

func (ts Date) MarshalBinary() ([]byte, error) {
	return timei.Date(ts.Days).MarshalBinary()
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (ts *Date) UnmarshalBinary(data []byte) error {
	return (*timei.Date)(&ts.Days).UnmarshalBinary(data)
}

func (ts Date) GobEncode() ([]byte, error) {
	return ts.MarshalBinary()
}

func (ts *Date) GobDecode(data []byte) error {
	return ts.UnmarshalBinary(data)
}

func (ts Date) MarshalJSON() ([]byte, error) {
	return timei.Date(ts.Days).MarshalJSON()
}

func (ts *Date) UnmarshalJSON(data []byte) error {
	return (*timei.Date)(&ts.Days).UnmarshalJSON(data)
}
func (x Date) MarshalGQL(w io.Writer) {
	timei.Date(x.Days).MarshalGQL(w)
}

func (x *Date) UnmarshalGQL(v interface{}) error {
	return (*timei.Date)(&x.Days).UnmarshalGQL(v)
}

func (d Date) MarshalText() ([]byte, error) {
	return timei.Date(d.Days).MarshalText()
}

func (d *Date) UnmarshalText(data []byte) error {
	return (*timei.Date)(&d.Days).UnmarshalText(data)
}
