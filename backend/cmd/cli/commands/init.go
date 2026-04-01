package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/shenfay/go-ddd-scaffold/cmd/cli/generators"
	"github.com/spf13/cobra"
)

func initCmd() *cobra.Command {
	var opts generators.InitOptions

	cmd := &cobra.Command{
		Use:   "init [project_name]",
		Short: "Initialize a new DDD project",
		Long:  `Create a new DDD project with standard Clean Architecture structure`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProjectName = args[0]

			// 从标志或配置中获取其他选项
			if !cmd.Flags().Changed("module-path") {
				opts.ModulePath = fmt.Sprintf("github.com/%s/%s", getCurrentUser(), opts.ProjectName)
			}

			generator := generators.NewInitGenerator(opts)
			return generator.Generate()
		},
	}

	cmd.Flags().StringVarP(&opts.ModulePath, "module-path", "m", "", "Go module path")
	cmd.Flags().StringVarP(&opts.Author, "author", "a", "", "Author name")
	cmd.Flags().StringVarP(&opts.Email, "email", "e", "", "Author email")
	cmd.Flags().StringVarP(&opts.License, "license", "l", "MIT", "License type")
	cmd.Flags().StringVarP(&opts.Template, "template", "t", "clean-architecture", "Project template")
	cmd.Flags().BoolVarP(&opts.SkipFrontend, "skip-frontend", "", false, "Skip frontend initialization")
	cmd.Flags().BoolVarP(&opts.WithDocker, "with-docker", "", false, "Include Docker configuration")
	cmd.Flags().BoolVarP(&opts.WithK8s, "with-k8s", "", false, "Include Kubernetes manifests")

	return cmd
}

func getCurrentUser() string {
	// 1. 尝试从环境变量获取
	if username := os.Getenv("USER"); username != "" {
		return username
	}

	// 2. 尝试从 git config 获取
	if output, err := exec.Command("git", "config", "user.name").Output(); err == nil {
		return strings.TrimSpace(string(output))
	}

	// 3. 回退到默认值
	return "unknown"
}
