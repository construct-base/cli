package construct

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
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
	Args: cobra.ExactArgs(1),
	Run:  runNewProject,
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func runNewProject(cmd *cobra.Command, args []string) {
	projectName := args[0]

	printBanner()
	fmt.Printf("🏗️  Creating new Construct project: %s\n\n", projectName)

	// Repository URL
	repoURL := "https://github.com/construct-base/core.git"

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

	fmt.Println("📥 Cloning Construct framework...")

	// Clone the repository
	cloneCmd := exec.Command("git", "clone", repoURL, projectName)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr

	if err := cloneCmd.Run(); err != nil {
		fmt.Printf("❌ Error cloning framework: %v\n", err)
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
	initCmd := exec.Command("git", "init")
	initCmd.Dir = projectName
	if err := initCmd.Run(); err != nil {
		fmt.Printf("⚠️  Warning: Could not initialize git repository: %v\n", err)
	} else {
		fmt.Println("   Initialized new git repository")
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
	fmt.Println("   construct g Post title:string content:text  # Generate CRUD")
	fmt.Println()
}