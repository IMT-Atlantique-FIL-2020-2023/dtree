package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "dtree",
	Short: "...",
	Long:  "...",
}

func Execute() {
	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "An error occured '%s'", err)
		os.Exit(1)
	}
}
