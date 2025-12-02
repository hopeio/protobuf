/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package time

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/hopeio/gox/encoding/binary"
	timex "github.com/hopeio/gox/time"
)

func (ts *Timestamp) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*ts = Timestamp{Timestamp: nullTime.Time.UnixMilli()}
	return
}

func (ts *Timestamp) Value() (driver.Value, error) {
	return time.UnixMilli(ts.Timestamp), nil
}

func (ts *Timestamp) Time() time.Time {
	return time.UnixMilli(ts.Timestamp)
}

func (ts *Timestamp) GormDataType() string {
	return "time"
}

func (ts *Timestamp) MarshalBinary() ([]byte, error) {
	return binary.ToBinary(ts.Timestamp), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (ts *Timestamp) UnmarshalBinary(data []byte) error {
	ts.Timestamp = binary.BinaryTo[int64](data)
	return nil
}

func (ts *Timestamp) GobEncode() ([]byte, error) {
	return ts.MarshalBinary()
}

func (ts *Timestamp) GobDecode(data []byte) error {
	return ts.UnmarshalBinary(data)
}

func (ts *Timestamp) MarshalJSON() ([]byte, error) {
	if ts == nil {
		return []byte("null"), nil
	}
	t := time.UnixMilli(ts.Timestamp)
	return timex.MarshalJSON(t)
}

func (ts *Timestamp) UnmarshalJSON(data []byte) error {
	var t time.Time
	err := timex.UnmarshalJSON(&t, data)
	if err != nil {
		return err
	}
	ts.Timestamp = t.UnixMilli()
	return err
}
