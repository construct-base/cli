package construct

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start production server",
	Long:  "Start the production server (requires running 'construct build' first)",
	Run: func(cmd *cobra.Command, args []string) {
		runStart()
	},
}

func runStart() {
	printBanner()
	fmt.Println("🚀 Starting production server...")
	fmt.Println()

	root, err := findProjectRoot()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}

	binaryName := "construct"
	if runtime.GOOS == "windows" {
		binaryName = "construct.exe"
	}

	binaryPath := filepath.Join(root, binaryName)
	if !fileExists(binaryPath) {
		fmt.Printf("❌ Binary not found: %s\n", binaryPath)
		fmt.Println("   Run 'construct build' first")
		os.Exit(1)
	}

	cmd := exec.Command(binaryPath)
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Failed to start server: %v\n", err)
		os.Exit(1)
	}
}