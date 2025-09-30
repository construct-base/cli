package construct

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start development servers",
	Long:  "Start both Go API (port 8100) and Vue dev server (port 3100) with hot reload",
	Run: func(cmd *cobra.Command, args []string) {
		runDev()
	},
}

func runDev() {
	printBanner()
	fmt.Println("🚀 Starting development servers...")
	fmt.Println()

	root, err := findProjectRoot()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}

	// Start Go API
	goCmd := exec.Command("go", "run", "main.go")
	goCmd.Dir = root
	goCmd.Stdout = NewPrefixWriter(os.Stdout, "🔷 [Go]   ")
	goCmd.Stderr = NewPrefixWriter(os.Stderr, "🔷 [Go]   ")
	goCmd.Env = os.Environ()

	if err := goCmd.Start(); err != nil {
		fmt.Printf("❌ Failed to start Go server: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Go API starting on http://localhost:8100")

	// Start Vue dev server
	vueDir := filepath.Join(root, "vue")
	vueCmd := detectPackageManager(vueDir, "dev")
	vueCmd.Dir = vueDir
	vueCmd.Stdout = NewPrefixWriter(os.Stdout, "🟢 [Vue]  ")
	vueCmd.Stderr = NewPrefixWriter(os.Stderr, "🟢 [Vue]  ")
	vueCmd.Env = os.Environ()

	if err := vueCmd.Start(); err != nil {
		fmt.Printf("❌ Failed to start Vue dev server: %v\n", err)
		goCmd.Process.Kill()
		os.Exit(1)
	}

	fmt.Println("✅ Vue dev server starting on http://localhost:3100")
	fmt.Println()
	fmt.Println("📝 Press Ctrl+C to stop both servers")
	fmt.Println()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n\n🛑 Shutting down...")
		goCmd.Process.Kill()
		vueCmd.Process.Kill()
		os.Exit(0)
	}()

	// Wait for both processes
	go goCmd.Wait()
	vueCmd.Wait()
}