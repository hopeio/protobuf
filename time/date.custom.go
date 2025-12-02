/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package time

import (
	"database/sql/driver"
	"io"
	"time"

	timex "github.com/hopeio/gox/time"
)

func (x *Date) Scan(value interface{}) error {
	return (*timex.Date)(&x.Days).Scan(value)
}

func (x *Date) Value() (driver.Value, error) {
	return timex.Date(x.Days).Value()
}

func (x *Date) GormDataType() string {
	return "time"
}

func (x *Date) Time() time.Time {
	return timex.Date(x.Days).Time()
}

func (x *Date) MarshalBinary() ([]byte, error) {
	return timex.Date(x.Days).MarshalBinary()
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (x *Date) UnmarshalBinary(data []byte) error {
	return (*timex.Date)(&x.Days).UnmarshalBinary(data)
}

func (x *Date) GobEncode() ([]byte, error) {
	return x.MarshalBinary()
}

func (x *Date) GobDecode(data []byte) error {
	return x.UnmarshalBinary(data)
}

func (x *Date) MarshalJSON() ([]byte, error) {
	return timex.Date(x.Days).MarshalJSON()
}

func (x *Date) UnmarshalJSON(data []byte) error {
	return (*timex.Date)(&x.Days).UnmarshalJSON(data)
}
func (x *Date) MarshalGQL(w io.Writer) {
	timex.Date(x.Days).MarshalGQL(w)
}

func (x *Date) UnmarshalGQL(v interface{}) error {
	return (*timex.Date)(&x.Days).UnmarshalGQL(v)
}

func (x *Date) MarshalText() ([]byte, error) {
	return timex.Date(x.Days).MarshalText()
}

func (x *Date) UnmarshalText(data []byte) error {
	return (*timex.Date)(&x.Days).UnmarshalText(data)
}
