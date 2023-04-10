package client

import (
	"time"
)

type Result[T any] struct {
	Result T
	Err    error
}

type ResultMap[T any, K comparable, V any] struct {
	Result T
	Err    error
	Map    map[K]V
}

func NullableString(str string) *string {
	if str == "" {
		return nil
	}
	return &str
}

func ExponentialBackoff(init, limit time.Duration) func() {
	internalTime := init
	return func() {
		time.Sleep(internalTime)
		switch {
		case internalTime*2 > limit:
			internalTime = limit
		case internalTime*2 <= limit:
			internalTime *= 2
		}
	}
}
