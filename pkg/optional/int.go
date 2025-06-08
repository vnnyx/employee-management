package optional

import (
	"database/sql"
	"fmt"
	"strconv"
)

type Int64 struct {
	Option[int64]
}

func NewInt64(values ...int64) Int64 {
	var i Int64
	i.SetEmpty()

	if len(values) > 0 {
		i.Set(values[0])
	}

	return i
}

func NewInt64FromRef(value *int64) Int64 {
	var i Int64
	i.SetEmpty()

	if value != nil {
		i.Set(*value)
	}

	return i
}

func (i Int64) Value() (int64, bool) {
	v, ok := i.Get()
	if !i.IsValueSet() || !ok {
		return 0, false
	}
	return v, true
}

func (i *Int64) Scan(value any) error {
	sqlInt := sql.NullInt64{}
	err := sqlInt.Scan(value)
	if err != nil {
		return err
	}

	if sqlInt.Valid {
		i.Set(sqlInt.Int64)
	}

	return nil
}

func (i *Int64) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" {
		i.SetEmpty()
		return nil
	}
	v, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return fmt.Errorf("optional.Int64: cannot parse '%s' as int64: %w", str, err)
	}
	i.Set(v)
	return nil
}

type Int32 struct {
	Option[int32]
}

func NewInt32(values ...int32) Int32 {
	var i Int32
	i.SetEmpty()

	if len(values) > 0 {
		i.Set(values[0])
	}

	return i
}

func NewInt32FromRef(value *int32) Int32 {
	var i Int32
	i.SetEmpty()

	if value != nil {
		i.Set(*value)
	}

	return i
}

func (i Int32) Value() (int32, bool) {
	v, ok := i.Get()
	if !i.IsValueSet() || !ok {
		return 0, false
	}
	return v, true
}

func (i *Int32) Scan(value any) error {
	sqlInt := sql.NullInt64{}
	err := sqlInt.Scan(value)
	if err != nil {
		return err
	}

	if sqlInt.Valid {
		i.Set(int32(sqlInt.Int64))
	}

	return nil
}

func (i *Int32) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" {
		i.SetEmpty()
		return nil
	}
	v, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return fmt.Errorf("optional.Int32: cannot parse '%s' as int32: %w", str, err)
	}
	i.Set(int32(v))
	return nil
}
