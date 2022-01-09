package main

import (
	"github.com/gefion-tech/tg-exchanger-server/cmd"
	"github.com/spf13/cobra"
)

func main() {
	cobra.CheckErr(cmd.NewRootCmd().Execute())
}
