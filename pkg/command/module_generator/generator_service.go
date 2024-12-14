package module_generator

import (
	"errors"
	"fmt"
	"golang.org/x/mod/modfile"
	"os"
	"path/filepath"
)

type GeneratorService struct {
	fc *FileCreator
}

func NewGeneratorService(fc *FileCreator) *GeneratorService {
	return &GeneratorService{fc: fc}
}

func (g GeneratorService) Execute(moduleName string) error {
	// Base path for the module
	basePath := filepath.Join("internal", "app", "module", moduleName)

	// Check if the module already exists
	if _, err := os.Stat(basePath); !os.IsNotExist(err) {
		return errors.New("module already exists")
	}

	pn, err := getProjectPackageName()
	if err != nil {
		return err
	}

	err = g.fc.Create(moduleName, basePath, pn)
	if err != nil {
		return err
	}

	return nil
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
