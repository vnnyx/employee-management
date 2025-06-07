package optional

import (
	"database/sql"
	"database/sql/driver"
)

type Bool struct {
	Option[bool]
}

func NewBool(values ...bool) Bool {
	var b Bool
	b.SetEmpty()

	if len(values) > 0 {
		b.Set(values[0])
	}

	return b
}

func NewBoolFromRef(value *bool) Bool {
	var b Bool
	b.SetEmpty()

	if value != nil {
		b.Set(*value)
	}

	return b
}

func (b Bool) Value() (driver.Value, error) {
	v, ok := b.Get()
	if !b.IsValueSet() || !ok {
		return nil, nil
	}
	return v, nil
}

func (b *Bool) Scan(value any) error {
	sqlBool := sql.NullBool{}
	err := sqlBool.Scan(value)
	if err != nil {
		return err
	}

	if sqlBool.Valid {
		b.Set(sqlBool.Bool)
	}

	return nil
}