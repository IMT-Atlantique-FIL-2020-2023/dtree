package tree

import (
	"fmt"

	"encoding/json"

	"gopkg.in/yaml.v3"
)

type Table = map[string]interface{}

type Tree struct {
	data Table
}

func (t Tree) Data() Table {
	return t.data
}

// Check if a query path is valid in the tree data
func (t Tree) Has(query []string) bool {
	var current interface{} = t.data

	for _, key := range query {
		ctable, ok := current.(Table)

		if !ok {
			return false
		}

		current = ctable[key]
	}

	return current != nil
}

// Query the tree data and return the result object
func (t Tree) Get(query []string) interface{} {
	var current interface{} = t.data

	for _, key := range query {
		ctable, ok := current.(Table)

		if !ok {
			panic(fmt.Sprintf("tried to search in a non queryable object type '%s'", current))
		}

		current = ctable[key]
	}

	return current
}

// Query the tree data for an array
func (t Tree) GetArray(query []string) []interface{} {

	result := t.Get(query) // with it's kinds as key

	keys, valid := result.([]interface{})

	if !valid {
		panic(fmt.Sprintf("the k8s templates file is invalid, object kind enum was wrongly formatted '%s'", result))
	}

	return keys
}

// Iterate over the tree props and call a function for each sub tree
// It will ignore simple props (flat fields)
func (t Tree) For(fn func(name string, tree Tree)) {

	// for each root object in the templates
	for name, entry := range t.data {

		// If the object can be queried on
		if table, valid := entry.(Table); valid {

			// Convert it to a tree
			tree := new(table)

			// Call the closure
			fn(name, tree)
		}
	}
}

// Parralel iteration over both trees based on this tree props
func (t Tree) ForBoth(other Tree, fn func(name string, a_tree Tree, b_tree Tree)) {

	t.For(func(a_name string, a_tree Tree) {
		other.For(func(b_name string, b_tree Tree) {

			// TODO: if better performances are needed, we can only iterate
			// over this tree (a) and search for props if they are present
			// in the the other tree (b)
			if a_name == b_name {
				fn(a_name, a_tree, b_tree)
			}
		})
	})
}

func new(data Table) Tree {
	t := Tree{data}
	return t
}

func FromJSON(file string) Tree {
	var data Table

	json.Unmarshal([]byte(file), &data)

	tree := new(data)

	return tree
}

func FromYAML(file string) Tree {
	var data Table

	yaml.Unmarshal([]byte(file), &data)

	tree := new(data)

	return tree
}

// Return a string formatted as a json object
func (t Tree) String() string {
	formatted, err := json.MarshalIndent(t.data, "", "  ")

	if err != nil {
		panic(fmt.Sprintf("An error occured while trying to pretty print the tree: '%s'", err))
	}

	return string(formatted)
}

// Try to instanciate a Tree from an interface
func From(obj interface{}) Tree {
	obj_table, valid := obj.(Table)

	if !valid {
		panic(fmt.Sprintf("The given interface could not be resolved into a tree:\n%v", obj))
	}

	return new(obj_table)
}
