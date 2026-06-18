package delete

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lucaswangdev/nmkill/csv"
)

func TestDeleter_DeletePaths_DryRun(t *testing.T) {
	tmpDir := t.TempDir()
	nmPath := filepath.Join(tmpDir, "node_modules")
	
	if err := os.MkdirAll(nmPath, 0755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(nmPath, "package.json"), []byte(`{}`), 0644)
	
	d := New(true) // dry-run
	
	success, fail, err := d.DeletePaths([]string{nmPath})
	if err != nil {
		t.Fatalf("DeletePaths failed: %v", err)
	}
	
	if success != 1 {
		t.Errorf("Expected 1 success, got %d", success)
	}
	if fail != 0 {
		t.Errorf("Expected 0 failures, got %d", fail)
	}
	
	// 验证目录仍然存在 (dry-run 不会删除)
	if _, err := os.Stat(nmPath); os.IsNotExist(err) {
		t.Error("Directory should still exist after dry-run")
	}
}

func TestDeleter_DeletePaths_Real(t *testing.T) {
	tmpDir := t.TempDir()
	nmPath := filepath.Join(tmpDir, "node_modules")
	
	if err := os.MkdirAll(nmPath, 0755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(nmPath, "package.json"), []byte(`{}`), 0644)
	
	d := New(false) // real delete
	
	success, fail, err := d.DeletePaths([]string{nmPath})
	if err != nil {
		t.Fatalf("DeletePaths failed: %v", err)
	}
	
	if success != 1 {
		t.Errorf("Expected 1 success, got %d", success)
	}
	if fail != 0 {
		t.Errorf("Expected 0 failures, got %d", fail)
	}
	
	// 验证目录已被删除
	if _, err := os.Stat(nmPath); !os.IsNotExist(err) {
		t.Error("Directory should be deleted")
	}
}

func TestDeleter_DeletePaths_NonExistent_DryRun(t *testing.T) {
	d := New(true) // dry-run
	
	// Dry-run 模式下，非存在的路径也被视为"成功"（只打印信息）
	success, fail, err := d.DeletePaths([]string{"/path/does/not/exist"})
	if err != nil {
		t.Fatalf("DeletePaths failed: %v", err)
	}
	
	if success != 1 {
		t.Errorf("Expected 1 success (dry-run), got %d", success)
	}
	if fail != 0 {
		t.Errorf("Expected 0 failures, got %d", fail)
	}
}

func TestDeleter_DeletePaths_NonExistent_Real(t *testing.T) {
	d := New(false) // real mode
	
	success, fail, err := d.DeletePaths([]string{"/path/does/not/exist"})
	if err != nil {
		t.Fatalf("DeletePaths failed: %v", err)
	}
	
	if success != 0 {
		t.Errorf("Expected 0 success, got %d", success)
	}
	if fail != 1 {
		t.Errorf("Expected 1 failure, got %d", fail)
	}
}

func TestDeleter_DeletePaths_NotNodeModules(t *testing.T) {
	tmpDir := t.TempDir()
	regularDir := filepath.Join(tmpDir, "regular_dir")
	
	if err := os.MkdirAll(regularDir, 0755); err != nil {
		t.Fatal(err)
	}
	
	d := New(false)
	
	success, fail, err := d.DeletePaths([]string{regularDir})
	if err != nil {
		t.Fatalf("DeletePaths failed: %v", err)
	}
	
	if success != 0 {
		t.Errorf("Expected 0 success, got %d", success)
	}
	if fail != 1 {
		t.Errorf("Expected 1 failure, got %d", fail)
	}
	
	// 验证目录仍然存在
	if _, err := os.Stat(regularDir); os.IsNotExist(err) {
		t.Error("Directory should still exist")
	}
}

func TestDeleter_DeleteFromCSV(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "test.csv")
	nmPath := filepath.Join(tmpDir, "node_modules")
	
	// 创建 CSV
	csvContent := "path,size,delete\n" + nmPath + ",100,yes\n"
	if err := os.WriteFile(csvPath, []byte(csvContent), 0644); err != nil {
		t.Fatal(err)
	}
	
	// 创建 node_modules 目录
	if err := os.MkdirAll(nmPath, 0755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(nmPath, "package.json"), []byte(`{}`), 0644)
	
	d := New(false)
	success, fail, err := d.DeleteFromCSV(csvPath)
	if err != nil {
		t.Fatalf("DeleteFromCSV failed: %v", err)
	}
	
	if success != 1 {
		t.Errorf("Expected 1 success, got %d", success)
	}
	if fail != 0 {
		t.Errorf("Expected 0 failures, got %d", fail)
	}
	
	// 验证目录已被删除
	if _, err := os.Stat(nmPath); !os.IsNotExist(err) {
		t.Error("Directory should be deleted")
	}
}

func init() {
	// 确保 csv 包已导入
	_ = csv.DeleteYes
}
