package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/lucaswangdev/nmkill/scanner"
)

const (
	HeaderPath   = "path"
	HeaderSize   = "size"
	HeaderDelete = "delete"
	DeleteYes    = "yes"
	DeleteNo     = "no"
	DefaultSizeMB = 500 // 默认 500MB 以上标记为 yes
)

// Record CSV 记录
type Record struct {
	Path   string
	Size   int64
	Delete string
}

// WriteCSV 写入 CSV 文件
func WriteCSV(path string, modules []scanner.NodeModule) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	if err := writer.Write([]string{HeaderPath, HeaderSize, HeaderDelete}); err != nil {
		return err
	}

	// 写入数据
	for _, m := range modules {
		sizeMB := m.Size / (1024 * 1024)
		deleteMark := DeleteNo
		if sizeMB >= DefaultSizeMB {
			deleteMark = DeleteYes
		}

		if err := writer.Write([]string{
			m.Path,
			strconv.FormatInt(m.Size, 10),
			deleteMark,
		}); err != nil {
			return err
		}
	}

	return nil
}

// ReadCSV 读取 CSV 文件
func ReadCSV(path string) ([]Record, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("读取 CSV 失败: %w", err)
	}

	var result []Record
	for i, record := range records {
		// 跳过表头
		if i == 0 {
			continue
		}
		if len(record) < 3 {
			continue
		}

		size, _ := strconv.ParseInt(record[1], 10, 64)
		result = append(result, Record{
			Path:   record[0],
			Size:   size,
			Delete: record[2],
		})
	}

	return result, nil
}

// GetFilesToDelete 获取需要删除的文件列表
func GetFilesToDelete(records []Record) []string {
	var toDelete []string
	for _, r := range records {
		if r.Delete == DeleteYes {
			toDelete = append(toDelete, r.Path)
		}
	}
	return toDelete
}
