package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd 是 CLI 入口
var RootCmd = &cobra.Command{
	Use:   "tokencli",
	Short: "Token 备份 & 恢复 CLI",
	Long:  `一个管理 Token 备份、恢复、自动清理的命令行工具`,
}

// Execute 运行 CLI
func Execute() error {
	return RootCmd.Execute()
}
