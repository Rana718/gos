package cmd

import (
	"fmt"
	"os/exec"

	"github.com/Rana718/gos/internal/db"
	"github.com/Rana718/gos/internal/resolver"

	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open [name|path|.]",
	Short: "Open in Explorer",
	Long:  `Open a saved path, direct path, or current directory in the file explorer.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		loadedPaths := db.LoadPathsMap()
		finalPath, err := resolver.ResolvePath(args[0], loadedPaths)
		if err != nil {
			return err
		}

		err = exec.Command("explorer", finalPath).Run()
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
				fmt.Println("Opened in Explorer:", finalPath)
				return nil
			}
			return fmt.Errorf("failed to open Explorer: %w", err)
		}
		fmt.Println("Opened in Explorer:", finalPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
}
