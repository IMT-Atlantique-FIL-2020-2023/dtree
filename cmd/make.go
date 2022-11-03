package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/IMT-Atlantique-FIL-2020-2023/dtree/internal/tree"
	"github.com/spf13/cobra"
)

var config string   // config file path
var template string // template file path
var format string   // output format

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

		switch format {
		case "config":
			output_config_list(required)
		case "template":
			output_template_list(required)
		case "tree":
			output_config_tree(required)
		default:
			panic("Wrong format given for the format parameter. Accepted values are: 'config', 'template' or 'tree'")
		}
	},
}

func init() {
	make.Flags().StringVarP(&config, "config", "c", "./config.yaml", "The yaml config file")
	make.Flags().StringVarP(&template, "template", "t", "./template.json", "The template definition file")
	make.Flags().StringVarP(&format, "format", "f", "tree", "The format of the output (config, template or tree)")

	make.MarkFlagFilename("config")
	make.MarkFlagFilename("template")

	make.MarkFlagRequired("config")
	make.MarkFlagRequired("template")

	root.AddCommand(make)
}

func output_template_list(required []tree.Requirement) {
	if len(required) < 1 {
		fmt.Println(" * None")
	}

	for _, req := range required {
		fmt.Println(" * " + req.Template_path + ".properties." + req.Label)
	}
}

func output_config_list(required []tree.Requirement) {
	if len(required) < 1 {
		fmt.Println(" * None")
	}

	for _, req := range required {
		fmt.Println(" * " + req.Config_path + "." + req.Label + " = " + fmt.Sprint(req.Layer.Get([]string{req.Label})))
	}
}

func output_config_tree(required []tree.Requirement) {
	if len(required) < 1 {
		fmt.Println("Empty")
	}

	generated := tree.New(tree.Table{})

	for _, req := range required {

		// TODO: insert not only the required field in the layer,
		// but the whole layer and it's sub levels (limited depth?)
		// Problem: we have to pass an object param to the Set function
		// otherwise the copy will be shallow and the object will
		// disapear (must check how to deep copy an interface...)
		query := strings.Split(req.Config_path+"."+req.Label, ".")[1:]
		generated.Set(query, fmt.Sprint(req.Layer.Get([]string{req.Label})))

		// generated.Set(query, req.Layer.Get(query[len(query)-1:]))
	}

	fmt.Println(generated)
}
