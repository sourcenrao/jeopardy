package main

import (
	"fmt"
)

type Clue struct {
	Round    uint8
	Value    uint32
	Double   string
	Category string
	Comments string
	Answer   string
	Question string
}

func main() {
	fmt.Println("Clue")
}
