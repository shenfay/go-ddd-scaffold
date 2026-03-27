package main

import (
	"fmt"
	"os"

	command "github.com/shenfay/go-ddd-scaffold/cmd/cli/commands"
	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"
	commit  = "dev"
	date    = "unknown"
)

func main() {
	rootCmd := newRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "go-ddd-scaffold",
		Short: "Go DDD Scaffold CLI Tool",
		Long: `A CLI tool for generating DDD scaffold code following best practices.
	
Features:
  - Project initialization with Clean Architecture
  - Code generation (entities, repositories, DAOs, services)
  - Database migration management
  - Documentation generation
  - Custom template support`,
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// 设置版本信息
	cmd.SetVersionTemplate(fmt.Sprintf("Version: %s\nCommit: %s\nBuilt: %s\n", version, commit, date))

	// 注册所有子命令
	command.RegisterAll(cmd)

	// 全局标志
	cmd.PersistentFlags().StringP("config", "c", "", "Config file path (default is $HOME/.go-ddd-scaffold.yaml)")
	cmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
	cmd.PersistentFlags().BoolP("dry-run", "n", false, "Show what would be done without actually doing it")

	return cmd
}
