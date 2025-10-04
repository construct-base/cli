package construct

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/base-go/mamba"
)

var newCmd = &mamba.Command{
	Use:   "new [project_name]",
	Short: "Create a new Construct project",
	Long: `Create a new Construct project by cloning the framework and setting up the directory.

This will clone the latest Construct framework and set up a new project with:
  • Go backend with HMVC architecture
  • Vue 3 frontend with file-based routing
  • Development and production tooling
  • Example structures and documentation

Example:
  construct new my-blog
  construct new ecommerce-app`,
	Run: func(cmd *mamba.Command, args []string) {
		if len(args) != 1 {
			cmd.Help()
			os.Exit(1)
		}
		runNewProject(cmd, args)
	},
}

func runNewProject(cmd *mamba.Command, args []string) {
	projectName := args[0]

	printBanner()
	fmt.Printf("🏗️  Creating new Construct project: %s\n\n", projectName)

	// Repository URL - clone from latest release
	repoURL := "https://github.com/construct-base/core.git"
	branch := "--branch=main" // Use main branch, releases are tagged

	// Check if project directory already exists
	if _, err := os.Stat(projectName); err == nil {
		fmt.Printf("❌ Error: Directory '%s' already exists\n", projectName)
		os.Exit(1)
	}

	// Check if git is installed
	if _, err := exec.LookPath("git"); err != nil {
		fmt.Printf("❌ Error: git is not installed or not in PATH\n")
		fmt.Printf("   Please install git: https://git-scm.com/downloads\n")
		os.Exit(1)
	}

	fmt.Println("📥 Cloning Construct framework from latest release...")
	fmt.Println("   This will download a clean template with:")
	fmt.Println("   • Go backend (Base framework)")
	fmt.Println("   • Vue 3 frontend")
	fmt.Println("   • Development tooling")
	fmt.Println()

	// Clone the repository with depth 1 for faster download
	cloneCmd := exec.Command("git", "clone", branch, "--depth=1", repoURL, projectName)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr

	if err := cloneCmd.Run(); err != nil {
		fmt.Printf("❌ Error cloning framework: %v\n", err)
		fmt.Printf("   Make sure you have internet connection and access to:\n")
		fmt.Printf("   %s\n", repoURL)
		os.Exit(1)
	}

	fmt.Println()

	// Remove .git directory to make it a fresh project
	fmt.Println("🔧 Setting up project...")
	gitDir := filepath.Join(projectName, ".git")
	if err := os.RemoveAll(gitDir); err != nil {
		fmt.Printf("⚠️  Warning: Could not remove .git directory: %v\n", err)
	}

	// Copy .env.example to .env
	fmt.Println("⚙️  Setting up environment...")
	envExamplePath := filepath.Join(projectName, ".env.example")
	envPath := filepath.Join(projectName, ".env")
	if err := copyFile(envExamplePath, envPath); err != nil {
		fmt.Printf("⚠️  Warning: Could not create .env file: %v\n", err)
	} else {
		fmt.Println("   Created .env file")
	}

	// Initialize new git repository
	fmt.Println("   Initializing git repository...")
	initCmd := exec.Command("git", "init")
	initCmd.Dir = projectName
	if err := initCmd.Run(); err != nil {
		fmt.Printf("⚠️  Warning: Could not initialize git repository: %v\n", err)
	} else {
		fmt.Println("   ✓ Git repository initialized")
	}

	// Install Vue dependencies
	fmt.Println()
	fmt.Println("📦 Installing Vue dependencies...")

	vueDir := filepath.Join(projectName, "vue")

	// Detect package manager (bun > pnpm > yarn > npm)
	var installCmd *exec.Cmd
	if _, err := exec.LookPath("bun"); err == nil {
		fmt.Println("   Using bun...")
		installCmd = exec.Command("bun", "install")
	} else if _, err := exec.LookPath("pnpm"); err == nil {
		fmt.Println("   Using pnpm...")
		installCmd = exec.Command("pnpm", "install")
	} else if _, err := exec.LookPath("yarn"); err == nil {
		fmt.Println("   Using yarn...")
		installCmd = exec.Command("yarn", "install")
	} else if _, err := exec.LookPath("npm"); err == nil {
		fmt.Println("   Using npm...")
		installCmd = exec.Command("npm", "install")
	} else {
		fmt.Println("⚠️  Warning: No package manager found (bun/pnpm/yarn/npm)")
		fmt.Println("   You'll need to manually install dependencies:")
		fmt.Printf("   cd %s/vue && npm install\n", projectName)
	}

	if installCmd != nil {
		installCmd.Dir = vueDir
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr

		if err := installCmd.Run(); err != nil {
			fmt.Printf("⚠️  Warning: Could not install dependencies: %v\n", err)
			fmt.Println("   You may need to install them manually:")
			fmt.Printf("   cd %s/vue && npm install\n", projectName)
		} else {
			fmt.Println("   ✓ Dependencies installed")
		}
	}

	// Get absolute path
	absPath, err := filepath.Abs(projectName)
	if err != nil {
		absPath = projectName
	}

	fmt.Println()
	fmt.Println("✅ Project created successfully!")
	fmt.Println()
	fmt.Printf("📁 Location: %s\n", absPath)
	fmt.Println()
	fmt.Println("📝 Next steps:")
	fmt.Printf("   cd %s\n", projectName)
	fmt.Println("   construct dev              # Start development servers")
	fmt.Println()
	fmt.Println("💡 Or generate your first CRUD module:")
	fmt.Println("   construct g Post title:string content:text published:bool")
	fmt.Println()
}