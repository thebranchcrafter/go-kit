package module_domain

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed templates/*
var templatesFS embed.FS

var (
	directories = []string{"domain", "infrastructure"}
	templates   = map[string]string{
		"domain/{{.PackageName}}.go":                             "domain/module_aggregate.tmpl",
		"domain/{{.PackageName}}_repository.go":                  "domain/module_repository.tmpl",
		"infrastructure/postgres_{{.PackageName}}_repository.go": "infrastructure/module_postgres_repository.tmpl",
		"module.go": "module.tmpl",
	}
)

type TemplateGenerator struct {
}

func NewTemplateGenerator() *TemplateGenerator {
	return &TemplateGenerator{}
}

func (vg *TemplateGenerator) Generate(basePath string, variables *ModuleGeneratorVariables) error {
	err := vg.createDirectories(basePath)
	if err != nil {
		return err
	}

	for filePathTemplate, fileTemplatePath := range templates {
		// Render the file path
		filePath := renderTemplate(filePathTemplate, variables.ToArray())

		// Read the template content from embedded FS
		fileContent, err := readTemplateFile(fileTemplatePath)
		if err != nil {
			return fmt.Errorf("failed to read template file '%s': %w", fileTemplatePath, err)
		}

		// Render the file content
		renderedContent := renderTemplate(fileContent, variables.ToArray())

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

func (vg *TemplateGenerator) createDirectories(basePath string) error {
	for _, dir := range directories {
		path := filepath.Join(basePath, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
		_, _ = color.New(color.FgCyan).Printf("Created directory: %s\n", path)
	}

	return nil
}

// readTemplateFile reads the content of a template file from the embedded FS
func readTemplateFile(filePath string) (string, error) {
	content, err := templatesFS.ReadFile(fmt.Sprintf("templates/%s", filePath))
	if err != nil {
		return "", fmt.Errorf("could not read template file: %w", err)
	}
	return string(content), nil
}

// renderTemplate renders a template string with the provided data
func renderTemplate(templateStr string, data map[string]string) string {
	titleCaser := cases.Title(language.English)

	tmpl, err := template.New("template").Funcs(template.FuncMap{
		"Title": titleCaser.String,
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
