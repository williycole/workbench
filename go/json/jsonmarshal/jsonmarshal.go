package jsonmarshal

import (
	"encoding/json"
	"log"
)

// "encoding/json"

func main() {}

/*
Complete the marshalAll function. It accepts a slice of "items", which can be of any type.
The expectation is that they are structs of various forms.

It should return a slice of slices of bytes (I didn't stutter), where each resulting slice of bytes represents
the JSON representation of an item in the input slice. If an item cannot be marshaled, the function should
immediately return an error.

Return the marshalled data in the same order as the input items.
*/
type All struct {
	T any `json:"T,omitempty"`
}

func marshalAll[T any](items []T) ([][]byte, error) {
	var all = [][]byte{}
	for _, i := range items {
		m, err := json.Marshal(i)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		all = append(all, m)

	}

	return all, nil
}
