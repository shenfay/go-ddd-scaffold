package command

import (
	"github.com/shenfay/go-ddd-scaffold/cmd/cli/generators"
	"github.com/spf13/cobra"
)

func generateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Short:   "Generate code for DDD scaffold",
		Long:    `Generate various types of code following DDD and Clean Architecture patterns`,
		Aliases: []string{"gen", "g"},
	}

	cmd.AddCommand(generateEntityCmd())
	cmd.AddCommand(generateDAOCmd())
	cmd.AddCommand(generateRepositoryCmd())
	cmd.AddCommand(generateServiceCmd())
	cmd.AddCommand(generateHandlerCmd())
	cmd.AddCommand(generateDTOCmd())

	return cmd
}

func generateEntityCmd() *cobra.Command {
	var opts generators.EntityOptions

	cmd := &cobra.Command{
		Use:   "entity [name]",
		Short: "Generate domain entity",
		Long:  `Generate domain entity with value objects and aggregate root`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]
			generator := generators.NewEntityGenerator(opts)
			return generator.Generate()
		},
	}

	cmd.Flags().StringVarP(&opts.Fields, "fields", "f", "", "Field definitions (format: name:type,name:type)")
	cmd.Flags().StringVarP(&opts.Methods, "methods", "m", "", "Business methods to generate")
	cmd.Flags().StringVarP(&opts.Package, "package", "p", "", "Domain package name")
	cmd.Flags().BoolVarP(&opts.WithVO, "with-vo", "", false, "Generate value objects")
	cmd.Flags().BoolVarP(&opts.WithAggregate, "with-aggregate", "", false, "Generate aggregate root")

	return cmd
}

func generateDAOCmd() *cobra.Command {
	opts := generators.DAOOptions{}

	cmd := &cobra.Command{
		Use:   "dao",
		Short: "Generate DAO layer from database",
		Long:  `Generate DAO layer code and models by reverse engineering existing database tables using gorm/gen`,
		RunE: func(cmd *cobra.Command, args []string) error {
			generator := generators.NewDAOGenerator(opts)
			return generator.Generate()
		},
	}

	cmd.Flags().StringVarP(&opts.OutputPath, "output", "o", "", "Output directory for generated code")
	cmd.Flags().StringVarP(&opts.DSN, "dsn", "d", "", "Database connection DSN (will use config if not provided)")
	cmd.Flags().StringVarP(&opts.ConfigFile, "config", "c", "", "Config file path")
	cmd.Flags().BoolVarP(&opts.WithUnitTest, "with-test", "", false, "Generate unit tests")
	cmd.Flags().BoolVarP(&opts.FieldNullable, "field-nullable", "", true, "Generate pointer for nullable fields")
	cmd.Flags().StringSliceVarP(&opts.Tables, "tables", "t", []string{}, "Specific tables to generate (comma-separated)")

	return cmd
}

func generateRepositoryCmd() *cobra.Command {
	var opts generators.RepositoryOptions

	cmd := &cobra.Command{
		Use:   "repository [name]",
		Short: "Generate repository layer",
		Long:  `Generate repository implementation that uses DAO`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]
			generator := generators.NewRepositoryGenerator(opts)
			return generator.Generate()
		},
	}

	cmd.Flags().StringVarP(&opts.Domain, "domain", "d", "", "Domain name")
	cmd.Flags().StringVarP(&opts.OutputDir, "output", "o", "", "Output directory")

	return cmd
}

func generateServiceCmd() *cobra.Command {
	var opts generators.ServiceOptions

	cmd := &cobra.Command{
		Use:   "service [name]",
		Short: "Generate application service",
		Long:  `Generate application service with CQRS pattern`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]
			generator := generators.NewServiceGenerator(opts)
			return generator.Generate()
		},
	}

	cmd.Flags().StringVarP(&opts.Type, "type", "t", "application", "Service type (application/domain)")
	cmd.Flags().StringVarP(&opts.Methods, "methods", "m", "", "Service methods")
	cmd.Flags().StringSliceVarP(&opts.Dependencies, "deps", "d", []string{}, "Service dependencies")

	return cmd
}

func generateHandlerCmd() *cobra.Command {
	var opts generators.HandlerOptions

	cmd := &cobra.Command{
		Use:   "handler [name]",
		Short: "Generate command/query handler",
		Long:  `Generate CQRS command or query handler`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]
			generator := generators.NewHandlerGenerator(opts)
			return generator.Generate()
		},
	}

	cmd.Flags().StringVarP(&opts.Type, "type", "t", "command", "Handler type (command/query)")
	cmd.Flags().StringVarP(&opts.Domain, "domain", "d", "", "Domain name")

	return cmd
}

func generateDTOCmd() *cobra.Command {
	var opts generators.DTOOptions

	cmd := &cobra.Command{
		Use:   "dto [name]",
		Short: "Generate DTO (Data Transfer Object)",
		Long:  `Generate DTO for API layer`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]
			generator := generators.NewDTOGenerator(opts)
			return generator.Generate()
		},
	}

	cmd.Flags().StringVarP(&opts.Type, "type", "t", "request", "DTO type (request/response)")
	cmd.Flags().StringVarP(&opts.Fields, "fields", "f", "", "Field definitions")
	cmd.Flags().BoolVarP(&opts.WithValidation, "with-validation", "", true, "Add validation tags")

	return cmd
}
