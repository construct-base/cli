package construct

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:     "generate [resource] [fields...]",
	Aliases: []string{"g", "gen"},
	Short:   "Generate full-stack CRUD",
	Long: `Generate complete CRUD for both Go backend and Vue frontend.

Creates:
  • Go: Model, Controller, Service, Routes
  • Vue: Page with list, create, edit, delete modals
  • Automatic API integration

Example:
  construct g Post title:string content:text published:bool
  construct g Article title:string category_id:uint author:belongsTo:User`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runGenerate(args)
	},
}

func runGenerate(args []string) {
	printBanner()
	fmt.Println("🔧 Generating full-stack CRUD...")
	fmt.Println()

	resourceName := args[0]
	fields := args[1:]

	root, err := findProjectRoot()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}

	// Step 1: Generate Go backend using base CLI
	fmt.Println("🔷 Generating Go backend...")
	if err := generateGoBackend(root, resourceName, fields); err != nil {
		fmt.Printf("❌ Go generation failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Go backend generated")
	fmt.Println()

	// Step 2: Generate Vue frontend
	fmt.Println("🟢 Generating Vue frontend...")
	if err := generateVueFrontend(root, resourceName, fields); err != nil {
		fmt.Printf("❌ Vue generation failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Vue frontend generated")
	fmt.Println()

	fmt.Println("🎉 Full-stack CRUD generated successfully!")
	fmt.Println()
	fmt.Printf("📝 Next steps:\n")
	fmt.Printf("   1. Start dev servers: construct dev\n")
	fmt.Printf("   2. Visit: http://localhost:3100/%s\n", strings.ToLower(pluralize(resourceName)))
	fmt.Printf("   3. API available at: /api/%s\n", strings.ToLower(pluralize(resourceName)))
}

func generateGoBackend(root, resourceName string, fields []string) error {
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

func generateVueFrontend(root, resourceName string, fields []string) error {
	return GenerateVueFiles(root, resourceName, fields)
}

type Field struct {
	Name     string
	Type     string
	VueType  string
	Relation string
}

func parseFields(fields []string) []Field {
	var result []Field

	for _, f := range fields {
		parts := strings.Split(f, ":")
		if len(parts) < 2 {
			continue
		}

		name := parts[0]
		goType := parts[1]
		vueType := mapGoTypeToVue(goType)
		relation := ""

		if len(parts) > 2 {
			relation = parts[2]
		}

		result = append(result, Field{
			Name:     name,
			Type:     goType,
			VueType:  vueType,
			Relation: relation,
		})
	}

	return result
}

func mapGoTypeToVue(goType string) string {
	switch goType {
	case "string", "text":
		return "string"
	case "int", "uint", "int64", "uint64":
		return "number"
	case "bool", "boolean":
		return "boolean"
	case "float", "float64":
		return "number"
	case "date", "datetime", "time":
		return "string"
	default:
		return "string"
	}
}

func generateVuePage(path, resourceName string, fields []Field) error {
	// Generate Vue SFC with CRUD modals
	content := generateVuePageTemplate(resourceName, fields)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(content), 0644)
}

func generateVuePageTemplate(resourceName string, fields []Field) string {
	pluralName := pluralize(resourceName)
	lowerPlural := strings.ToLower(pluralName)
	lowerSingular := strings.ToLower(resourceName)

	return fmt.Sprintf(`<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { use%s } from './composable'
import type { %s } from './types'

const { %s, loading, fetch%s, create%s, update%s, delete%s } = use%s()

onMounted(() => {
  fetch%s()
})

const columns = [
%s
  {
    key: 'actions',
    label: 'Actions'
  }
]

const showAddModal = ref(false)
const showEditModal = ref(false)
const showDeleteModal = ref(false)
const selectedItem = ref<% | null>(null)

const handleEdit = (item: %s) => {
  selectedItem.value = item
  showEditModal.value = true
}

const handleDelete = (item: %s) => {
  selectedItem.value = item
  showDeleteModal.value = true
}
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar title="%s">
        <template #right>
          <UButton @click="showAddModal = true" icon="i-lucide-plus">
            Add %s
          </UButton>
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <UTable
        :rows="%s"
        :columns="columns"
        :loading="loading"
      >
        <template #actions-data="{ row }">
          <div class="flex gap-2">
            <UButton
              size="xs"
              color="primary"
              variant="ghost"
              icon="i-lucide-pencil"
              @click="handleEdit(row)"
            />
            <UButton
              size="xs"
              color="error"
              variant="ghost"
              icon="i-lucide-trash"
              @click="handleDelete(row)"
            />
          </div>
        </template>
      </UTable>
    </template>
  </UDashboardPanel>

  <!-- Add Modal -->
  <UModal v-model="showAddModal" title="Add %s">
    <!-- TODO: Add form -->
  </UModal>

  <!-- Edit Modal -->
  <UModal v-model="showEditModal" title="Edit %s">
    <!-- TODO: Edit form -->
  </UModal>

  <!-- Delete Modal -->
  <UModal v-model="showDeleteModal" title="Delete %s">
    <!-- TODO: Delete confirmation -->
  </UModal>
</template>
`,
		pluralName,
		pluralName,
		resourceName,
		lowerSingular,
		lowerPlural,
		pluralName,
		resourceName,
		resourceName,
		resourceName,
		pluralName,
		generateColumnsList(fields),
		resourceName,
		resourceName,
		resourceName,
		pluralName,
		resourceName,
		lowerPlural,
		resourceName,
		resourceName,
		resourceName,
	)
}

func generateColumnsList(fields []Field) string {
	var columns []string
	for _, f := range fields {
		columns = append(columns, fmt.Sprintf("  { key: '%s', label: '%s' },",
			f.Name, titleCase(f.Name)))
	}
	return strings.Join(columns, "\n")
}

func generateVueTypes(path, resourceName string, fields []Field) error {
	content := generateVueTypesContent(resourceName, fields)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(content), 0644)
}

func generateVueTypesContent(resourceName string, fields []Field) string {
	var fieldDefs []string
	for _, f := range fields {
		fieldDefs = append(fieldDefs, fmt.Sprintf("  %s: %s", f.Name, f.VueType))
	}

	return fmt.Sprintf(`export interface %s {
  id: number
%s
  created_at: string
  updated_at: string
}

export interface %sCreateRequest {
%s
}

export interface %sUpdateRequest {
%s
}
`,
		resourceName,
		strings.Join(fieldDefs, "\n"),
		resourceName,
		strings.Join(fieldDefs, "\n"),
		resourceName,
		strings.Join(fieldDefs, "\n"),
	)
}

func generateVueComposable(path, resourceName string, fields []Field) error {
	content := generateVueComposableContent(resourceName)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(content), 0644)
}

func generateVueComposableContent(resourceName string) string {
	pluralName := pluralize(resourceName)
	lowerPlural := strings.ToLower(pluralName)

	return fmt.Sprintf(`import { ref } from 'vue'
import { apiClient } from '@core/api/client'
import type { %s, %sCreateRequest, %sUpdateRequest } from '@/types/%s'

export function use%s() {
  const %s = ref<%s[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const fetch%s = async () => {
    loading.value = true
    error.value = null
    try {
      const response = await apiClient.get('/%s')
      %s.value = response.data
    } catch (e: any) {
      error.value = e.message
    } finally {
      loading.value = false
    }
  }

  const create%s = async (data: %sCreateRequest) => {
    loading.value = true
    error.value = null
    try {
      await apiClient.post('/%s', data)
      await fetch%s()
    } catch (e: any) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  const update%s = async (id: number, data: %sUpdateRequest) => {
    loading.value = true
    error.value = null
    try {
      await apiClient.put('/%s/' + id, data)
      await fetch%s()
    } catch (e: any) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  const delete%s = async (id: number) => {
    loading.value = true
    error.value = null
    try {
      await apiClient.delete('/%s/' + id)
      await fetch%s()
    } catch (e: any) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  return {
    %s,
    loading,
    error,
    fetch%s,
    create%s,
    update%s,
    delete%s
  }
}
`,
		resourceName,
		resourceName,
		resourceName,
		strings.ToLower(resourceName),
		pluralName,
		lowerPlural,
		resourceName,
		pluralName,
		lowerPlural,
		lowerPlural,
		resourceName,
		resourceName,
		lowerPlural,
		pluralName,
		resourceName,
		resourceName,
		lowerPlural,
		pluralName,
		resourceName,
		lowerPlural,
		pluralName,
		lowerPlural,
		pluralName,
		resourceName,
		resourceName,
		resourceName,
	)
}

func pluralize(word string) string {
	// Simple pluralization
	if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "x") ||
	   strings.HasSuffix(word, "ch") || strings.HasSuffix(word, "sh") {
		return word + "es"
	}
	if strings.HasSuffix(word, "y") {
		return word[:len(word)-1] + "ies"
	}
	return word + "s"
}

func titleCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}