package config

import (
	"os"
	"path/filepath"
	"sort"
	"time"
)

// AppName 应用名称
var AppName = "nmkill"

// ConfigDir 配置目录 ~/.nmkill
var ConfigDir = filepath.Join(os.Getenv("HOME"), "." + AppName)

// EnsureConfigDir 确保配置目录存在
func EnsureConfigDir() error {
	return os.MkdirAll(ConfigDir, 0755)
}

// GenerateCSVPath 生成带时间戳的 CSV 文件路径
func GenerateCSVPath() string {
	timestamp := time.Now().Format("20060102_150405")
	return filepath.Join(ConfigDir, "node_modules_"+timestamp+".csv")
}

// GetCSVPath 返回默认 CSV 文件路径（最新生成的文件）
func GetCSVPath() string {
	matches, err := filepath.Glob(filepath.Join(ConfigDir, "node_modules_*.csv"))
	if err != nil || len(matches) == 0 {
		return ""
	}
	
	// 按文件名排序（时间戳在文件名中，排序后最后一个是最新的）
	sort.Strings(matches)
	return matches[len(matches)-1]
}

// GetLatestCSV 获取最新的 CSV 文件路径，如果没有则返回空字符串
func GetLatestCSV() string {
	return GetCSVPath()
}
