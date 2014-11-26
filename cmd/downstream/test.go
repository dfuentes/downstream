package main

import (
	"log"

	"github.com/spf13/cobra"
)

var verbose bool

var testCmd = &cobra.Command{
	Use:   "test [module]",
	Short: "Test downstream node module",
	Long:  "Test downstream node module, if no module is specified, will test all downstreams.",
	Run:   doTest,
}

func init() {
	testCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "show all test output")
}

func doTest(cmd *cobra.Command, args []string) {
	_, downstreamDependencies := getDownstreams()
	workingDirectory := getWorkingDirectory()

	for _, ds := range downstreamDependencies {
		if len(args) > 0 && !contains(args, ds.Name) {
			continue
		}
		if err := ds.Test(workingDirectory, verbose); err != nil {
			log.Fatal(err)
		}
	}
}
