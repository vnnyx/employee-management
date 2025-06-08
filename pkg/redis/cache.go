package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/vnnyx/employee-management/pkg/observability/instrumentation"
)

func SetWithExpiration(ctx context.Context, key string, value any, expiration time.Duration) error {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"redis.SetWithExpiration()",
	)
	defer span.End()

	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return Client.Set(ctx, key, jsonValue, expiration).Err()
}

func GetAndUnmarshal[T any](ctx context.Context, key string) (T, error) {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"redis.GetAndUnmarshal()",
	)
	defer span.End()

	val, err := Client.Get(ctx, key).Result()
	if err != nil {
		var zero T
		return zero, err
	}

	var result T
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return result, err
	}

	return result, nil
}
