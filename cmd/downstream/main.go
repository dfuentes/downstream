package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func getWorkingDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("Failed to get working directory: %s", err)
	}
	return dir
}

func main() {
	var rootCmd = &cobra.Command{}
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.Execute()
}
