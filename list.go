package main

import (
	"github.com/spf13/cobra"
)

var helper *Helper

var rootCmd = &cobra.Command{
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		helper, err = NewHelper()

		return err
	},
}

var listCmd = &cobra.Command{
	Use: "list",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		helper.List(ctx)
		return nil
	},
}

var listEngineCmd = &cobra.Command{
	Use:  "list-engine",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		helper.ListEngine(ctx, args[0])
		return nil
	},
}
