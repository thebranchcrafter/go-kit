package module_generator

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/iancoleman/strcase"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	directories = []string{"domain", "infrastructure"}
	templates   = map[string]string{
		"domain/{{.PackageName}}.go":                             "domain/module_aggregate.tmpl",
		"domain/{{.PackageName}}_repository.go":                  "domain/module_repository.tmpl",
		"infrastructure/postgres_{{.PackageName}}_repository.go": "infrastructure/module_postgres_repository.tmpl",
		"module.go": "module.tmpl",
	}
)

type FileCreator struct {
	tmplDir string
}

func NewFileCreator(tmplDir string) *FileCreator {
	return &FileCreator{
		tmplDir: tmplDir,
	}
}

func (c FileCreator) Create(moduleName string, basePath string, packageName string) error {
	err := c.createDirectories(basePath)
	if err != nil {
		return err
	}

	// Template data
	data := map[string]string{
		"ClassName":      strcase.ToCamel(moduleName),
		"ClassInitial":   toAbbreviation(moduleName),
		"FullModuleName": packageName,
		"ModuleName":     strings.Title(moduleName),
		"PackageName":    strings.ToLower(moduleName),
		"RepositoryName": fmt.Sprintf("%sRepository", strings.Title(moduleName)),
	}

	// Create files using templates
	for filePathTemplate, fileTemplatePath := range templates {
		// Render the file path
		filePath := renderTemplate(filePathTemplate, data)

		fullTemplatePath := c.tmplDir + fileTemplatePath
		// Read the file template content
		fileContent, err := readTemplateFile(fullTemplatePath)
		if err != nil {
			return fmt.Errorf("failed to read template file '%s': %w", fullTemplatePath, err)
		}

		// Render the file content
		renderedContent := renderTemplate(fileContent, data)

		// Full file path
		fullPath := filepath.Join(basePath, filePath)

		// Ensure the directory exists
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory '%s': %w", dir, err)
		}

		// Write the rendered content to the file
		if err := os.WriteFile(fullPath, []byte(renderedContent), 0644); err != nil {
			return fmt.Errorf("failed to write file '%s': %w", fullPath, err)
		}

		fmt.Printf("Created file: %s\n", fullPath)
	}

	return nil
}

func (c *FileCreator) createDirectories(basePath string) error {
	for _, dir := range directories {
		path := filepath.Join(basePath, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
		color.New(color.FgCyan).Printf("Created directory: %s\n", path)
	}

	return nil
}

// readTemplateFile reads the content of a template file
func readTemplateFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("could not read template file: %w", err)
	}
	return string(content), nil
}

// renderTemplate renders a template string with the provided data
func renderTemplate(templateStr string, data map[string]string) string {
	tmpl, err := template.New("template").Funcs(template.FuncMap{
		"Title": strings.Title,
	}).Parse(templateStr)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse template: %s", err))
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		panic(fmt.Sprintf("Failed to execute template: %s", err))
	}

	return buf.String()
}

func toAbbreviation(input string) string {
	parts := strings.Split(input, "_")
	var abbreviation string
	for _, part := range parts {
		if len(part) > 0 {
			abbreviation += string(part[0])
		}
	}
	return strings.ToLower(abbreviation)
}
