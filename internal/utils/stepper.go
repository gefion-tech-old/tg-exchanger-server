package utils

import (
	"fmt"

	"github.com/fatih/color"
)

func SetSuccessStep(msg string) {
	gb := color.New(color.FgGreen, color.Bold)
	fmt.Printf(fmt.Sprintf("%s... ", msg))
	gb.Printf("ok\n")
}

func SetAttentionStep(msg string) {
	yb := color.New(color.FgYellow, color.Bold)
	fmt.Printf(fmt.Sprintf("%s...", msg))
	yb.Printf("attention\n")
}

func SetErrorStep(msg string) {
	rb := color.New(color.FgRed, color.Bold)
	fmt.Printf(fmt.Sprintf("%s...", msg))
	rb.Printf("error\n")
}
