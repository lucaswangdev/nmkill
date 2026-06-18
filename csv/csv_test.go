package csv

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/lucaswangdev/nmkill/scanner"
)

func TestWriteCSV(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "test.csv")
	
	modules := []scanner.NodeModule{
		{Path: "/path/to/node_modules1", Size: 600 * 1024 * 1024, LastMod: time.Now()},
		{Path: "/path/to/node_modules2", Size: 400 * 1024 * 1024, LastMod: time.Now()},
		{Path: "/path/to/node_modules3", Size: 500 * 1024 * 1024, LastMod: time.Now()},
	}
	
	err := WriteCSV(csvPath, modules)
	if err != nil {
		t.Fatalf("WriteCSV failed: %v", err)
	}
	
	// 读取验证
	data, err := os.ReadFile(csvPath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	
	content := string(data)
	
	// 验证表头包含 size(MB)
	if !strings.Contains(content, "size(MB)") {
		t.Error("CSV header should contain size(MB)")
	}
	
	// 验证内容包含 yes/no
	if !strings.Contains(content, "yes") {
		t.Error("Expected 'yes' in CSV")
	}
	if !strings.Contains(content, "no") {
		t.Error("Expected 'no' in CSV")
	}
	
	// 验证 MB 值正确
	if !strings.Contains(content, "600.00") {
		t.Error("Expected '600.00' MB in CSV")
	}
	if !strings.Contains(content, "400.00") {
		t.Error("Expected '400.00' MB in CSV")
	}
}

func TestReadCSV(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "test.csv")
	
	// 写入测试 CSV (MB格式)
	content := `path,size(MB),delete
/path/to/nm1,500.00,yes
/path/to/nm2,100.00,no
`
	err := os.WriteFile(csvPath, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}
	
	records, err := ReadCSV(csvPath)
	if err != nil {
		t.Fatalf("ReadCSV failed: %v", err)
	}
	
	if len(records) != 2 {
		t.Errorf("Expected 2 records, got %d", len(records))
	}
	
	if records[0].Path != "/path/to/nm1" {
		t.Errorf("Expected /path/to/nm1, got %s", records[0].Path)
	}
	// 500 MB = 524288000 bytes
	if records[0].Size != 524288000 {
		t.Errorf("Expected 524288000, got %d", records[0].Size)
	}
	if records[0].Delete != "yes" {
		t.Errorf("Expected yes, got %s", records[0].Delete)
	}
	
	if records[1].Delete != "no" {
		t.Errorf("Expected no, got %s", records[1].Delete)
	}
}

func TestGetFilesToDelete(t *testing.T) {
	records := []Record{
		{Path: "/path/to/nm1", Delete: "yes"},
		{Path: "/path/to/nm2", Delete: "no"},
		{Path: "/path/to/nm3", Delete: "yes"},
	}
	
	toDelete := GetFilesToDelete(records)
	
	if len(toDelete) != 2 {
		t.Errorf("Expected 2 files to delete, got %d", len(toDelete))
	}
	
	if toDelete[0] != "/path/to/nm1" {
		t.Errorf("Expected /path/to/nm1, got %s", toDelete[0])
	}
	if toDelete[1] != "/path/to/nm3" {
		t.Errorf("Expected /path/to/nm3, got %s", toDelete[1])
	}
}
