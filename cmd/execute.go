package cmd

import (
	"fmt"
	"os"

	"github.com/lucaswangdev/nmkill/config"
	"github.com/lucaswangdev/nmkill/delete"
	"github.com/spf13/cobra"
)

var executeCmd = &cobra.Command{
	Use:   "execute [csv]",
	Short: "根据 CSV 文件删除 node_modules",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var csvPath string
		
		if len(args) > 0 {
			// 用户指定了 CSV 文件路径
			csvPath = args[0]
		} else {
			// 使用最新的 CSV 文件
			csvPath = config.GetLatestCSV()
			if csvPath == "" {
				fmt.Fprintln(os.Stderr, "错误: 未找到 CSV 文件")
				fmt.Fprintln(os.Stderr, "请先运行 'nmkill query' 命令")
				os.Exit(1)
			}
			fmt.Printf("使用最新 CSV 文件: %s\n", csvPath)
		}

		// 检查文件是否存在
		if _, err := os.Stat(csvPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "错误: CSV 文件不存在: %s\n", csvPath)
			fmt.Fprintf(os.Stderr, "请先运行 'nmkill query' 命令\n")
			os.Exit(1)
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		skipConfirm, _ := cmd.Flags().GetBool("yes")

		if !dryRun && !skipConfirm {
			fmt.Println("警告: 即将删除 CSV 文件中标记为 'yes' 的 node_modules 目录!")
			fmt.Print("确认删除? (输入 'yes' 继续): ")
			var confirm string
			fmt.Scanln(&confirm)
			if confirm != "yes" {
				fmt.Println("操作已取消")
				os.Exit(0)
			}
		}

		d := delete.New(dryRun)
		success, fail, err := d.DeleteFromCSV(csvPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: 删除失败: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\n删除完成: 成功 %d, 失败 %d\n", success, fail)
	},
}

func init() {
	rootCmd.AddCommand(executeCmd)
	executeCmd.Flags().BoolP("dry-run", "n", false, "模拟运行，不实际删除")
	executeCmd.Flags().BoolP("yes", "y", false, "跳过确认直接删除")
}
