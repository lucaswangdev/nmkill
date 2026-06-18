package cmd

import (
	"fmt"
	"os"

	"github.com/lucaswangdev/nmkill/config"
	"github.com/spf13/cobra"
)

var version = "dev" // 默认版本，构建时通过 ldflags 注入

var rootCmd = &cobra.Command{
	Use:   "nmkill",
	Short: "nmkill - 清理 node_modules 的 CLI 工具",
	Long: `nmkill 是一个用于查找和删除 node_modules 目录的 CLI 工具。
它可以帮助你释放磁盘空间。

使用方法:
  nmkill query [path]    扫描目录下的 node_modules
  nmkill execute [csv]   根据 CSV 文件删除 node_modules
  nmkill help            显示帮助信息`,
	Version: version,
}

func Execute() error {
	if err := config.EnsureConfigDir(); err != nil {
		fmt.Fprintf(os.Stderr, "警告: 无法创建配置目录: %v\n", err)
	}
	return rootCmd.Execute()
}
