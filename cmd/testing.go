package cmd

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	exit := m.Run()
	os.Exit(exit)
}
