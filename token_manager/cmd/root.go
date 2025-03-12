package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"token_manager/tokenmanager"
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

var backupCmd = &cobra.Command{
	Use:   "backup [token_id]",
	Short: "备份 Token",
	Run: func(cmd *cobra.Command, args []string) {
		tm := tokenmanager.NewTokenManager("v2")
		tokenID := args[0]
		err := tm.BackupToken(tokenID)
		if err != nil {
			return
		}
		fmt.Println("Backup successful! File saved at:")
	},
}

var restoreCmd = &cobra.Command{
	Use:   "restore [token_id]",
	Short: "恢复 Token 备份",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tm := tokenmanager.NewTokenManager("v2")
		tokenID := args[0]
		err := tm.RestoreToken(tokenID)
		if err != nil {
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(backupCmd)
	RootCmd.AddCommand(restoreCmd)
}
