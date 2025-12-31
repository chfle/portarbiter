package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: portarbiter <port>")
		os.Exit(1)
	}

	port := os.Args[1]
	fmt.Printf("portarbiter: inspecting port %s\n", port)
}

