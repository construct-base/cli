package construct

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// TemplateData holds all data needed for code generation
type TemplateData struct {
	ResourceName     string
	LowerResourceName string
	PluralName       string
	LowerPluralName  string
	ModuleName       string
	DisplayField     string
	Fields           []TemplateField
}

// TemplateField represents a field in the structure
type TemplateField struct {
	Name           string
	FieldName      string // PascalCase for Go
	Label          string
	Type           string
	TypeScriptType string
	GoType         string
	IsBool         bool
	IsPointer      bool
	Sortable       bool
	ZeroValue      string
	TrueLabel      string
	FalseLabel     string
}

// NewTemplateData creates template data from resource name and fields
func NewTemplateData(resourceName string, fieldArgs []string) *TemplateData {
	pluralName := pluralize(resourceName)
	displayField := "title" // Default, can be made smarter

	fields := parseFieldsToTemplateFields(fieldArgs)

	return &TemplateData{
		ResourceName:      resourceName,
		LowerResourceName: strings.ToLower(resourceName),
		PluralName:        pluralName,
		LowerPluralName:   strings.ToLower(pluralName),
		ModuleName:        strings.ToLower(pluralName),
		DisplayField:      displayField,
		Fields:            fields,
	}
}

func parseFieldsToTemplateFields(fieldArgs []string) []TemplateField {
	var fields []TemplateField

	for _, f := range fieldArgs {
		parts := strings.Split(f, ":")
		if len(parts) < 2 {
			continue
		}

		name := parts[0]
		goType := parts[1]
		tsType := mapGoTypeToTypeScript(goType)
		isBool := goType == "bool" || goType == "boolean"
		isPointer := false

		// For update requests, booleans are pointers
		if isBool {
			isPointer = true
		}

		fieldName := titleCase(name)
		label := titleCase(name)

		zeroValue := getZeroValue(goType)

		fields = append(fields, TemplateField{
			Name:           name,
			FieldName:      fieldName,
			Label:          label,
			Type:           goType,
			TypeScriptType: tsType,
			GoType:         goType,
			IsBool:         isBool,
			IsPointer:      isPointer,
			Sortable:       true,
			ZeroValue:      zeroValue,
			TrueLabel:      "Yes",
			FalseLabel:     "No",
		})
	}

	return fields
}

func mapGoTypeToTypeScript(goType string) string {
	switch goType {
	case "string", "text":
		return "string"
	case "int", "uint", "int64", "uint64", "float", "float64":
		return "number"
	case "bool", "boolean":
		return "boolean"
	case "date", "datetime", "time":
		return "string"
	default:
		return "string"
	}
}

func getZeroValue(goType string) string {
	switch goType {
	case "string", "text":
		return "\"\""
	case "int", "uint", "int64", "uint64", "float", "float64":
		return "0"
	case "bool", "boolean":
		return "false"
	default:
		return "\"\""
	}
}

// GenerateVueFiles generates all Vue files for a structure
func GenerateVueFiles(root, resourceName string, fields []string) error {
	data := NewTemplateData(resourceName, fields)

	vueDir := filepath.Join(root, "vue")
	structureDir := filepath.Join(vueDir, "structures", data.LowerPluralName)

	// Create structure directory
	if err := os.MkdirAll(structureDir, 0755); err != nil {
		return fmt.Errorf("failed to create structure directory: %w", err)
	}

	// Generate index.vue
	if err := generateFromTemplate(
		filepath.Join(structureDir, "index.vue"),
		vueIndexTemplate,
		data,
	); err != nil {
		return fmt.Errorf("failed to generate index.vue: %w", err)
	}

	// Generate composable.ts
	if err := generateFromTemplate(
		filepath.Join(structureDir, "composable.ts"),
		vueComposableTemplate,
		data,
	); err != nil {
		return fmt.Errorf("failed to generate composable.ts: %w", err)
	}

	// Generate types.ts (in structure)
	if err := generateFromTemplate(
		filepath.Join(structureDir, "types.ts"),
		vueTypesTemplate,
		data,
	); err != nil {
		return fmt.Errorf("failed to generate types.ts: %w", err)
	}

	// Also generate types in view/types for global access
	globalTypesDir := filepath.Join(vueDir, "view", "types")
	if err := os.MkdirAll(globalTypesDir, 0755); err != nil {
		return fmt.Errorf("failed to create types directory: %w", err)
	}

	if err := generateFromTemplate(
		filepath.Join(globalTypesDir, data.LowerResourceName+".ts"),
		vueTypesTemplate,
		data,
	); err != nil {
		return fmt.Errorf("failed to generate global types: %w", err)
	}

	return nil
}

func generateFromTemplate(outputPath, templateContent string, data interface{}) error {
	tmpl, err := template.New("").Parse(templateContent)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}