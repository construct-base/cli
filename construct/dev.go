package construct

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/base-go/mamba"
)

var devCmd = &mamba.Command{
	Use:   "dev",
	Short: "Start development servers",
	Long:  "Start both Go API (port 8100) and Vue dev server (port 3100) with hot reload",
	Run: func(cmd *mamba.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		runDev(verbose)
	},
}

func init() {
	devCmd.Flags().BoolP("verbose", "v", false, "show all logs (verbose mode)")
}

func runDev(verbose bool) {
	printBanner()
	ShowProgress("Starting development servers...")
	fmt.Println()

	root, err := findProjectRoot()
	if err != nil {
		ShowError(fmt.Sprintf("Error: %v", err))
		os.Exit(1)
	}

	// Kill any process using ports
	if !verbose {
		WithSpinner("Cleaning up ports 8100 and 3100", func() error {
			killCmd := exec.Command("sh", "-c", "lsof -ti:8100 | xargs kill -9 2>/dev/null || true")
			killCmd.Run()
			killCmd2 := exec.Command("sh", "-c", "lsof -ti:3100 | xargs kill -9 2>/dev/null || true")
			killCmd2.Run()
			return nil
		})
	} else {
		ShowProgress("Checking for processes on port 8100...")
		killCmd := exec.Command("sh", "-c", "lsof -ti:8100 | xargs kill -9 2>/dev/null || true")
		killCmd.Run()
		ShowProgress("Checking for processes on port 3100...")
		killCmd2 := exec.Command("sh", "-c", "lsof -ti:3100 | xargs kill -9 2>/dev/null || true")
		killCmd2.Run()
	}
	fmt.Println()

	// Generate Swagger documentation
	if !verbose {
		WithSpinner("Generating API documentation", func() error {
			swagCmd := exec.Command("swag", "init", "--dir", "./", "--output", "./docs", "--parseDependency", "--parseInternal", "--parseVendor", "--parseDepth", "1", "--generatedTime", "false")
			swagCmd.Dir = root
			swagCmd.Stdout = nil
			swagCmd.Stderr = nil
			// Ignore errors - not critical if swagger isn't installed
			swagCmd.Run()
			return nil
		})
	} else {
		ShowProgress("Generating API documentation...")
		swagCmd := exec.Command("swag", "init", "--dir", "./", "--output", "./docs", "--parseDependency", "--parseInternal", "--parseVendor", "--parseDepth", "1", "--generatedTime", "false")
		swagCmd.Dir = root
		swagCmd.Stdout = NewPrefixWriter(os.Stdout, "   ")
		swagCmd.Stderr = NewPrefixWriter(os.Stderr, "   ")
		if err := swagCmd.Run(); err != nil {
			ShowProgress("Warning: Failed to generate docs (swag might not be installed)")
		} else {
			ShowSuccess("API documentation generated")
		}
	}
	fmt.Println()

	// Run go mod tidy to sync dependencies
	if !verbose {
		WithSpinner("Syncing Go dependencies", func() error {
			tidyCmd := exec.Command("go", "mod", "tidy")
			tidyCmd.Dir = root
			tidyCmd.Stdout = nil
			tidyCmd.Stderr = nil
			return tidyCmd.Run()
		})
	} else {
		ShowProgress("Syncing Go dependencies...")
		tidyCmd := exec.Command("go", "mod", "tidy")
		tidyCmd.Dir = root
		tidyCmd.Stdout = NewPrefixWriter(os.Stdout, "   ")
		tidyCmd.Stderr = NewPrefixWriter(os.Stderr, "   ")
		if err := tidyCmd.Run(); err != nil {
			ShowProgress("Warning: Failed to sync dependencies")
		} else {
			ShowSuccess("Go dependencies synced")
		}
	}
	fmt.Println()

	// Install Vue dependencies first (before starting servers)
	vueDir := filepath.Join(root, "vue")
	nodeModulesPath := filepath.Join(vueDir, "node_modules")
	if _, err := os.Stat(nodeModulesPath); os.IsNotExist(err) {
		installCmd := detectPackageManager(vueDir, "install")
		installCmd.Dir = vueDir

		if verbose {
			// Verbose: show all output
			ShowProgress("Installing Vue dependencies...")
			installCmd.Stdout = NewPrefixWriter(os.Stdout, "   ")
			installCmd.Stderr = NewPrefixWriter(os.Stderr, "   ")
			if err := installCmd.Run(); err != nil {
				ShowError(fmt.Sprintf("Failed to install dependencies: %v", err))
				os.Exit(1)
			}
			ShowSuccess("Dependencies installed")
		} else {
			// Normal: use spinner
			err := WithSpinner("Installing Vue dependencies", func() error {
				installCmd.Stdout = nil
				installCmd.Stderr = nil
				return installCmd.Run()
			})
			if err != nil {
				ShowError(fmt.Sprintf("Failed to install dependencies: %v", err))
				os.Exit(1)
			}
		}
		fmt.Println()
	}

	// Start Go API
	goCmd := exec.Command("go", "run", "main.go")
	goCmd.Dir = root

	// Configure output based on verbose flag
	if verbose {
		// Verbose mode: show all logs with prefix
		goCmd.Stdout = NewPrefixWriter(os.Stdout, "🔷 [Go]   ")
		goCmd.Stderr = NewPrefixWriter(os.Stderr, "🔷 [Go]   ")
	} else {
		// Normal mode: filter logs - extract key info only
		goCmd.Stdout = NewFilteredWriter(os.Stdout, "🔷 [Go]   ", []string{
			"Server starting",
			"ERROR",
			"FATAL",
			"panic",
		})
		goCmd.Stderr = NewFilteredWriter(os.Stderr, "🔷 [Go]   ", []string{
			"ERROR",
			"FATAL",
			"panic",
		})
	}
	goCmd.Env = os.Environ()

	if err := goCmd.Start(); err != nil {
		ShowError("Failed to start Go server")
		os.Exit(1)
	}

	ShowSuccess("Go API starting on http://localhost:8100")

	// Start Vue dev server
	vueCmd := detectPackageManager(vueDir, "dev")
	vueCmd.Dir = vueDir

	if verbose {
		// Verbose mode: show all logs
		vueCmd.Stdout = NewPrefixWriter(os.Stdout, "🟢 [Vue]  ")
		vueCmd.Stderr = NewPrefixWriter(os.Stderr, "🟢 [Vue]  ")
	} else {
		// Normal mode: filter logs
		vueCmd.Stdout = NewFilteredWriter(os.Stdout, "🟢 [Vue]  ", []string{
			"VITE",
			"ready in",
			"Local:",
			"Network:",
			"ERROR",
			"warning",
			"✓",
		})
		vueCmd.Stderr = NewFilteredWriter(os.Stderr, "🟢 [Vue]  ", []string{
			"ERROR",
			"warning",
			"error",
		})
	}
	vueCmd.Env = os.Environ()

	if err := vueCmd.Start(); err != nil {
		ShowError("Failed to start Vue dev server")
		goCmd.Process.Kill()
		os.Exit(1)
	}

	ShowSuccess("Vue dev server starting on http://localhost:3100")
	fmt.Println()
	ShowProgress("Press Ctrl+C to stop both servers")
	fmt.Println()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println()
		ShowProgress("Shutting down servers...")
		goCmd.Process.Kill()
		vueCmd.Process.Kill()
		os.Exit(0)
	}()

	// Wait for both processes
	go goCmd.Wait()
	vueCmd.Wait()
}