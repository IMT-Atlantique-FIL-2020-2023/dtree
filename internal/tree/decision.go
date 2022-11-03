package tree

import (
	"fmt"
)

// Required field with some usefull datas:
//   - Config_path: absolute path in the user config data tree
//   - Template_path: path of the k8s template for this config object
//   - Layer: current config level, in which a required value was detected
//   - Label: property name that was found required for this template
type Requirement struct {
	Config_path   string
	Template_path string
	Layer         Tree
	Label         string
}

// References found in the FindReferences function:
// - Name: reference property name
// - Path: path of the referenced template
// - Conf: configuration tree under the prop
type Reference struct {
	Name string
	Path string
	Conf Tree
}

// Iterate over config and template to resolve required field,
// return an array of each found required field as a Requirement.
func Required(config Tree, template Tree) []Requirement {

	// Step 1: Enter the root of the definitions tree
	root := []string{"definitions"}
	definitions := From(template.Get(root))

	// Step 2: Identify the config file entry point
	kind, valid := config.Get([]string{"kind"}).(string)

	if !valid {
		panic(fmt.Sprintf("the user config file is invalid, object kind can't be converted to key:\n%s", config))
	}

	entries := EntryPoints(definitions)            // get the map of entry points
	entry := entries[kind]                         // identify the current one
	definition := definitions.Get([]string{entry}) // find the entry template

	// Step 3: recursively traverse the tree and mark required fields
	definition_tree := From(definition)

	return RecursiveTraversal("", entry, config, definition_tree, definitions)
}

// To identify the entry point from the config file
// we check the 'kind' property. We can reference all the
// possible entry points from the templates by checking
// for each of them if they possess the 'kind' property
// and if the kind enum is equal to the provided one.
func EntryPoints(template Tree) map[string]string {

	dictionnary := map[string]string{}
	query := []string{"properties", "kind", "enum"}

	// for each root queryable object in the templates
	template.For(func(name string, tree Tree) {

		// Check if the object possesses a "properties>kind>enum" prop
		if tree.Has(query) {

			// Reference it's name in the map
			keys := tree.GetArray(query) // with it's kinds as key
			value := name                // and it's name as value

			if len(keys) < 1 {
				panic(fmt.Sprintf("the k8s templates file is invalid, expected at least one value in kind enum '%s'", keys))
			}

			key := fmt.Sprint(keys[0]) // The enum should always have one label which serve to identify the config object

			dictionnary[key] = value
		}
	})

	return dictionnary
}

// Recursively iterate over referenced dependancies in the template
// and find every required field.
func RecursiveTraversal(config_path string, template_path string, config Tree, template Tree, definitions Tree) []Requirement {

	required := []Requirement{}

	required_query := []string{"required"}
	props_query := []string{"properties"}
	items_query := []string{"items"}

	// Check if there is a required field
	if template.Has(required_query) {
		req := template.GetArray(required_query)

		for _, label := range req {
			required = append(required, Requirement{config_path, template_path, config, fmt.Sprint(label)})
		}
	}

	if !template.Has(props_query) {
		panic(fmt.Sprintf("the template file in incorrect, no 'properties' field cannot be found:\n%s", template))
	}

	// Get the properties fields for the template tree
	props_template := template.Get(props_query)
	tree_template := From(props_template)

	// Find all the sub references in fields
	references := FindReferences(config, tree_template)

	for _, reference := range references {

		sub_name := reference.Name
		sub_path := reference.Path
		sub_conf := reference.Conf

		// Find the ressource from the template definitions
		name_query := []string{sub_path}
		if !definitions.Has(name_query) {
			panic(fmt.Sprintf("The reference name was not found in the template definitions '%s'", sub_path))
		}

		definition := From(definitions.Get(name_query))

		// If config is an array, do a recursive traversal for the first element
		// TODO: do we have to check for each sub elem? For max required ?
		// Or even do the summation of all required fields to expand ???

		if sub_conf.Has(items_query) {
			items := sub_conf.GetArray(items_query)

			if len(items) < 1 {
				panic(fmt.Sprintf("the k8s templates file is invalid, expected at least one item in the array '%s'", items))
			}

			// convert the first item to a tree
			sub_conf = From(items[0])
		}

		// Recursive call
		sub_required := RecursiveTraversal(config_path+"."+sub_name, sub_path, sub_conf, definition, definitions)

		required = append(required, sub_required...)
	}

	return required
}

// Recursively explore all fields present in both tree
// and extract the references field found
func FindReferences(conf Tree, temp Tree) []Reference {
	references := []Reference{}

	reference_query := []string{"$ref"}
	items_query := []string{"items"}

	// /!\ According to the k8s specs, replace the arrays
	// in conf with an "items" object to find sub references
	for name, entry := range conf.data {

		// If the prop is an array
		if array, valid := entry.([]interface{}); valid {

			// replace the [...] by { "items": [...] }
			conf.data[name] = map[string]interface{}{"items": array}
		}
	}

	// Iterate over the sub field in the properties (config & template)
	conf.ForBoth(temp, func(name string, sub_conf Tree, sub_temp Tree) {

		to_search := sub_temp

		// If array, check in items
		if to_search.Has(items_query) {
			to_search = From(to_search.Get(items_query))
		}

		// If there is a reference, find the template and keep iterating in it
		if to_search.Has(reference_query) {

			ref := fmt.Sprint(to_search.Get(reference_query)) // $ref are given as '#/definitions/...'
			path := ref[14:]                                  // extract the ressource name

			references = append(references, Reference{name, path, sub_conf})
		}
	})

	return references
}
