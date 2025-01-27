package module_application

import (
	"context"
	"errors"
	module_domain "github.com/thebranchcrafter/go-kit/internal/module/domain"
	"github.com/thebranchcrafter/go-kit/pkg/application"
	"os"
	"path/filepath"
)

type GenerateModuleCommand struct {
	ModuleName string
}

func (g GenerateModuleCommand) Id() string {
	return "generate-module-command"
}

type GenerateModuleCommandHandler struct {
	vg *module_domain.VariablesGenerator
	tg *module_domain.TemplateGenerator
}

func NewGenerateModuleCommandHandler(vg *module_domain.VariablesGenerator, tg *module_domain.TemplateGenerator) *GenerateModuleCommandHandler {
	return &GenerateModuleCommandHandler{vg: vg, tg: tg}
}

func (g GenerateModuleCommandHandler) Handle(_ context.Context, c application.Command) error {
	cmd, ok := c.(GenerateModuleCommand)
	if !ok {
		return errors.New("invalid command")
	}

	v, err := g.vg.Generate(cmd.ModuleName)
	if err != nil {
		return err
	}

	// Base path for the module
	basePath := filepath.Join("internal", "app", "module", cmd.ModuleName)

	// Check if the module already exists
	if _, err := os.Stat(basePath); !os.IsNotExist(err) {
		return errors.New("module already exists")
	}

	return g.tg.Generate(basePath, v)
}
