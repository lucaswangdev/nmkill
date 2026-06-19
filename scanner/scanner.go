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
	results     []NodeModule
	dirsVisited int64
}

// New 创建新的扫描器
func New() *Scanner {
	return &Scanner{
		results: make([]NodeModule, 0),
	}
}

// ProgressCallback 进度回调函数类型
type ProgressCallback func(dirsVisited int64)

// Scan 从指定根目录开始扫描 node_modules
// progressCallback 可选回调，每隔一段时间报告已扫描的目录数
func (s *Scanner) Scan(root string, progressCallback ...ProgressCallback) ([]NodeModule, error) {
	results := make([]NodeModule, 0)
	s.dirsVisited = 0

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		s.dirsVisited++

		// 每1000个目录触发一次回调
		if len(progressCallback) > 0 && s.dirsVisited%1000 == 0 {
			progressCallback[0](s.dirsVisited)
		}

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

	// 最终回调
	if len(progressCallback) > 0 && s.dirsVisited > 0 {
		progressCallback[0](s.dirsVisited)
	}

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
