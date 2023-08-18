# words

[![Go Reference](https://pkg.go.dev/badge/github.com/wttw/words.svg)](https://pkg.go.dev/github.com/wttw/words)

Encode an integer into a short list of memorable words and decode it back again.

```go
	encoded, _ := Encode(5000)
	decoded, _ := Decode(encoded)
	fmt.Printf("5000 -> %v -> %d", encoded, decoded)
	// Output:
	// 5000 -> [rub fog] -> 5000
```

For the simple case the list of words will be of variable length - single words for small integers, longer
sentences as they get larger. For some use cases having a fixed length of sentence for all valid input
integers can be convenient. Including a range in the parameters to Encode will do that.

```go
	encoded, _ := Encode(5000, 0, 10_000_000_000)
	decoded, _ := Decode(encoded, 0)
	fmt.Printf("5000 -> %v -> %d", encoded, decoded)
	// Output:
	// 5000 -> [cannon cannon rub fog] -> 5000
```

It's slightly more convenient, and slightly more efficient, to create and reuse a Coder to do the same thing.

```go
	coder, _ := New(0, 10_000_000)
	encoded, _ := coder.Encode(42)
	decoded, _ := coder.Decode(encoded)
	fmt.Printf("42 -> %v -> %d", encoded, decoded)
	// Output:
	// 42 -> [cannon cannon tank] -> 42
```

words uses the word list from https://github.com/wcodes-org/wordlist