package main

import (
	"fmt"

	//lint:ignore ST1001 Ignoring dot-imports for the example-dir usage
	. "github.com/j-flat/go-veikkaus/goveikkaus"
)

func main() {
	fmt.Println("Displaying the API Glossary")
	fmt.Println()

	client := NewClient(nil)

	glossary := client.Glossary.Get()

	glossary.Print()
}
