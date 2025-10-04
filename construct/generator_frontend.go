package construct

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// GenerateFrontend generates all Vue frontend files in self-contained module structure
func GenerateFrontend(root, resourceName string, fields []string) error {
	data := NewTemplateData(resourceName, fields)

	// Self-contained module under app/{module}/
	moduleDir := filepath.Join(root, "vue", "app", data.LowerPluralName)

	// 1. Generate types in app/{module}/types/
	if err := generateFileFromTemplate(
		filepath.Join(moduleDir, "types", data.LowerResourceName+".ts"),
		vueTypesTemplate,
		data,
	); err != nil {
		return fmt.Errorf("failed to generate types: %w", err)
	}
	fmt.Printf("  ✓ Generated app/%s/types/%s.ts\n", data.LowerPluralName, data.LowerResourceName)

	// 2. Generate composable in app/{module}/composables/
	if err := generateFileFromTemplate(
		filepath.Join(moduleDir, "composables", "use"+data.PluralName+".ts"),
		vueComposableTemplate,
		data,
	); err != nil {
		return fmt.Errorf("failed to generate composable: %w", err)
	}
	fmt.Printf("  ✓ Generated app/%s/composables/use%s.ts\n", data.LowerPluralName, data.PluralName)

	// 3. Generate store in app/{module}/stores/
	if err := generateFileFromTemplate(
		filepath.Join(moduleDir, "stores", data.LowerPluralName+".ts"),
		vueStoreTemplate,
		data,
	); err != nil {
		return fmt.Errorf("failed to generate store: %w", err)
	}
	fmt.Printf("  ✓ Generated app/%s/stores/%s.ts\n", data.LowerPluralName, data.LowerPluralName)

	// 4. Generate modal components in app/{module}/components/
	componentsDir := filepath.Join(moduleDir, "components")
	if err := os.MkdirAll(componentsDir, 0755); err != nil {
		return fmt.Errorf("failed to create components directory: %w", err)
	}

	// Add Modal
	if err := generateFileFromTemplate(
		filepath.Join(componentsDir, data.PluralName+"AddModal.vue"),
		vueAddModalTemplate,
		data,
	); err != nil {
		return fmt.Errorf("failed to generate AddModal: %w", err)
	}
	fmt.Printf("  ✓ Generated app/%s/components/%sAddModal.vue\n", data.LowerPluralName, data.PluralName)

	// Delete Modal
	if err := generateFileFromTemplate(
		filepath.Join(componentsDir, data.PluralName+"DeleteModal.vue"),
		vueDeleteModalTemplate,
		data,
	); err != nil {
		return fmt.Errorf("failed to generate DeleteModal: %w", err)
	}
	fmt.Printf("  ✓ Generated app/%s/components/%sDeleteModal.vue\n", data.LowerPluralName, data.PluralName)

	// 5. Generate main page (index) in app/{module}/pages/
	if err := generateFileFromTemplate(
		filepath.Join(moduleDir, "pages", "index.vue"),
		vueIndexTemplate,
		data,
	); err != nil {
		return fmt.Errorf("failed to generate page: %w", err)
	}
	fmt.Printf("  ✓ Generated app/%s/pages/index.vue\n", data.LowerPluralName)

	return nil
}

func generateFileFromTemplate(outputPath, templateContent string, data interface{}) error {
	// Parse template
	tmpl, err := template.New("").Parse(templateContent)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write file
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

