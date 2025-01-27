package module_domain

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"golang.org/x/mod/modfile"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"strings"
)

type ModuleGeneratorVariables struct {
	ClassName      string
	ClassInitial   string
	FullModuleName string
	ModuleName     string
	PackageName    string
	RepositoryName string
}

func (gen *ModuleGeneratorVariables) ToArray() map[string]string {
	return map[string]string{
		"ClassName":      gen.ClassName,
		"ClassInitial":   gen.ClassInitial,
		"FullModuleName": gen.FullModuleName,
		"ModuleName":     gen.ModuleName,
		"PackageName":    gen.PackageName,
		"RepositoryName": gen.RepositoryName,
	}
}

type VariablesGenerator struct {
}

func NewVariablesGenerator() *VariablesGenerator {
	return &VariablesGenerator{}
}

func (vg *VariablesGenerator) Generate(moduleName string) (*ModuleGeneratorVariables, error) {
	pn, err := getProjectPackageName()
	if err != nil {
		return nil, err
	}

	titleCaser := cases.Title(language.English)

	return &ModuleGeneratorVariables{
		ClassName:      strcase.ToCamel(moduleName),
		ClassInitial:   toAbbreviation(moduleName),
		FullModuleName: pn,
		ModuleName:     titleCaser.String(moduleName),
		PackageName:    strings.ToLower(moduleName),
		RepositoryName: fmt.Sprintf("%sRepository", titleCaser.String(moduleName)),
	}, nil
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

func getProjectPackageName() (string, error) {
	// Read the go.mod file
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "", fmt.Errorf("error reading go.mod file: %w", err)
	}

	// Parse the go.mod file
	modFile, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return "", fmt.Errorf("error parsing go.mod file: %w", err)
	}

	// Return the module name
	return modFile.Module.Mod.Path, nil
}
