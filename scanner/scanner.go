package scanner

import (
	"os"
	"path/filepath"
	"time"
)

// NodeModule node_modules 目录信息
type NodeModule struct {
	Path    string
	Size    int64
	LastMod time.Time
}

// Scanner 目录扫描器
type Scanner struct {
	results []NodeModule
}

// New 创建新的扫描器
func New() *Scanner {
	return &Scanner{
		results: make([]NodeModule, 0),
	}
}

// Scan 从指定根目录开始扫描 node_modules
func (s *Scanner) Scan(root string) ([]NodeModule, error) {
	results := make([]NodeModule, 0)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 跳过无法访问的目录
		}

		// 如果是目录且名称为 node_modules
		if info.IsDir() && info.Name() == "node_modules" {
			// 计算目录大小
			size, err := s.calculateDirSize(path)
			if err != nil {
				return nil
			}

			// 获取最后修改时间
			lastMod := info.ModTime()

			results = append(results, NodeModule{
				Path:    path,
				Size:    size,
				LastMod: lastMod,
			})

			// 跳过 node_modules 内部遍历
			return filepath.SkipDir
		}

		return nil
	})

	return results, err
}

// calculateDirSize 计算目录大小
func (s *Scanner) calculateDirSize(path string) (int64, error) {
	var size int64

	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size, err
}
