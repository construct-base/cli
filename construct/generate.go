package construct

import (
	"fmt"
	"os"
	"strings"

	"github.com/base-go/mamba"
)

var generateCmd = &mamba.Command{
	Use:     "generate [resource] [fields...]",
	Aliases: []string{"g", "gen", "g:b", "gen:b", "g:f", "gen:f", "generate:b", "generate:f"},
	Short:   "Generate full-stack CRUD",
	Long: `Generate complete CRUD for both Go backend and Vue frontend.

Creates:
  • Go: Model, Controller, Service, Routes
  • Vue: Page, Store, Composables, Modals, Types
  • Automatic API integration

Examples:
  construct g Post title:string content:text published:bool
  construct g:b Product name:string price:float stock:uint
  construct g:f Category name:string description:text

Syntax:
  g or generate    Generate both backend and frontend
  g:b or gen:b     Generate backend only
  g:f or gen:f     Generate frontend only`,
	Run: func(cmd *mamba.Command, args []string) {
		// Validate args
		if len(args) < 1 {
			cmd.Help()
			os.Exit(1)
		}
		runGenerate(cmd.Name(), args)
	},
}

func runGenerate(command string, args []string) {
	printBanner()

	resourceName := args[0]
	fields := args[1:]

	root, err := findProjectRoot()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}

	// Determine what to generate based on command suffix
	generateBackend := true
	generateFrontend := true

	// Check for :b or :f suffix
	if strings.HasSuffix(command, ":b") {
		generateBackend = true
		generateFrontend = false
		fmt.Println("🔧 Generating backend only...")
	} else if strings.HasSuffix(command, ":f") {
		generateBackend = false
		generateFrontend = true
		fmt.Println("🔧 Generating frontend only...")
	} else {
		fmt.Println("🔧 Generating full-stack CRUD...")
	}
	fmt.Println()

	// Step 1: Generate Go backend (if needed)
	if generateBackend {
		fmt.Println("🔷 Generating Go backend...")
		if err := GenerateBackend(root, resourceName, fields); err != nil {
			fmt.Printf("❌ Go generation failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Go backend generated")
		fmt.Println()
	}

	// Step 2: Generate Vue frontend (if needed)
	if generateFrontend {
		fmt.Println("🟢 Generating Vue frontend...")
		if err := GenerateFrontend(root, resourceName, fields); err != nil {
			fmt.Printf("❌ Vue generation failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Vue frontend generated")
		fmt.Println()
	}

	// Success message
	if generateBackend && generateFrontend {
		fmt.Println("🎉 Full-stack CRUD generated successfully!")
	} else if generateBackend {
		fmt.Println("🎉 Backend generated successfully!")
	} else {
		fmt.Println("🎉 Frontend generated successfully!")
	}

	fmt.Println()
	fmt.Printf("📝 Next steps:\n")
	if generateBackend && generateFrontend {
		fmt.Printf("   1. Start dev servers: construct dev\n")
		fmt.Printf("   2. Visit: http://localhost:3100/%s\n", strings.ToLower(pluralize(resourceName)))
		fmt.Printf("   3. API available at: /api/%s\n", strings.ToLower(pluralize(resourceName)))
	} else if generateBackend {
		fmt.Printf("   1. Test API: curl http://localhost:8100/api/%s\n", strings.ToLower(pluralize(resourceName)))
		fmt.Printf("   2. Generate frontend: construct g %s %s --frontend\n", resourceName, strings.Join(fields, " "))
	} else {
		fmt.Printf("   1. Generate backend: construct g %s %s --backend\n", resourceName, strings.Join(fields, " "))
		fmt.Printf("   2. Start dev: construct dev\n")
	}
}
