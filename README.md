# nmkill

> A CLI tool to find and clean `node_modules` directories, freeing up disk space.

## Features

- 🚀 **Fast scanning** - Multi-threaded directory traversal
- 📊 **CSV export** - Generate timestamped reports
- 🔒 **Safe deletion** - Dry-run mode and confirmation prompts
- ✏️ **Editable** - Manually edit CSV to choose what to delete
- 🌐 **Cross-platform** - macOS, Linux, Windows supported

## Installation

### macOS / Linux

```bash
# Using Go
go install github.com/lucaswangdev/nmkill@latest

# Or download from GitHub Releases
curl -L https://github.com/lucaswangdev/nmkill/releases/latest/download/nmkill-linux-amd64 -o nmkill
chmod +x nmkill
```

### Windows

```powershell
# Using Go
go install github.com/lucaswangdev/nmkill@latest

# Or download from GitHub Releases
# Download nmkill-windows-amd64.exe from the releases page
```

## Usage

### 1. Scan for node_modules

```bash
nmkill query ~/projects
```

Output:
```
正在扫描: /home/user/projects
找到 3 个 node_modules 目录
结果已保存到: ~/.nmkill/node_modules_20250618_143000.csv

说明: delete 字段为 'yes' 表示将被删除 (>= 500MB)
      delete 字段为 'no' 表示不删除 (< 500MB)
      你可以手动编辑 CSV 文件来修改 delete 字段
```

### 2. Edit the CSV file

Open `~/.nmkill/node_modules_*.csv` and change `yes`/`no` to your preference.

```csv
path,size,delete
/home/user/projects/project1/node_modules,629145616,yes
/home/user/projects/project2/node_modules,104857600,no
/home/user/projects/project3/node_modules,209715200,yes
```

### 3. Preview deletion (dry-run)

```bash
nmkill execute --dry-run
```

### 4. Execute deletion

```bash
# With confirmation
nmkill execute

# Skip confirmation (for automation)
nmkill execute --yes
```

## CSV File Location

All CSV files are stored in `~/.nmkill/` with timestamps:

```
~/.nmkill/
├── node_modules_20250618_143000.csv
├── node_modules_20250618_144500.csv
└── node_modules_20250618_151000.csv
```

## Options

| Command | Description |
|---------|-------------|
| `nmkill query [path]` | Scan for node_modules (default: current directory) |
| `nmkill execute [csv]` | Delete marked node_modules (uses latest CSV if not specified) |

| Flag | Description |
|------|-------------|
| `-n, --dry-run` | Preview deletion without actually deleting |
| `-y, --yes` | Skip confirmation prompt |

## Size Threshold

By default, node_modules >= 500MB are marked `yes` for deletion.

## License

MIT
