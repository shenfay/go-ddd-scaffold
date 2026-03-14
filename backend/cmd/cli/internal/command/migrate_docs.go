package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func migrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Database migration commands",
		Long:  `Manage database migrations`,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "up",
		Short: "Run all pending migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Running migrations...")
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "down",
		Short: "Rollback last migration",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Rolling back migration...")
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "create [name]",
		Short: "Create a new migration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Creating migration: %s\n", args[0])
			return nil
		},
	})

	return cmd
}

func docsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docs",
		Short: "Generate documentation",
		Long:  `Generate API and project documentation`,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "swagger",
		Short: "Generate Swagger documentation",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Generating Swagger docs...")
			return nil
		},
	})

	return cmd
}
