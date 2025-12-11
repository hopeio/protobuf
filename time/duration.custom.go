/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package time

import (
	"time"

	timex "github.com/hopeio/gox/time"
	"google.golang.org/protobuf/runtime/protoimpl"
)

// AsDuration converts x to a time.Duration,
// returning the closest duration value in the event of overflow.
func (x *Duration) AsDuration() time.Duration {
	return time.Duration(x.Micro) * 1000
}

// IsValid reports whether the duration is valid.
// It is equivalent to CheckValid == nil.
func (x *Duration) IsValid() bool {
	return x.check() == 0
}

// CheckValid returns an error if the duration is invalid.
// In particular, it checks whether the value is within the range of
// -10000 years to +10000 years inclusive.
// An error is reported for a nil Duration.
func (x *Duration) CheckValid() error {
	switch x.check() {
	case invalidNil:
		return protoimpl.X.NewError("invalid nil Duration")
	case invalidUnderflow:
		return protoimpl.X.NewError("duration (%v) exceeds -10000 years", x)
	case invalidOverflow:
		return protoimpl.X.NewError("duration (%v) exceeds +10000 years", x)
	case invalidNanosRange:
		return protoimpl.X.NewError("duration (%v) has out-of-range nanos", x)
	case invalidNanosSign:
		return protoimpl.X.NewError("duration (%v) has seconds and nanos with different signs", x)
	default:
		return nil
	}
}

const (
	_ = iota
	invalidNil
	invalidUnderflow
	invalidOverflow
	invalidNanosRange
	invalidNanosSign
)

func (x *Duration) check() uint {
	const absDuration = 315576000000000000 // 10000yr * 365.25day/yr * 24hr/day * 60min/hr * 60sec/min * 1e6
	switch {
	case x.Micro < -absDuration:
		return invalidUnderflow
	case x.Micro > +absDuration:
		return invalidOverflow
	default:
		return 0
	}
}

func (t Duration) MarshalJSON() ([]byte, error) {
	return timex.Duration(t.Micro * 1000).MarshalJSON()
}

// UnmarshalJSON implements the [encoding/json.Unmarshaler] interface.
// The time must be a quoted string in the RFC 3339 format.
func (t *Duration) UnmarshalJSON(data []byte) error {
	var t1 timex.Duration
	err := (*timex.Duration)(&t1).UnmarshalJSON(data)
	if err != nil {
		return err
	}
	t.Micro = int64(t1 / 1000)
	return nil
}
