package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/Rana718/gos/internal/config"
	"github.com/Rana718/gos/internal/db"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "gos",
	Short:   "gos - a fast CLI to manage paths, aliases, and more",
	Long:    `gos is a CLI tool to manage paths, aliases, and more.`,
	Version: "0.2.0",
}

func Execute() {
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "@") {
		aliasName := strings.TrimPrefix(os.Args[1], "@")
		handleAlias(aliasName)
		return
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func handleAlias(name string) {
	if name == "" {
		fmt.Fprintln(os.Stderr, "Usage: gos @<alias_name>")
		os.Exit(1)
	}

	if name == "config" {
		config.OpenInEditor()
		return
	}

	command, err := db.GetAlias(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unknown alias '@%s'. Use 'gos alias ls' to see aliases.\n", name)
		os.Exit(1)
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "Error running '@%s': %v\n", name, err)
		os.Exit(1)
	}
}
