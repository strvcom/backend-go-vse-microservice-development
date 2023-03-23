package main

import (
	"fmt"
)

func main() {
	data1 := map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}

	data2 := map[string]string{
		"EUR": "€",
		"USD": "$",
	}

	fmt.Println(MapKeys(data1))
	fmt.Println(MapKeys(data2))
}

func MapKeys[K comparable, V any](input map[K]V) []K {
	result := make([]K, 0, len(input))
	for k, _ := range input {
		result = append(result, k)
	}
	return result
}
