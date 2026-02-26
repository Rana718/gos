package cmd

import (
	"fmt"
	"os/exec"

	"github.com/Rana718/gos/internal/db"
	"github.com/Rana718/gos/internal/resolver"

	"github.com/spf13/cobra"
)

var codeCmd = &cobra.Command{
	Use:   "code [name|path|.]",
	Short: "Open in VS Code",
	Long:  `Open a saved path, direct path, or current directory in Visual Studio Code.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		loadedPaths := db.LoadPathsMap()
		finalPath, err := resolver.ResolvePath(args[0], loadedPaths)
		if err != nil {
			return err
		}

		err = exec.Command("code", finalPath).Run()
		if err != nil {
			return fmt.Errorf("failed to open VS Code: %w", err)
		}
		fmt.Println("Opened in VS Code:", finalPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(codeCmd)
}
