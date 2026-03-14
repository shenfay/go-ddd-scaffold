package command

import (
	"github.com/shenfay/go-ddd-scaffold/cmd/cli/internal/generators"
	"github.com/spf13/cobra"
)

// RegisterAll 注册所有子命令到根命令
func RegisterAll(rootCmd *cobra.Command) {
	rootCmd.AddCommand(initCmd())
	rootCmd.AddCommand(generateCmd())
	rootCmd.AddCommand(migrateCmd())
	rootCmd.AddCommand(docsCmd())
	rootCmd.AddCommand(cleanCmd())
	rootCmd.AddCommand(versionCmd())
}

// versionCmd 显示版本信息
func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("go-ddd-scaffold version %s\n", cmd.Root().Version)
		},
	}
}

// cleanCmd 清理生成的文件
func cleanCmd() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "clean [path]",
		Short: "Clean generated files",
		Long:  `Remove generated code files from the project`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			targetPath := "."
			if len(args) > 0 {
				targetPath = args[0]
			}

			return generators.CleanGeneratedFiles(targetPath, dryRun)
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be deleted")
	return cmd
}
