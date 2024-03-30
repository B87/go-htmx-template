package main

import (
	"fmt"

	"github.com/B87/go-htmx-template/cmd"
)

func main() {
	err := cmd.RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}
