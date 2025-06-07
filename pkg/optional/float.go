package optional

import "database/sql"

type Float64 struct {
	Option[float64]
}

func NewFloat64(values ...float64) Float64 {
	var f Float64
	f.SetEmpty()

	if len(values) > 0 {
		f.Set(values[0])
	}

	return f
}

func NewFloat64FromRef(value *float64) Float64 {
	var f Float64
	f.SetEmpty()

	if value != nil {
		f.Set(*value)
	}

	return f
}

func (f Float64) Value() (float64, bool) {
	v, ok := f.Get()
	if !f.IsValueSet() || !ok {
		return 0, false
	}
	return v, true
}

func (f *Float64) Scan(value any) error {
	sqlFloat := sql.NullFloat64{}
	err := sqlFloat.Scan(value)
	if err != nil {
		return err
	}

	if sqlFloat.Valid {
		f.Set(sqlFloat.Float64)
	}

	return nil
}

type Float32 struct {
	Option[float32]
}

func NewFloat32(values ...float32) Float32 {
	var f Float32
	f.SetEmpty()

	if len(values) > 0 {
		f.Set(values[0])
	}

	return f
}

func NewFloat32FromRef(value *float32) Float32 {
	var f Float32
	f.SetEmpty()

	if value != nil {
		f.Set(*value)
	}

	return f
}

func (f Float32) Value() (float32, bool) {
	v, ok := f.Get()
	if !f.IsValueSet() || !ok {
		return 0, false
	}
	return v, true
}

func (f *Float32) Scan(value any) error {
	sqlFloat := sql.NullFloat64{}
	err := sqlFloat.Scan(value)
	if err != nil {
		return err
	}

	if sqlFloat.Valid {
		f.Set(float32(sqlFloat.Float64))
	}

	return nil
}