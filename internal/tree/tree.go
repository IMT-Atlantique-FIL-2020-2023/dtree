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

// Query the tree data and return the result object
func (t *Tree) Set(query []string, value interface{}) {
	if len(query) > 1 {
		current := t.data[query[0]]

		if _, ok := current.(Table); !ok {
			current = Table{} // create if path does not exist
		}

		stree := From(current)
		stree.Set(query[1:], value)

		t.data[query[0]] = stree.data

	} else {
		t.data[query[0]] = value
	}
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
			tree := New(table)

			// Call the closure
			fn(name, tree)
		}
	}
}

// Parallel iteration over both trees based on this tree props
func (t Tree) ForBoth(other Tree, fn func(name string, a_tree Tree, b_tree Tree)) {

	t.For(func(a_name string, a_tree Tree) {
		other_query := []string{a_name}

		// Check for prop in other tree
		if other.Has(other_query) {

			// If the prop can also be queried on
			if table, valid := other.Get(other_query).(Table); valid {

				// Convert it to a tree
				b_tree := New(table)

				// Call the closure
				fn(a_name, a_tree, b_tree)
			}
		}
	})
}

func New(data Table) Tree {
	t := Tree{data}
	return t
}

func FromJSON(file string) Tree {
	var data Table

	json.Unmarshal([]byte(file), &data)

	tree := New(data)

	return tree
}

func FromYAML(file string) Tree {
	var data Table

	yaml.Unmarshal([]byte(file), &data)

	tree := New(data)

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

	return New(obj_table)
}
