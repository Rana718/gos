package cmd

import (
	"fmt"

	"github.com/Rana718/gos/internal/db"

	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm [name]",
	Short: "Remove a saved path",
	Long:  `Remove a saved path by its name.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := db.RemovePath(args[0]); err != nil {
			return err
		}
		fmt.Printf("Removed '%s'.\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
