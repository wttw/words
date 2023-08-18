package words

import (
	"fmt"
	"slices"
	"sync"
)

// Generate our list of words from the list from https://github.com/wcodes-org/wordlist
//go:generate go run generate-list.go Wordlist.tsv wordlist.go

const maxInt64 = 1<<63 - 1

type Coder struct {
	wordCount int
	minVal    int64
	maxVal    int64
}

var reverse map[string]int
var reverseOnce sync.Once

// New creates a new Coder. If passed two values, minimum and maximum, it will create slices of strings
// of a fixed length, adequate to encode all integers in that range. If passed a single value, minimum, it will
// create slices of strings of variable length, able to encode any number greater than or equal to minimum.
// If called with no parameters it creates a Coder that encodes any non-negative integer into a variable
// length slice of strings.
func New(valueRange ...int64) (Coder, error) {
	switch len(valueRange) {
	case 0:
		return Coder{
			wordCount: -1,
			minVal:    0,
			maxVal:    maxInt64,
		}, nil
	case 1:
		return Coder{
			wordCount: -1,
			minVal:    valueRange[0],
			maxVal:    maxInt64,
		}, nil
	case 2:
		minVal, maxVal := valueRange[0], valueRange[1]
		if minVal >= maxVal {
			return Coder{}, fmt.Errorf("invalid range (%d, %d)", minVal, maxVal)
		}
		span := maxVal - minVal
		wordCount := 1
		for span > listSize {
			wordCount++
			span /= listSize
		}
		return Coder{
			wordCount: wordCount,
			minVal:    minVal,
			maxVal:    maxVal,
		}, nil
	default:
		return Coder{}, fmt.Errorf("wrong number of parameters passed to words.New: %d", len(valueRange))
	}
}

// Length returns the fixed length of slices this Coder generates, or -1 if it generates variable lengths
func (c Coder) Length() int {
	return c.wordCount
}

// Encode converts an integer into a list of words
func (c Coder) Encode(value int64) ([]string, error) {
	if c.wordCount < 0 {
		return c.encodeDynamic(value)
	}
	if value < c.minVal || value > c.maxVal {
		return nil, fmt.Errorf("%d is outside the range (%d, %d)", value, c.minVal, c.maxVal)
	}
	result := make([]string, c.wordCount)
	value -= c.minVal
	wordCount := c.wordCount
	for wordCount > 0 {
		wordCount--
		result[wordCount] = List[value%listSize]
		value /= listSize
	}
	return result, nil
}

func (c Coder) encodeDynamic(value int64) ([]string, error) {
	if value < c.minVal {
		return nil, fmt.Errorf("%d is less than minimum %d", value, c.minVal)
	}
	value -= c.minVal
	result := make([]string, 0)
	for {
		result = append(result, List[value%listSize])
		value /= listSize
		if value == 0 {
			slices.Reverse(result)
			return result, nil
		}
	}
}

// Decode converts a list of words back into an integer
func (c Coder) Decode(words []string) (int64, error) {
	reverseOnce.Do(func() {
		reverse = make(map[string]int, listSize)
		for i, v := range List {
			reverse[v] = i
		}
	})
	if c.wordCount != -1 && len(words) != c.wordCount {
		return 0, fmt.Errorf("expected %d words, not %d", c.wordCount, len(words))
	}
	var result int64 = 0
	for _, v := range words {
		num, ok := reverse[v]
		if !ok {
			return 0, fmt.Errorf("unexpected word '%s' found", v)
		}
		result = (result * listSize) + int64(num)
	}
	return result + c.minVal, nil
}

// Length returns the number of words required to encode integers in the given range
func Length(minVal, maxVal int64) (int, error) {
	c, err := New(minVal, maxVal)
	if err != nil {
		return 0, err
	}
	return c.Length(), nil
}

// Encode a value in the range (minVal..maxVal) into a list of words
func Encode(value int64, ranges ...int64) ([]string, error) {
	c, err := New(ranges...)
	if err != nil {
		return nil, err
	}
	return c.Encode(value)
}

// Decode a list of words back into an integer
func Decode(words []string, ranges ...int64) (int64, error) {
	c, err := New(ranges...)
	if err != nil {
		return 0, err
	}
	return c.Decode(words)
}
