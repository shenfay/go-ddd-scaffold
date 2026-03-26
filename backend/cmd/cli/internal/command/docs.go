package command

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func docsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docs",
		Short: "Generate documentation",
		Long:  `Generate API and project documentation`,
	}

	cmd.AddCommand(swaggerCmd())

	return cmd
}

// swaggerCmd 生成 Swagger API 文档
func swaggerCmd() *cobra.Command {
	var dir string
	var output string

	cmd := &cobra.Command{
		Use:   "swagger",
		Short: "Generate Swagger documentation",
		Long:  `Generate Swagger API documentation from code comments`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateSwagger(dir, output)
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", ".", "项目根目录")
	cmd.Flags().StringVarP(&output, "output", "o", "./docs/swagger", "输出目录")

	return cmd
}

// generateSwagger 执行 swag 命令生成文档
func generateSwagger(dir, output string) error {
	// 检查是否安装了 swag 工具
	if _, err := exec.LookPath("swag"); err != nil {
		fmt.Println("未找到 swag 工具，请先安装：")
		fmt.Println("  go install github.com/swaggo/swag/cmd/swag@latest")
		return fmt.Errorf("swag tool not found")
	}

	// 确保输出目录存在
	if err := os.MkdirAll(output, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 执行 swag init 命令
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	fmt.Printf("正在生成 Swagger 文档...\n")
	fmt.Printf("项目目录：%s\n", absDir)
	fmt.Printf("输出目录：%s\n", output)

	swagCmd := exec.Command("swag", "init",
		"-g", "cmd/api/main.go",
		"-d", absDir,
		"-o", output,
		"--parseDependency",
		"--parseInternal",
	)

	swagCmd.Stdout = os.Stdout
	swagCmd.Stderr = os.Stderr

	if err := swagCmd.Run(); err != nil {
		return fmt.Errorf("swag init failed: %w", err)
	}

	fmt.Println("✓ Swagger 文档生成成功！")
	fmt.Printf("📄 查看文档：%s/docs.go\n", output)
	fmt.Printf("🌐 访问 UI: http://localhost:8080/swagger/index.html\n")

	return nil
}
