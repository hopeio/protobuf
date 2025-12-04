/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package deletedAt

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"strings"
	"time"

	timex "github.com/hopeio/gox/time"
	"github.com/jinzhu/now"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

func (x *DeletedAt) Time() time.Time {
	return time.Unix(x.Seconds, int64(x.Nanos))
}

func (x *DeletedAt) IsValid() bool {
	return x != nil && timex.Check(x) == 0
}

// Scan implements the Scanner interface.
func (x *DeletedAt) Scan(value interface{}) error {
	nullTime := &sql.NullTime{}
	err := nullTime.Scan(value)
	if err != nil {
		return err
	}
	if nullTime.Valid {
		*x = DeletedAt{Seconds: nullTime.Time.Unix(), Nanos: int32(nullTime.Time.Nanosecond())}
	}
	return nil
}

// Value implements the driver Valuer interface.
func (x *DeletedAt) Value() (driver.Value, error) {
	if x == nil || timex.Check(x) != 0 {
		return nil, nil
	}
	return time.Unix(x.Seconds, int64(x.Nanos)), nil
}

func (x *DeletedAt) GormDataType() string {
	return "time"
}

var (
	FlagDeleted = 1
	FlagActived = 0
)

func (*DeletedAt) QueryClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{gorm.SoftDeleteQueryClause{Field: f, ZeroValue: parseZeroValueTag(f)}}
}

func parseZeroValueTag(f *schema.Field) sql.NullString {
	if v, ok := f.TagSettings["ZEROVALUE"]; ok {
		if _, err := now.Parse(v); err == nil {
			return sql.NullString{String: v, Valid: true}
		}
	}
	return sql.NullString{Valid: false}
}

func (*DeletedAt) DeleteClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{gorm.SoftDeleteDeleteClause{Field: f, ZeroValue: parseZeroValueTag(f)}}
}

func (*DeletedAt) UpdateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{gorm.SoftDeleteUpdateClause{Field: f}}
}

func getTimeType(settings map[string]string) schema.TimeType {
	if settings["NANO"] != "" {
		return schema.UnixNanosecond
	}

	if settings["MILLI"] != "" {
		return schema.UnixMillisecond
	}

	fieldUnit := strings.ToUpper(settings["DELETEDATFIELDUNIT"])

	if fieldUnit == "NANO" {
		return schema.UnixNanosecond
	}

	if fieldUnit == "MILLI" {
		return schema.UnixMillisecond
	}

	return schema.UnixSecond
}

func (x *DeletedAt) MarshalGQL(w io.Writer) {
	data, _ := timex.MarshalText(x.Time())
	w.Write(data)
}

func (x *DeletedAt) UnmarshalGQL(v interface{}) error {
	var t time.Time
	if i, ok := v.(string); ok {
		err := timex.UnmarshalText(&t, []byte(i))
		if err != nil {
			return err
		}
		*x = DeletedAt{Seconds: t.Unix(), Nanos: int32(t.Nanosecond())}
		return nil
	}
	return errors.New("enum need integer type")
}

func (x *DeletedAt) MarshalJSON() ([]byte, error) {
	if x == nil {
		return []byte("null"), nil
	}
	return timex.MarshalJSON(x.Time())
}

func (x *DeletedAt) UnmarshalJSON(data []byte) error {
	var st time.Time
	if err := timex.UnmarshalJSON(&st, data); err != nil {
		return err
	}
	*x = DeletedAt{Seconds: st.Unix(), Nanos: int32(st.Nanosecond())}
	return nil
}

type DeletedAtInput = DeletedAt
