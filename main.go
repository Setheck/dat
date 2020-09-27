package main

import (
	"os"
)

func main() {
	rootCmd := NewRootCommand()
	rootCmd.ParseFlags()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
