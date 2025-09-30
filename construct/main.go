package construct

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "1.0.0"

var rootCmd = &cobra.Command{
	Use:   "construct",
	Short: "Construct - Full-stack Vue + Go framework",
	Long: `Construct CLI - A modern full-stack framework combining Vue 3 and Base Go.

Build powerful web applications with the best of both worlds:
  • Vue 3 for reactive, component-based UI
  • Base Go for fast, type-safe backend

One framework. One command. One binary.`,
	Version: version,
}

func init() {
	rootCmd.AddCommand(devCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(generateCmd)

	// Future commands:
	// rootCmd.AddCommand(migrateCmd)
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func printBanner() {
	banner := `
   ____                _                   _
  / ___|___  _ __  ___| |_ _ __ _   _  ___| |_
 | |   / _ \| '_ \/ __| __| '__| | | |/ __| __|
 | |__| (_) | | | \__ \ |_| |  | |_| | (__| |_
  \____\___/|_| |_|___/\__|_|   \__,_|\___|\__|

  Full-stack Vue + Go Framework v%s
`
	fmt.Printf(banner, version)
}