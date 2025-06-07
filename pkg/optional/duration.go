package optional

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Duration struct {
	Option[time.Duration]
}

func NewDuration(values ...time.Duration) Duration {
	var d Duration
	d.SetEmpty()

	if len(values) > 0 {
		d.Set(values[0])
	}

	return d
}

func NewDurationFromRef(value *time.Duration) Duration {
	var d Duration
	d.SetEmpty()

	if value != nil {
		d.Set(*value)
	}

	return d
}

func (d Duration) Value() (driver.Value, error) {
	if !d.IsValueSet() {
		return nil, nil
	}

	v, ok := d.Get()
	if !ok {
		return nil, nil
	}

	return v, nil
}

func (d *Duration) Scan(value any) error {
	if value == nil {
		d.SetEmpty()
		return nil
	}

	switch v := value.(type) {
	case time.Duration:
		d.Set(v)
		return nil
	case string:
		parsed, err := parsePostgresInterval(v)
		if err != nil {
			return err
		}
		d.Set(parsed)
		return nil
	case []byte:
		parsed, err := parsePostgresInterval(string(v))
		if err != nil {
			return err
		}
		d.Set(parsed)
		return nil
	default:
		return errors.New("unsupported type for Duration scan")
	}
}

func parsePostgresInterval(s string) (time.Duration, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid interval format: %s", s)
	}

	formatted := fmt.Sprintf("%sh%sm%ss", parts[0], parts[1], parts[2])
	duration, err := time.ParseDuration(formatted)
	if err != nil {
		return 0, fmt.Errorf("failed to parse interval: %w", err)
	}

	return duration, nil
}
