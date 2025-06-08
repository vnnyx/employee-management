package optional

import (
	"reflect"

	"github.com/goccy/go-json"
)

func ToRef[T any](v T) *T {
	return &v
}

type Option[T any] struct {
	value *T
	isSet bool
}

// IsValueSet returns true if the value is set programmatically
// or set from external source like JSON, SQL, etc.
func (o *Option[T]) IsValueSet() bool {
	return o.isSet
}

// Set sets the value of the option.
func (o *Option[T]) Set(v T) {
	o.value = &v
	o.isSet = true
}

// SetEmpty sets the value of the option to nil.
func (o *Option[T]) SetEmpty() Option[T] {
	var value *T
	o.value = value
	o.isSet = true

	return *o
}

// IsPresent returns true if the value is not nil.
func (o Option[T]) IsPresent() bool {
	return o.value != nil
}

/*
IfPresent calls the function fn if the value is not nil

Example:

	var opt optional.Option[int]

	opt.IfPresent(func(v int) {
		fmt.Println(v)
	})
*/
func (o Option[T]) IfPresent(fn func(T)) Option[T] {
	if o.IsPresent() {
		fn(*o.value)
	}
	return o
}

// Get returns the value and a boolean indicating if the value is present.
func (o Option[T]) Get() (T, bool) {
	if !o.IsPresent() {
		var zeroValue T
		return zeroValue, false
	}

	return *o.value, true
}

// GetOrDefault returns the value if it is present, otherwise it returns the default value.
func (o Option[T]) GetOrDefault(defaultValues ...T) T {
	var defaultValue T
	if len(defaultValues) > 0 {
		defaultValue = defaultValues[0]
	}

	if !o.IsPresent() {
		return defaultValue
	}

	return *o.value
}

// MustGet returns the value if it is present, otherwise it panics.
func (o Option[T]) MustGet() T {
	if !o.IsPresent() {
		panic("value is not present")
	}

	return *o.value
}

// MarshalJSON if value is nil, it returns null, otherwise it returns the value.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.IsPresent() {
		return json.Marshal(*o.value)
	}

	return json.Marshal(nil)
}

// UnmarshalJSON sets the value from the JSON data.
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.SetEmpty()
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	o.Set(value)

	return nil
}

// MarshalEncryption serializes the underlying value to JSON.
func (o Option[T]) MarshalEncryption() ([]byte, error) {
	if !o.IsPresent() {
		return []byte("null"), nil
	}
	return json.Marshal(*o.value)
}

// UnmarshalEncryption deserializes the JSON data and sets the underlying value.
func (o *Option[T]) UnmarshalEncryption(data []byte) error {
	// If data is "null", treat it as an empty value.
	if string(data) == "null" {
		o.SetEmpty()
		return nil
	}
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	o.Set(v)
	return nil
}

func (o Option[T]) IsZero() bool {
	if !o.isSet || o.value == nil {
		return true
	}
	var zero T
	return reflect.DeepEqual(*o.value, zero)
}

func (o Option[T]) SetAndReturn(v T) Option[T] {
	o.Set(v)
	return o
}
