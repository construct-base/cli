package construct

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GenerateBackend generates Go backend files using Base CLI
func GenerateBackend(root, resourceName string, fields []string) error {
	// Use base CLI to generate backend in api/ directory
	args := append([]string{"generate", resourceName}, fields...)

	cmd := exec.Command("base", args...)
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Set environment to tell base to use api/ instead of app/
	cmd.Env = append(os.Environ(), "BASE_MODULE_DIR=api")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("base CLI failed: %w", err)
	}

	// Move from app/ to api/ if base created in app/
	appModulePath := filepath.Join(root, "app", strings.ToLower(pluralize(resourceName)))
	apiModulePath := filepath.Join(root, "api", strings.ToLower(pluralize(resourceName)))

	if _, err := os.Stat(appModulePath); err == nil {
		// Move module from app/ to api/
		if err := os.Rename(appModulePath, apiModulePath); err != nil {
			return fmt.Errorf("failed to move module to api/: %w", err)
		}

		// Update api/init.go instead of app/init.go
		updateInitFile(filepath.Join(root, "api", "init.go"), resourceName)
	}

	return nil
}

func updateInitFile(initPath string, resourceName string) error {
	moduleName := strings.ToLower(pluralize(resourceName))

	// Read current init.go
	content, err := os.ReadFile(initPath)
	if err != nil {
		return err
	}

	contentStr := string(content)

	// Check if module already exists
	moduleInit := fmt.Sprintf("modules[\"%s\"] = %s.Init(deps)", moduleName, moduleName)
	if strings.Contains(contentStr, moduleInit) {
		fmt.Printf("✅ Module already registered in api/init.go\n")
		return nil // Already added
	}

	// Add import
	importLine := fmt.Sprintf("\"base/api/%s\"", moduleName)
	if !strings.Contains(contentStr, importLine) {
		// Find the import block and add the import
		importBlockStart := strings.Index(contentStr, "import (")
		if importBlockStart == -1 {
			return fmt.Errorf("could not find import block in api/init.go")
		}

		// Find the end of existing imports (before ")")
		importBlockEnd := strings.Index(contentStr[importBlockStart:], ")")
		if importBlockEnd == -1 {
			return fmt.Errorf("could not find end of import block in api/init.go")
		}
		importBlockEnd += importBlockStart

		// Insert new import before the closing )
		newImport := fmt.Sprintf("\t%s\n", importLine)
		contentStr = contentStr[:importBlockEnd] + newImport + contentStr[importBlockEnd:]
	}

	// Add module initialization
	returnIndex := strings.Index(contentStr, "return modules")
	if returnIndex == -1 {
		return fmt.Errorf("could not find 'return modules' in api/init.go")
	}

	// Insert the module initialization before return
	moduleInitLine := fmt.Sprintf("\t// %s module\n\t%s\n\n\t", titleCase(moduleName), moduleInit)
	contentStr = contentStr[:returnIndex] + moduleInitLine + contentStr[returnIndex:]

	// Write back to file
	if err := os.WriteFile(initPath, []byte(contentStr), 0644); err != nil {
		return fmt.Errorf("failed to write api/init.go: %w", err)
	}

	fmt.Printf("✅ Added module to api/init.go\n")
	return nil
}
