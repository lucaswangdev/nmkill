package delete

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lucaswangdev/nmkill/csv"
)

// Deleter 删除器
type Deleter struct {
	DryRun bool
}

// New 创建删除器
func New(dryRun bool) *Deleter {
	return &Deleter{DryRun: dryRun}
}

// DeleteFromCSV 从 CSV 文件删除标记为 yes 的目录
func (d *Deleter) DeleteFromCSV(csvPath string) (int, int, error) {
	// 读取 CSV
	records, err := csv.ReadCSV(csvPath)
	if err != nil {
		return 0, 0, err
	}

	toDelete := csv.GetFilesToDelete(records)
	return d.DeletePaths(toDelete)
}

// DeletePaths 删除指定的路径
func (d *Deleter) DeletePaths(paths []string) (int, int, error) {
	successCount := 0
	failCount := 0

	for _, path := range paths {
		if d.DryRun {
			fmt.Printf("[DRY-RUN] 将删除: %s\n", path)
			successCount++
			continue
		}

		// 检查路径是否存在
		info, err := os.Stat(path)
		if err != nil {
			fmt.Printf("[跳过] 路径不存在: %s\n", path)
			failCount++
			continue
		}

		if !info.IsDir() {
			fmt.Printf("[跳过] 不是目录: %s\n", path)
			failCount++
			continue
		}

		// 确认是 node_modules
		if filepath.Base(path) != "node_modules" {
			fmt.Printf("[跳过] 不是 node_modules: %s\n", path)
			failCount++
			continue
		}

		// 删除目录
		if err := os.RemoveAll(path); err != nil {
			fmt.Printf("[失败] %s: %v\n", path, err)
			failCount++
		} else {
			fmt.Printf("[删除成功] %s\n", path)
			successCount++
		}
	}

	return successCount, failCount, nil
}
