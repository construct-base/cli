package construct

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build production app",
	Long:  "Build Vue SPA and compile Go binary into dist/ directory for production deployment",
	Run: func(cmd *cobra.Command, args []string) {
		runBuild()
	},
}

func runBuild() {
	printBanner()
	fmt.Println("🔨 Building Construct for production...")
	fmt.Println()

	root, err := findProjectRoot()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}

	// Create dist directory structure
	distDir := filepath.Join(root, "dist")
	publicDir := filepath.Join(distDir, "public")

	fmt.Println("📁 Creating dist/ directory structure...")
	if err := os.RemoveAll(distDir); err != nil {
		fmt.Printf("❌ Failed to clean dist directory: %v\n", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(publicDir, 0755); err != nil {
		fmt.Printf("❌ Failed to create dist directory: %v\n", err)
		os.Exit(1)
	}

	// Build Vue SPA
	fmt.Println("📦 Building Vue SPA → dist/public/")
	vueDir := filepath.Join(root, "vue")
	buildCmd := detectPackageManager(vueDir, "build")
	buildCmd.Dir = vueDir
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	if err := buildCmd.Run(); err != nil {
		fmt.Printf("❌ Vue build failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Vue SPA built successfully")
	fmt.Println()

	// Build Go binary
	fmt.Println("🔷 Building Go binary → dist/construct")

	binaryName := "construct"
	if runtime.GOOS == "windows" {
		binaryName = "construct.exe"
	}

	binaryPath := filepath.Join(distDir, binaryName)
	goCmd := exec.Command("go", "build", "-o", binaryPath, "main.go")
	goCmd.Dir = root
	goCmd.Stdout = os.Stdout
	goCmd.Stderr = os.Stderr

	if err := goCmd.Run(); err != nil {
		fmt.Printf("❌ Go build failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Go binary built successfully")
	fmt.Println()

	// Copy necessary runtime files
	fmt.Println("📋 Copying runtime files...")

	// Copy .env.example if exists
	envExample := filepath.Join(root, ".env.example")
	if _, err := os.Stat(envExample); err == nil {
		copyFile(envExample, filepath.Join(distDir, ".env.example"))
	}

	// Create required directories
	os.MkdirAll(filepath.Join(distDir, "storage"), 0755)
	os.MkdirAll(filepath.Join(distDir, "logs"), 0755)

	fmt.Println("✅ Runtime files copied")
	fmt.Println()
	fmt.Println("🎉 Build complete!")
	fmt.Println()
	fmt.Println("📦 Production build location: dist/")
	fmt.Println()
	fmt.Println("To start production server:")
	fmt.Printf("  cd dist && ./%s\n", binaryName)
	fmt.Println()
	fmt.Println("Or use:")
	fmt.Println("  construct start")
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}