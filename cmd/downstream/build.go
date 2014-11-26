package main

import (
	"log"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build [module]",
	Short: "Build downstream node modules",
	Long:  "Build downstream node modules, if no module is specified, builds all downstreams.",
	Run:   doBuild,
}

func doBuild(cmd *cobra.Command, args []string) {
	_, downstreamDependencies := getDownstreams()
	workingDirectory := getWorkingDirectory()

	for _, ds := range downstreamDependencies {
		if len(args) > 0 && !contains(args, ds.Name) {
			continue
		}
		if err := ds.Build(workingDirectory); err != nil {
			log.Fatalf("Failed to build %s: %s", ds.Name, err)
		}
	}
}

func contains(arr []string, val string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}
