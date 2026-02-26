package cmd

import (
	"fmt"
	"strings"

	"github.com/Rana718/gos/internal/db"

	"github.com/spf13/cobra"
)

var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Manage aliases",
	Long:  `Create, remove, and list command aliases. Run with gos @name.`,
}

var aliasSetCmd = &cobra.Command{
	Use:   "set [name] [command...]",
	Short: "Set an alias",
	Long:  `Set an alias that maps a name to a command. Usage: gos alias set build "go build ./..."`,
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		command := strings.Join(args[1:], " ")

		if err := db.SetAlias(name, command); err != nil {
			return fmt.Errorf("failed to set alias: %w", err)
		}
		fmt.Printf("Alias '@%s' -> %s\n", name, command)
		return nil
	},
}

var aliasRmCmd = &cobra.Command{
	Use:   "rm [name]",
	Short: "Remove an alias",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := db.RemoveAlias(args[0]); err != nil {
			return err
		}
		fmt.Printf("Alias '@%s' removed.\n", args[0])
		return nil
	},
}

var aliasLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all aliases",
	RunE: func(cmd *cobra.Command, args []string) error {
		aliases, err := db.ListAliases()
		if err != nil {
			return err
		}
		if len(aliases) == 0 {
			fmt.Println("No aliases.")
			return nil
		}
		fmt.Println("Aliases:")
		for _, a := range aliases {
			fmt.Printf("  @%s -> %s\n", a.AliasName, a.Command)
		}
		return nil
	},
}

func init() {
	aliasCmd.AddCommand(aliasSetCmd)
	aliasCmd.AddCommand(aliasRmCmd)
	aliasCmd.AddCommand(aliasLsCmd)
	rootCmd.AddCommand(aliasCmd)
}
