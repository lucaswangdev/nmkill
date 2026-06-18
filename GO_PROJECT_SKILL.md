# Go CLI 项目开发流程

> 基于 nmkill 项目的完整开发经验总结，适用于快速创建、开发和发布 Go CLI 工具

---

## 目录

1. [项目初始化](#1-项目初始化)
2. [代码结构](#2-代码结构)
3. [核心功能开发](#3-核心功能开发)
4. [测试](#4-测试)
5. [CLI 框架](#5-cli-框架)
6. [构建和发布](#6-构建和发布)
7. [GitHub Actions CI/CD](#7-github-actions-cicd)
8. [README 编写](#8-readme-编写)
9. [快速开始模板](#9-快速开始模板)

---

## 1. 项目初始化

### 创建项目目录
```bash
mkdir myproject && cd myproject
```

### 初始化 Go 模块
```bash
go mod init github.com/你的用户名/myproject
```

### 添加依赖
```bash
go get github.com/spf13/cobra
```

---

## 2. 代码结构

### 推荐目录结构
```
myproject/
├── main.go              # 入口文件
├── cmd/                 # 命令模块
│   ├── root.go         # 根命令
│   ├── query.go        # 子命令1
│   └── execute.go      # 子命令2
├── internal/            # 内部包
│   ├── scanner/        # 扫描逻辑
│   ├── csv/            # CSV处理
│   └── delete/         # 删除逻辑
├── config/             # 配置
├── .github/
│   └── workflows/
│       └── release.yml # 发布流程
├── README.md
├── CLAUDE.md           # 项目维护文档
└── go.mod
```

### main.go 模板
```go
package main

import (
    "fmt"
    "os"
    "github.com/你的用户名/myproject/cmd"
)

func main() {
    if err := cmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

---

## 3. 核心功能开发

### Scanner (目录扫描)
```go
package scanner

import (
    "os"
    "path/filepath"
)

type Result struct {
    Path string
    Size int64
}

type Scanner struct{}

func New() *Scanner {
    return &Scanner{}
}

func (s *Scanner) Scan(root string) ([]Result, error) {
    var results []Result
    filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return nil
        }
        if info.IsDir() && info.Name() == "target" {
            size, _ := s.calculateSize(path)
            results = append(results, Result{Path: path, Size: size})
            return filepath.SkipDir
        }
        return nil
    })
    return results, nil
}
```

---

## 4. 测试

### 编写测试
```go
package scanner

import "testing"

func TestScanner_Scan(t *testing.T) {
    s := New()
    results, err := s.Scan("/tmp/test")
    if err != nil {
        t.Fatalf("Scan failed: %v", err)
    }
    if len(results) != 1 {
        t.Errorf("Expected 1 result, got %d", len(results))
    }
}
```

### 运行测试
```bash
go test ./... -v
```

---

## 5. CLI 框架

### 使用 Cobra

#### cmd/root.go
```go
package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "myapp",
    Short: "My CLI application",
    Version: "0.1.0",
}

func Execute() error {
    return rootCmd.Execute()
}
```

#### cmd/subcommand.go
```go
package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var subCmd = &cobra.Command{
    Use:   "sub",
    Short: "Sub command",
    Args:  cobra.MaximumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) > 0 {
            fmt.Println(args[0])
        }
    },
}

func init() {
    rootCmd.AddCommand(subCmd)
    subCmd.Flags().BoolP("dry-run", "n", false, "Dry run")
}
```

---

## 6. 构建和发布

### 本地构建
```bash
# 普通构建
go build -o myapp .

# 去除调试信息 (更小)
go build -ldflags="-s -w" -o myapp .

# Windows exe
GOOS=windows go build -o myapp.exe .
```

### 跨平台构建矩阵
```bash
# Linux
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o myapp-linux-amd64 .

# macOS
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o myapp-darwin-amd64 .
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o myapp-darwin-arm64 .

# Windows
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o myapp-windows-amd64.exe .
```

---

## 7. GitHub Actions CI/CD

### .github/workflows/release.yml
```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        goos: [linux, darwin, windows]
        goarch: [amd64, arm64]
        include:
          - goos: windows
            extension: .exe

    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Build
        env:
          GOOS: ${{ matrix.goos }}
          GOARCH: ${{ matrix.goarch }}
          CGO_ENABLED: 0
        run: |
          ext=""
          [ "${{ matrix.goos }}" = "windows" ] && ext=".exe"
          filename="myapp-${{ matrix.goos }}-${{ matrix.goarch }}${ext}"
          go build -ldflags="-s -w" -o "${filename}" .

      - name: Upload Release
        uses: softprops/action-gh-release@v1
        with:
          files: myapp-*
```

### 发布流程
```bash
# 1. 提交代码
git add . && git commit -m "description"
git push origin main

# 2. 创建 tag
git tag v0.1.0
git push origin v0.1.0
```

---

## 8. README 编写

### 必须包含的内容
1. **项目简介** - 一句话说清楚
2. **特性列表** - 3-5 个核心功能
3. **安装方式** - 每平台单独命令
4. **使用方法** - 典型使用流程
5. **命令说明** - 所有 flag
6. **License**

### 安装命令格式
```bash
# Linux
curl -L https://github.com/user/repo/releases/download/latest/myapp-linux-amd64 -o /tmp/myapp && chmod +x /tmp/myapp && sudo mv /tmp/myapp /usr/local/bin/myapp

# macOS
curl -L https://github.com/user/repo/releases/download/latest/myapp-darwin-amd64 -o /tmp/myapp && chmod +x /tmp/myapp && sudo mv /tmp/myapp /usr/local/bin/myapp

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/user/repo/releases/download/latest/myapp-windows-amd64.exe" -OutFile "C:\Windows\System32\myapp.exe"
```

---

## 9. 快速开始模板

### 一键创建项目
```bash
# 创建项目
mkdir -p myproject/{cmd,internal,config,.github/workflows}
cd myproject

# 初始化
go mod init github.com/username/myproject
go get github.com/spf13/cobra

# 创建文件
touch main.go cmd/root.go cmd/sub.go

# 开发测试
go run main.go --help

# 构建
go build -o myapp .

# 发布
git add . && git commit -m "init"
git tag v0.1.0 && git push origin v0.1.0
```

---

## nmkill 项目关键决策记录

| 决策 | 选择 | 原因 |
|------|------|------|
| CLI框架 | cobra | 事实标准，功能完整 |
| CSV格式 | 路径、MB、yes/no | 用户友好，可编辑 |
| CSV存储 | ~/.appname/ | 符合XDG规范 |
| 删除确认 | 两次确认 | 防止误删 |
| 大小阈值 | 500MB | 合理默认值 |
| 发布方式 | GitHub Release | 官方、免费 |

---

## 常用命令速查

```bash
# 开发
go run main.go          # 运行
go build .              # 构建
go test ./...           # 测试

# 发布
git tag v0.x.x         # 打标签
git push origin v0.x.x # 推送标签

# 清理
go mod tidy            # 整理依赖
rm -f myapp            # 删除二进制
```

---

## 下次开发建议

1. **项目初始化时**立即创建 `CLAUDE.md`，记录设计决策
2. **先写测试**，再实现功能 (TDD)
3. **每完成一个功能**就 commit，不要最后一次性提交
4. **README 在第一个版本发布前**完成
5. **使用语义化版本号** vMAJOR.MINOR.PATCH

---

## ⚠️ 经验教训: 版本号管理

### ❌ 错误做法: 硬编码版本号

```go
// cmd/root.go
var rootCmd = &cobra.Command{
    Version: "0.1.0",  // ❌ 硬编码！每次发布都要手动改
}
```

**问题**:
- 发布新版本时忘记改这里
- GitHub Actions 构建时版本号不变
- 用户下载的是新版本，但显示的是旧版本号

### ✅ 正确做法: ldflags 注入版本号

**1. cmd/root.go** - 定义变量，不赋值
```go
var version = "dev" // 构建时注入
```

**2. release.yml** - 构建时注入
```yaml
- name: Extract version from tag
  id: version
  run: echo "VERSION=${GITHUB_REF#refs/tags/v}" >> $GITHUB_OUTPUT

- name: Build
  env:
    GOOS: ${{ matrix.goos }}
    GOARCH: ${{ matrix.goarch }}
    CGO_ENABLED: 0
  run: |
    ext=""
    [ "${{ matrix.goos }}" = "windows" ] && ext=".exe"
    filename="myapp-${{ matrix.goos }}-${{ matrix.goarch }}${ext}"
    # 关键：-X flag 注入版本号
    go build -ldflags="-s -w -X myproject/cmd.version=${{ steps.version.outputs.VERSION }}" -o "${filename}" .
```

**3. 本地测试**
```bash
go build -ldflags="-X myproject/cmd.version=0.1.0" -o myapp .
./myapp --version  # 输出: myapp version 0.1.0
```

### 🔑 关键点

| 项目 | 说明 |
|------|------|
| `-X` flag | Go 编译时注入变量值 |
| 格式 | `-X package.path.variable=value` |
| 示例 | `-X github.com/user/myproject/cmd.version=0.1.0` |
| 变量类型 | 只能注入 string 类型变量 |

### 发布流程检查清单

- [ ] `cmd/root.go` 中 version 变量默认值为 "dev"
- [ ] `release.yml` 使用 ldflags 注入版本号
- [ ] 打 tag 前本地测试: `go build -ldflags="-X..."`
- [ ] 验证: 下载 release 二进制，运行 `--version`
- [ ] README 更新为新版本号

### GitHub Actions 完整模板

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        goos: [linux, darwin, windows]
        goarch: [amd64, arm64]

    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Extract version
        id: version
        run: echo "VERSION=${GITHUB_REF#refs/tags/v}" >> $GITHUB_OUTPUT

      - name: Build
        env:
          GOOS: ${{ matrix.goos }}
          GOARCH: ${{ matrix.goarch }}
          CGO_ENABLED: 0
        run: |
          ext=""
          [ "${{ matrix.goos }}" = "windows" ] && ext=".exe"
          filename="myapp-${{ matrix.goos }}-${{ matrix.goarch }}${ext}"
          go build -ldflags="-s -w -X myproject/cmd.version=${{ steps.version.outputs.VERSION }}" -o "${filename}" .

      - name: Upload Release
        uses: softprops/action-gh-release@v1
        with:
          files: myapp-*
```

---

## 🔧 调试技巧

### 本地验证版本号注入
```bash
# 测试不同版本
go build -ldflags="-X myproject/cmd.version=1.0.0" -o myapp .
./myapp --version

# 验证注入成功
strings myapp | grep "version"
```

### 检查 release.yml 是否正确执行
```bash
# 查看 Actions 日志中的 version 输出
# 确认: VERSION=0.1.2
```
