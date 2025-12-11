/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package duration

import (
	"math"
	"time"

	timex "github.com/hopeio/gox/time"
	"google.golang.org/protobuf/runtime/protoimpl"
)

// New constructs a new Duration from the provided time.Duration.
func New(d time.Duration) *Duration {
	nanos := d.Nanoseconds()
	secs := nanos / 1e9
	nanos -= secs * 1e9
	return &Duration{Seconds: secs, Nanos: int32(nanos)}
}

// AsDuration converts x to a time.Duration,
// returning the closest duration value in the event of overflow.
func (x *Duration) AsDuration() time.Duration {
	secs := x.GetSeconds()
	nanos := x.GetNanos()
	d := time.Duration(secs) * time.Second
	overflow := d/time.Second != time.Duration(secs)
	d += time.Duration(nanos) * time.Nanosecond
	overflow = overflow || (secs < 0 && nanos < 0 && d > 0)
	overflow = overflow || (secs > 0 && nanos > 0 && d < 0)
	if overflow {
		switch {
		case secs < 0:
			return time.Duration(math.MinInt64)
		case secs > 0:
			return time.Duration(math.MaxInt64)
		}
	}
	return d
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
	const absDuration = 315576000000 // 10000yr * 365.25day/yr * 24hr/day * 60min/hr * 60sec/min
	secs := x.GetSeconds()
	nanos := x.GetNanos()
	switch {
	case x == nil:
		return invalidNil
	case secs < -absDuration:
		return invalidUnderflow
	case secs > +absDuration:
		return invalidOverflow
	case nanos <= -1e9 || nanos >= +1e9:
		return invalidNanosRange
	case (secs > 0 && nanos < 0) || (secs < 0 && nanos > 0):
		return invalidNanosSign
	default:
		return 0
	}
}

func (t Duration) MarshalJSON() ([]byte, error) {
	return timex.Duration(t.AsDuration()).MarshalJSON()
}

// UnmarshalJSON implements the [encoding/json.Unmarshaler] interface.
// The time must be a quoted string in the RFC 3339 format.
func (t *Duration) UnmarshalJSON(data []byte) error {
	var t1 timex.Duration
	if err := t1.UnmarshalJSON(data); err != nil {
		return err
	}
	nanos := time.Duration(t1).Nanoseconds()
	t.Seconds = nanos / 1e9
	t.Nanos -= int32(t.Seconds * 1e9)
	return nil
}
