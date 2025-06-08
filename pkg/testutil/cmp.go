package testutil

import (
	"log"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/vnnyx/employee-management/pkg/optional"
)

func EqualVerbose(x interface{}, y interface{}, opts ...cmp.Option) bool {
	opts = append(opts,
		comparerForOptional[int64](),
		comparerForOptional[int32](),
		comparerForOptional[float64](),
		comparerForOptional[float32](),
		comparerForOptional[bool](),
		comparerForOptional[string](),
		comparerForOptional[time.Duration](),

		comparerForWrapper[int64, optional.Int64](),
		comparerForWrapper[int32, optional.Int32](),
		comparerForWrapper[float64, optional.Float64](),
		comparerForWrapper[float32, optional.Float32](),
		comparerForWrapper[bool, optional.Bool](),
		comparerForWrapper[string, optional.String](),
		comparerForWrapper[time.Duration, optional.Duration](),
	)

	if diff := cmp.Diff(x, y, opts...); diff != "" {
		log.Println(diff)
	}
	return cmp.Equal(x, y, opts...)
}

func comparerForOptional[T comparable]() cmp.Option {
	return cmp.Comparer(func(a, b optional.Option[T]) bool {
		aVal, aOk := a.Get()
		bVal, bOk := b.Get()
		if aOk != bOk {
			return false
		}
		if !aOk {
			return true
		}
		return aVal == bVal
	})
}

func comparerForWrapper[T comparable, W interface {Get() (T, bool)}]() cmp.Option {
	return cmp.Comparer(func(a, b W) bool {
		av, aok := a.Get()
		bv, bok := b.Get()
		if aok != bok {
			return false
		}
		if !aok {
			return true
		}
		return av == bv
	})
}
