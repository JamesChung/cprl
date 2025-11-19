package client

import (
	"errors"
	"testing"
	"time"
)

func TestNullableString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantNil  bool
		wantVal  string
	}{
		{
			name:    "empty string returns nil",
			input:   "",
			wantNil: true,
		},
		{
			name:    "non-empty string returns pointer",
			input:   "test",
			wantNil: false,
			wantVal: "test",
		},
		{
			name:    "whitespace string returns pointer",
			input:   "   ",
			wantNil: false,
			wantVal: "   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NullableString(tt.input)
			if tt.wantNil {
				if got != nil {
					t.Errorf("NullableString(%q) = %v, want nil", tt.input, got)
				}
			} else {
				if got == nil {
					t.Errorf("NullableString(%q) = nil, want non-nil", tt.input)
				} else if *got != tt.wantVal {
					t.Errorf("NullableString(%q) = %q, want %q", tt.input, *got, tt.wantVal)
				}
			}
		})
	}
}

func TestExponentialBackoff(t *testing.T) {
	tests := []struct {
		name       string
		init       time.Duration
		limit      time.Duration
		iterations int
		wantDelays []time.Duration
	}{
		{
			name:       "doubles until limit",
			init:       100 * time.Millisecond,
			limit:      1 * time.Second,
			iterations: 5,
			wantDelays: []time.Duration{
				100 * time.Millisecond,
				200 * time.Millisecond,
				400 * time.Millisecond,
				800 * time.Millisecond,
				1 * time.Second, // hits limit
			},
		},
		{
			name:       "respects limit immediately",
			init:       500 * time.Millisecond,
			limit:      600 * time.Millisecond,
			iterations: 3,
			wantDelays: []time.Duration{
				500 * time.Millisecond,
				600 * time.Millisecond, // limit reached
				600 * time.Millisecond, // stays at limit
			},
		},
		{
			name:       "single iteration",
			init:       50 * time.Millisecond,
			limit:      1 * time.Second,
			iterations: 1,
			wantDelays: []time.Duration{
				50 * time.Millisecond,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backoff := ExponentialBackoff(tt.init, tt.limit)

			for i := 0; i < tt.iterations; i++ {
				start := time.Now()
				backoff()
				elapsed := time.Since(start)

				// Allow 10ms tolerance for timing variations
				tolerance := 10 * time.Millisecond
				expected := tt.wantDelays[i]

				if elapsed < expected-tolerance || elapsed > expected+tolerance {
					t.Errorf("iteration %d: backoff() took %v, want ~%v (±%v)",
						i, elapsed, expected, tolerance)
				}
			}
		})
	}
}

func TestExponentialBackoff_StatefulBehavior(t *testing.T) {
	// Test that the backoff function maintains state across calls
	backoff := ExponentialBackoff(10*time.Millisecond, 100*time.Millisecond)

	// First call should sleep for ~10ms
	start := time.Now()
	backoff()
	first := time.Since(start)

	// Second call should sleep for ~20ms
	start = time.Now()
	backoff()
	second := time.Since(start)

	// Second should be approximately double the first
	// (with generous tolerance for timing variations)
	if second < first || second > first*3 {
		t.Errorf("backoff not doubling correctly: first=%v, second=%v", first, second)
	}
}

func TestResult_StructureAndTypes(t *testing.T) {
	// Test Result with string type
	t.Run("Result with string", func(t *testing.T) {
		result := Result[string]{
			Result: "test-value",
			Err:    nil,
		}

		if result.Result != "test-value" {
			t.Errorf("Result.Result = %q, want %q", result.Result, "test-value")
		}
		if result.Err != nil {
			t.Errorf("Result.Err = %v, want nil", result.Err)
		}
	})

	// Test Result with error
	t.Run("Result with error", func(t *testing.T) {
		testErr := errors.New("test error")
		result := Result[int]{
			Result: 0,
			Err:    testErr,
		}

		if result.Result != 0 {
			t.Errorf("Result.Result = %d, want 0", result.Result)
		}
		if result.Err != testErr {
			t.Errorf("Result.Err = %v, want %v", result.Err, testErr)
		}
	})

	// Test Result with complex type
	t.Run("Result with slice", func(t *testing.T) {
		result := Result[[]int]{
			Result: []int{1, 2, 3},
			Err:    nil,
		}

		if len(result.Result) != 3 {
			t.Errorf("len(Result.Result) = %d, want 3", len(result.Result))
		}
	})
}

func TestResultMap_StructureAndTypes(t *testing.T) {
	// Test ResultMap with various types
	t.Run("ResultMap with string map", func(t *testing.T) {
		testMap := map[string]int{"a": 1, "b": 2}
		result := ResultMap[string, string, int]{
			Result: "success",
			Err:    nil,
			Map:    testMap,
		}

		if result.Result != "success" {
			t.Errorf("ResultMap.Result = %q, want %q", result.Result, "success")
		}
		if result.Err != nil {
			t.Errorf("ResultMap.Err = %v, want nil", result.Err)
		}
		if len(result.Map) != 2 {
			t.Errorf("len(ResultMap.Map) = %d, want 2", len(result.Map))
		}
		if result.Map["a"] != 1 {
			t.Errorf("ResultMap.Map[\"a\"] = %d, want 1", result.Map["a"])
		}
	})

	// Test ResultMap with error
	t.Run("ResultMap with error", func(t *testing.T) {
		testErr := errors.New("map error")
		result := ResultMap[bool, int, string]{
			Result: false,
			Err:    testErr,
			Map:    nil,
		}

		if result.Result != false {
			t.Errorf("ResultMap.Result = %v, want false", result.Result)
		}
		if result.Err != testErr {
			t.Errorf("ResultMap.Err = %v, want %v", result.Err, testErr)
		}
		if result.Map != nil {
			t.Errorf("ResultMap.Map = %v, want nil", result.Map)
		}
	})
}
