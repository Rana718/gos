package cmd

import (
	"fmt"

	"github.com/Rana718/gos/internal/db"

	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls [name]",
	Short: "List saved paths",
	Long:  `List all saved paths, or show the path for a specific name.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			paths, err := db.ListPaths()
			if err != nil {
				return fmt.Errorf("failed to list paths: %w", err)
			}
			if len(paths) == 0 {
				fmt.Println("No saved paths.")
				return nil
			}
			fmt.Println("Saved paths:")
			for _, p := range paths {
				fmt.Printf("  %s: %s\n", p.Name, p.Path)
			}
			return nil
		}

		path, err := db.GetPath(args[0])
		if err != nil {
			return err
		}
		fmt.Printf("%s: %s\n", args[0], path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
