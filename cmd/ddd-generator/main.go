package main

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	module_application "github.com/thebranchcrafter/go-kit/internal/module/application"
	module_domain "github.com/thebranchcrafter/go-kit/internal/module/domain"
	"github.com/thebranchcrafter/go-kit/pkg/utils"
	"os"
	"strings"
)

var binaryName = "ddd-generator"

func main() {
	// Root command
	rootCmd := &cobra.Command{
		Use:   binaryName,
		Short: "CLI tool for generating modules and components",
	}

	// Add the `generate` subcommand
	rootCmd.AddCommand(NewGenerateCmd())

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		_, _ = color.New(color.FgRed).Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// NewGenerateCmd creates the `generate` command
func NewGenerateCmd() *cobra.Command {
	// Create the `generate` command
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate modules and other components",
	}

	// Add the `module` subcommand to `generate`
	generateCmd.AddCommand(NewGenerateModuleCmd())

	return generateCmd
}

// NewGenerateModuleCmd creates the `generate module` command
func NewGenerateModuleCmd() *cobra.Command {
	var moduleName string

	cmd := &cobra.Command{
		Use:   "module",
		Short: color.New(color.Bold).Sprint("Generate a new module with a predefined folder structure"),
		Long: color.New(color.FgCyan).Sprint(`
Generate a new module with a predefined folder structure for scalable and maintainable applications.

The generated module includes layers for domain logic, infrastructure, and module-specific functionality:

┌───────────────────────────────────────────────┐
│                 Folder Structure              │
├───────────────────────────────────────────────┤
│ internal/app/module/<module-name>             │
│   ├── domain                                  │
│   │   ├── <module-name>.go                    │
│   │   └── <module-name>_repository.go         │
│   ├── infrastructure                          │
│   │   └── postgres_<module-name>_repository.go│
│   ├── module.go                               │
└───────────────────────────────────────────────┘

Details:
- domain:          Contains core business logic and interfaces.
- infrastructure:  Includes adapters for database and external services.
- module.go:       Entry point for module-specific logic and initialization.

Example:
Generate a module named 'user':
  internal/app/module/user/
    ├── domain/user.go
    ├── domain/user_repository.go
    ├── infrastructure/postgres_user_repository.go
    ├── module.go
`),
		Run: func(cmd *cobra.Command, args []string) {
			// Prompt user for module name
			utils.PrintYellow("Enter the module name: ")
			_, err := fmt.Scanln(&moduleName)
			if err != nil {
				panic(err)
			}

			// Trim and validate input
			moduleName = strings.ToLower(strings.TrimSpace(moduleName))
			if moduleName == "" {
				utils.PrintRed("Error: Module name cannot be empty")
				os.Exit(1)
			}

			// Print progress
			utils.PrintCyan(fmt.Sprintf("Generating module: %s", moduleName))

			handler := module_application.NewGenerateModuleCommandHandler(
				module_domain.NewVariablesGenerator(),
				module_domain.NewTemplateGenerator(),
			)

			// Handle errors
			if err := handler.Handle(context.Background(), module_application.GenerateModuleCommand{ModuleName: moduleName}); err != nil {
				utils.PrintRed(fmt.Sprintf("Error generating module: %v", err))
				os.Exit(1)
			}

			// Print success
			utils.PrintGreen(fmt.Sprintf("Module '%s' generated successfully!", moduleName))
		},
	}

	// Define the `--name` flag
	cmd.Flags().StringVarP(&moduleName, "name", "n", "", "Name of the module to generate (optional)")

	return cmd
}
