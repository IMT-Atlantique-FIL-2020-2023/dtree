package cmd

import (
	"fmt"
	"os"

	"github.com/IMT-Atlantique-FIL-2020-2023/dtree/internal/tree"
	"github.com/spf13/cobra"
)

var config string   // config file path
var template string // template file path

func ReadFile(path string) string {

	var file, err = os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	return string(file)
}

var make = &cobra.Command{
	Use:   "make",
	Short: "Generate the dependancy tree for provided files",
	Long:  "Generate the dependancy tree for the given .yaml config and .json struct",
	Run: func(cmd *cobra.Command, args []string) {

		var file_config = ReadFile(config)
		var file_template = ReadFile(template)

		var tree_config = tree.FromYAML(file_config)
		var tree_template = tree.FromJSON(file_template)

		required := tree.Required(tree_config, tree_template)

		fmt.Println("\nRequired fields for the given configuration: ")

		if len(required) < 1 {
			fmt.Println(" * None")
		}

		for _, v := range required {
			fmt.Println(" * " + v)
		}
	},
}

func init() {
	make.Flags().StringVarP(&config, "config", "c", "./config.yaml", "The yaml config file")
	make.Flags().StringVarP(&template, "template", "t", "./template.json", "The template definition file")

	make.MarkFlagFilename("config")
	make.MarkFlagFilename("template")

	make.MarkFlagRequired("config")
	make.MarkFlagRequired("template")

	root.AddCommand(make)
}
