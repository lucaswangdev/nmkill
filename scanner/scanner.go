package scanner

import (
	"os"
	"path/filepath"
	"sync"
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
	mu          sync.Mutex
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
	s.results = make([]NodeModule, 0)
	s.dirsVisited = 0

	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var entryErr error

	// 并行扫描每个子目录
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// 检查是否是 node_modules
		if entry.Name() == "node_modules" {
			path := filepath.Join(root, entry.Name())
			wg.Add(1)
			s.scanNodeModules(path, &wg, progressCallback)
			continue
		}

		path := filepath.Join(root, entry.Name())
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			s.walkDir(p, progressCallback)
		}(path)
	}

	wg.Wait()

	// 最终回调
	if len(progressCallback) > 0 && s.dirsVisited > 0 {
		progressCallback[0](s.dirsVisited)
	}

	return s.results, entryErr
}

// walkDir 遍历目录查找 node_modules
func (s *Scanner) walkDir(root string, progressCallback []ProgressCallback) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		s.mu.Lock()
		s.dirsVisited++
		dirsVisited := s.dirsVisited
		s.mu.Unlock()

		if len(progressCallback) > 0 && dirsVisited%1000 == 0 {
			progressCallback[0](dirsVisited)
		}

		if entry.Name() == "node_modules" {
			path := filepath.Join(root, entry.Name())
			s.scanNodeModules(path, nil, progressCallback)
			continue
		}

		path := filepath.Join(root, entry.Name())
		s.walkDir(path, progressCallback)
	}
}

// scanNodeModules 扫描单个 node_modules 并添加到结果
func (s *Scanner) scanNodeModules(path string, wg *sync.WaitGroup, progressCallback []ProgressCallback) {
	if wg != nil {
		defer wg.Done()
	}

	info, err := os.Stat(path)
	if err != nil {
		return
	}

	size, _ := s.calculateDirSize(path)

	s.mu.Lock()
	s.results = append(s.results, NodeModule{
		Path:    path,
		Size:    size,
		LastMod: info.ModTime(),
	})
	s.mu.Unlock()
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
