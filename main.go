package main

import (
	"os"

	"github.com/JamesChung/cprl/cmd"
	"github.com/JamesChung/cprl/internal/config"
	"github.com/JamesChung/cprl/pkg/util"
)

func init() {
	err := config.Read()
	util.ExitOnErr(err)
}

func main() {
	cli := cmd.NewCmd()

	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
