package main

import (
	"os"

	"github.com/JamesChung/cprl/cmd"
	"github.com/JamesChung/cprl/internal/config"
	"github.com/JamesChung/cprl/pkg/util"
	"github.com/pterm/pterm"
)

func init() {
	err := config.Read()
	if err != nil {
		pterm.Error.Println(err)
		err = config.Create()
	}
	util.ExitOnErr(err)
}

func main() {
	cli := cmd.NewCmd()

	// Generate cli documents
	// err := doc.GenMarkdownTree(cli, "./docs")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
