package optional

import (
	"database/sql"
	"database/sql/driver"
	"strings"
)

type String struct {
	Option[string]
}

func NewString(values ...string) String {
	var s String
	s.SetEmpty()

	if len(values) > 0 {
		s.Set(values[0])
	}

	return s
}

func NewStringFromRef(value *string) String {
	var s String
	s.SetEmpty()

	if value != nil {
		s.Set(*value)
	}

	return s
}

func (s String) Value() (driver.Value, error) {
	str, ok := s.Get()
	if !s.IsValueSet() || !ok {
		return nil, nil
	}
	return str, nil
}

func (s *String) Scan(value any) error {
	sqlStr := sql.NullString{}
	err := sqlStr.Scan(value)
	if err != nil {
		return err
	}

	if sqlStr.Valid {
		s.Set(sqlStr.String)
	}

	return nil
}

func (s *String) TrimSpace() String {
	if s.IsValueSet() {
		s.Set(strings.TrimSpace(*s.value))
	}

	return *s
}

func (s *String) TrimAllSpace() String {
	if s.IsPresent() {
		s.Set(strings.Join(strings.Fields(*s.value), ""))
	}

	return *s
}