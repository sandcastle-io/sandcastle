package main

import (
	"os"

	"github.com/mbobrovskyi/kube-sandcastle/internal/cmd"
)

func main() {
	if err := cmd.NewDefaultSandcastlectlCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
