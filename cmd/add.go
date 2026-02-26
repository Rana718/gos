package cmd

import (
	"fmt"

	"github.com/Rana718/gos/internal/db"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [name] [path]",
	Short: "Save a path with a name",
	Long:  `Save a directory path with a custom name for quick access later.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		path := args[1]

		if err := db.AddPath(name, path); err != nil {
			return fmt.Errorf("failed to save path: %w", err)
		}
		fmt.Printf("Saved '%s' -> %s\n", name, path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
