package utils

import (
	"fmt"

	"github.com/fatih/color"
	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
)

func SetSuccessStep(msg AppType.AppStep) {
	gb := color.New(color.FgGreen, color.Bold)
	fmt.Printf("%s... ", msg)
	gb.Printf("ok\n")
}

func SetAttentionStep(msg AppType.AppStep) {
	yb := color.New(color.FgYellow, color.Bold)
	fmt.Printf("%s...", msg)
	yb.Printf("attention\n")
}

func SetErrorStep(msg AppType.AppStep) {
	rb := color.New(color.FgRed, color.Bold)
	fmt.Printf("%s...", msg)
	rb.Printf("error\n")
}
