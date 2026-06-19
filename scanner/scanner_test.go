package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkScan(b *testing.B) {
	// 创建大量测试目录模拟真实场景
	tmpDir := b.TempDir()

	// 创建多个项目目录
	for i := 0; i < 50; i++ {
		projectDir := filepath.Join(tmpDir, "project"+string(rune('A'+i)))
		os.MkdirAll(filepath.Join(projectDir, "node_modules", "pkg1"), 0755)
		os.MkdirAll(filepath.Join(projectDir, "node_modules", "pkg2"), 0755)
		os.MkdirAll(filepath.Join(projectDir, "src"), 0755)
		os.WriteFile(filepath.Join(projectDir, "package.json"), []byte(`{"name":"test"}`), 0644)

		// 添加一些文件到 node_modules
		os.WriteFile(filepath.Join(projectDir, "node_modules", "pkg1", "index.js"), []byte(`console.log("test")`), 0644)
		os.WriteFile(filepath.Join(projectDir, "node_modules", "pkg2", "index.js"), []byte(`console.log("test")`), 0644)
	}

	s := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Scan(tmpDir)
	}
}

func TestScanner_Scan(t *testing.T) {
	// 创建临时测试目录结构
	tmpDir := t.TempDir()
	
	// 创建 node_modules 目录
	nm1 := filepath.Join(tmpDir, "project1", "node_modules")
	nm2 := filepath.Join(tmpDir, "project2", "node_modules")
	
	if err := os.MkdirAll(nm1, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(nm2, 0755); err != nil {
		t.Fatal(err)
	}
	
	// 在 node_modules 中创建一些文件
	os.WriteFile(filepath.Join(nm1, "package.json"), []byte(`{"name":"test"}`), 0644)
	os.WriteFile(filepath.Join(nm2, "package.json"), []byte(`{"name":"test"}`), 0644)
	
	s := New()
	results, err := s.Scan(tmpDir)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}
	
	if len(results) != 2 {
		t.Errorf("Expected 2 node_modules, got %d", len(results))
	}
	
	// 验证路径
	paths := make(map[string]bool)
	for _, r := range results {
		paths[r.Path] = true
		if filepath.Base(r.Path) != "node_modules" {
			t.Errorf("Expected path ending with node_modules, got %s", r.Path)
		}
	}
	
	if !paths[nm1] {
		t.Error("Expected to find nm1")
	}
	if !paths[nm2] {
		t.Error("Expected to find nm2")
	}
}

func TestScanner_Scan_NoResults(t *testing.T) {
	tmpDir := t.TempDir()
	
	s := New()
	results, err := s.Scan(tmpDir)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}
	
	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

func TestScanner_Scan_SkipsNestedNodeModules(t *testing.T) {
	tmpDir := t.TempDir()
	
	// 创建嵌套的 node_modules
	nm1 := filepath.Join(tmpDir, "node_modules")
	nm2 := filepath.Join(nm1, "some-package", "node_modules")
	
	if err := os.MkdirAll(nm1, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(nm2, 0755); err != nil {
		t.Fatal(err)
	}
	
	s := New()
	results, err := s.Scan(tmpDir)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}
	
	// 应该只找到顶层的 node_modules
	if len(results) != 1 {
		t.Errorf("Expected 1 node_modules, got %d", len(results))
	}
	
	if results[0].Path != nm1 {
		t.Errorf("Expected %s, got %s", nm1, results[0].Path)
	}
}
