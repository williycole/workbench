package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(Reverse("hello world"))
}

// Reverse reverses a string left to right
// Notice that we need to capitalize the first letter of the function
// If we don't then we won't be able to access this function outside of the
// mystrings package
func Reverse(s string) string {
	result := ""
	for _, v := range s {
		result = string(v) + result
	}
	return result
}

// easy, got done in a matter of seconds
func tolower() {
	s := "Hello, World!!!"
	s = strings.ToLower(s)
	fmt.Println(s)
}
