package main

import (
	"fmt"
	"os"

	"github.com/construct-go/cli/construct"
)

func main() {
	if err := construct.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}