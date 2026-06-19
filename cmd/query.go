package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lucaswangdev/nmkill/config"
	"github.com/lucaswangdev/nmkill/csv"
	"github.com/lucaswangdev/nmkill/scanner"
	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   "query [path]",
	Short: "扫描目录下的 node_modules",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		root := "."
		if len(args) > 0 {
			root = args[0]
		}

		absRoot, err := filepath.Abs(root)
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: 无法获取路径: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("正在扫描: %s\n", absRoot)

		s := scanner.New()
		progressChan := make(chan int64)
		var totalDirs int64

		// 启动进度显示 goroutine
		go func() {
			spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
			var lastCount int64
			for {
				select {
				case count, ok := <-progressChan:
					if !ok {
						return
					}
					lastCount = count
					idx := (int(count / 1000)) % len(spinners)
					fmt.Printf("\r%s 已扫描 %d 个目录", spinners[idx], count)
				case <-time.After(100 * time.Millisecond):
					if lastCount > 0 {
						fmt.Printf("\r已扫描 %d 个目录", lastCount)
					}
				}
			}
		}()

		results, err := s.Scan(absRoot, func(dirsVisited int64) {
			totalDirs = dirsVisited
			select {
			case progressChan <- dirsVisited:
			default:
			}
		})

		close(progressChan)

		// 清除进度行并输出结果
		fmt.Printf("\r%s\r", strings.Repeat(" ", 60))
		fmt.Printf("共扫描 %d 个目录\n", totalDirs)

		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: 扫描失败: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("找到 %d 个 node_modules 目录\n", len(results))

		if len(results) == 0 {
			fmt.Println("没有找到 node_modules 目录")
			return
		}

		// 生成带时间戳的 CSV 文件
		outputPath := config.GenerateCSVPath()

		if err := csv.WriteCSV(outputPath, results); err != nil {
			fmt.Fprintf(os.Stderr, "错误: 写入 CSV 失败: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("结果已保存到: %s\n", outputPath)
		fmt.Println("\n说明: delete 字段为 'yes' 表示将被删除 (>= 500MB)")
		fmt.Println("      delete 字段为 'no' 表示不删除 (< 500MB)")
		fmt.Println("      你可以手动编辑 CSV 文件来修改 delete 字段")
		fmt.Printf("\n执行删除: nmkill execute %s\n", outputPath)
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
}
