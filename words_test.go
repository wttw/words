package words

import (
	"fmt"
	"testing"
)

func TestErrors(t *testing.T) {
	tests := map[string]struct {
		Val int64
		Min int64
		Max int64
	}{
		"below_range":    {1, 5, 10},
		"above_range":    {15, 5, 10},
		"inverted_range": {5, 10, 1},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := Encode(test.Val, test.Min, test.Max)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestObject(t *testing.T) {
	tests := map[string]struct {
		Range  []int64
		Values []int64
	}{
		"small":          {[]int64{0, 500}, []int64{0, 500, 250}},
		"medium":         {[]int64{0, 10000}, []int64{0, 500, 5000, 10000}},
		"big":            {[]int64{0, 10_000_000}, []int64{0, 5000, 500_000, 5_000_000, 10_000_000}},
		"dynamic_small":  {[]int64{}, []int64{0, 500, 250}},
		"dynamic_medium": {[]int64{}, []int64{0, 500, 5000, 10000}},
		"dynamic_big":    {[]int64{}, []int64{0, 5000, 500_000, 5_000_000, 10_000_000}},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			for _, plain := range test.Values {
				coder, err := New(test.Range...)
				if err != nil {
					t.Errorf("failed to create coder: %s", err)
					continue
				}
				enc, err := coder.Encode(plain)
				if err != nil {
					t.Errorf("unexpected error for %d: %s", plain, err)
					continue
				}
				dec, err := coder.Decode(enc)
				if err != nil {
					t.Errorf("unexpected error decoding %v: %s", enc, err)
					continue
				}
				if plain != dec {
					t.Errorf("failed to round trip %d -> %v -> %d", plain, enc, dec)
				}
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	tests := map[string]struct {
		Min    int64
		Max    int64
		Values []int64
	}{
		"small":  {0, 500, []int64{0, 500, 250}},
		"medium": {0, 10000, []int64{0, 500, 5000, 10000}},
		"big":    {0, 10_000_000, []int64{0, 5000, 500_000, 5_000_000, 10_000_000}},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			for _, plain := range test.Values {
				enc, err := Encode(plain, test.Min, test.Max)
				if err != nil {
					t.Errorf("unexpected error for %d: %s", plain, err)
					continue
				}
				dec, err := Decode(enc, test.Min)
				if err != nil {
					t.Errorf("unexpected error decoding %v: %s", enc, err)
					continue
				}
				if plain != dec {
					t.Errorf("failed to round trip %d -> %v -> %d", plain, enc, dec)
				}
			}
		})
	}
}

func TestDynamicRoundTrip(t *testing.T) {
	tests := map[string]struct {
		Min    int64
		Max    int64
		Values []int64
	}{
		"small":  {0, 500, []int64{0, 500, 250}},
		"medium": {0, 10000, []int64{0, 500, 5000, 10000}},
		"big":    {0, 10_000_000, []int64{0, 5000, 500_000, 5_000_000, 10_000_000}},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			for _, plain := range test.Values {
				enc, err := Encode(plain, test.Min)
				if err != nil {
					t.Errorf("unexpected error for %d: %s", plain, err)
					continue
				}
				dec, err := Decode(enc, test.Min)
				if err != nil {
					t.Errorf("unexpected error decoding %v: %s", enc, err)
					continue
				}
				if plain != dec {
					t.Errorf("failed to round trip %d -> %v -> %d", plain, enc, dec)
				}
				fmt.Printf("%d -> %v\n", plain, enc)
			}
		})
	}
}

func TestExhaustive(t *testing.T) {
	tests := []struct {
		Min int64
		Max int64
	}{
		{1, 10},
		{1, 100},
		{100000, 100100},
		{-400, -300},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("range_%d_%d", test.Min, test.Max), func(t *testing.T) {
			for i := test.Min; i <= test.Max; i++ {
				enc, err := Encode(i, test.Min, test.Max)
				if err != nil {
					t.Errorf("encode %d failed: %s", i, err)
					continue
				}
				fmt.Printf("%d -> %v\n", i, enc)
				dec, err := Decode(enc, test.Min)
				if err != nil {
					t.Errorf("decode %d failed: %s", dec, err)
					continue
				}
				if dec != i {
					t.Errorf("round trip failed %d -> %v -> %d", i, enc, dec)
				}
			}
		})
	}
}

func Example_fixed_range() {
	encoded, _ := Encode(5000, 0, 10_000_000_000)
	decoded, _ := Decode(encoded, 0)
	fmt.Printf("5000 -> %v -> %d", encoded, decoded)
	// Output:
	// 5000 -> [cannon cannon rub fog] -> 5000
}

func Example() {
	encoded, _ := Encode(5000)
	decoded, _ := Decode(encoded)
	fmt.Printf("5000 -> %v -> %d", encoded, decoded)
	// Output:
	// 5000 -> [rub fog] -> 5000
}

func ExampleCoder() {
	coder, _ := New(0, 10_000_000)
	encoded, _ := coder.Encode(42)
	decoded, _ := coder.Decode(encoded)
	fmt.Printf("42 -> %v -> %d", encoded, decoded)
	// Output:
	// 42 -> [cannon cannon tank] -> 42
}
