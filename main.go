package main

import (
	"fmt"
	"os"

	"github.com/JamesChung/cprl/cmd"
	"github.com/JamesChung/cprl/internal/config"
)

func init() {
	err := config.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

func main() {
	cli := cmd.NewCmd()

	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
