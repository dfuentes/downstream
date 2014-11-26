package main

import (
	"fmt"
	"log"
	"path"

	"github.com/dfuentes/downstream/downstream"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List downstream node modules.",
	Long:  "List downstream node modules, use the -p flag to show only prod dependencies.",
	Run:   doList,
}

var hideDev bool

func init() {
	listCmd.Flags().BoolVarP(&hideDev, "prod", "p", false, "Show only prod dependencies.")
}

func doList(cmd *cobra.Command, args []string) {

	pack, downstreamDependencies := getDownstreams()

	for _, ds := range downstreamDependencies {
		if ds.DevDep && hideDev {
			continue
		}
		fmt.Printf("%s@%s depends on %s of %s\n", ds.Name, ds.Version, ds.DependsOn, pack.Name)
	}
}

func getDownstreams() (downstream.Package, []downstream.Downstream) {
	workingDirectory := getWorkingDirectory()

	if !downstream.IsNodeDir(workingDirectory) {
		log.Fatal("downstream must be run from a node module directory.")
	}

	pack, err := downstream.LoadPackage(path.Join(workingDirectory, "package.json"))
	if err != nil {
		log.Fatalf("Could not load package.json: %s", err)
	}
	codeDir := path.Join(workingDirectory, "..")
	downstreamDependencies, err := downstream.List(codeDir, pack.Name)
	if err != nil {
		log.Fatalf("Unable to list downstream dependencies: %s", err)
	}
	return pack, downstreamDependencies
}
