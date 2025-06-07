package optional

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

type Time struct {
	Option[time.Time]
}

func NewTime(values ...time.Time) Time {
	var t Time
	t.SetEmpty()

	if len(values) > 0 {
		t.Set(values[0])
	}

	return t
}

func NewTimeFromRef(value *time.Time) Time {
	var t Time
	t.SetEmpty()

	if value != nil {
		t.Set(*value)
	}

	return t
}

func (t Time) Value() (driver.Value, error) {
	v, ok := t.Get()
	if !t.IsValueSet() || !ok {
		return nil, nil
	}
	return v, nil
}

func (t *Time) Scan(value interface{}) error {
	sqlTime := sql.NullTime{}
	err := sqlTime.Scan(value)
	if err != nil {
		return err
	}

	if sqlTime.Valid {
		t.Set(sqlTime.Time)
	}

	return nil
}
