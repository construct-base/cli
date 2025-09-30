package main

import (
	"fmt"
	"os"

	"github.com/construct-go/cli/construct"
)

func main() {
	if err := construct.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}