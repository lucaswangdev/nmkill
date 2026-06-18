# nmkill

> A CLI tool to find and clean `node_modules` directories, freeing up disk space.

## Features

- 🚀 Fast scanning
- 📊 CSV export with timestamps
- 🔒 Safe deletion (dry-run mode)
- ✏️ Editable deletion list
- 🌐 Cross-platform (macOS, Linux, Windows)

---

## Installation

### Option 1: Go Install (Recommended)

```bash
go install github.com/lucaswangdev/nmkill@latest
```

> Requires [Go](https://go.dev/dl/) installed. Auto-updates to latest version.

---

### Option 2: Download Binary

Download the latest release from [GitHub Releases](https://github.com/lucaswangdev/nmkill/releases):

| Platform | Download |
|----------|----------|
| Linux amd64 | [nmkill-linux-amd64](https://github.com/lucaswangdev/nmkill/releases/latest/download/nmkill-linux-amd64) |
| Linux arm64 | [nmkill-linux-arm64](https://github.com/lucaswangdev/nmkill/releases/latest/download/nmkill-linux-arm64) |
| macOS Intel | [nmkill-darwin-amd64](https://github.com/lucaswangdev/nmkill/releases/latest/download/nmkill-darwin-amd64) |
| macOS Apple Silicon | [nmkill-darwin-arm64](https://github.com/lucaswangdev/nmkill/releases/latest/download/nmkill-darwin-arm64) |
| Windows | [nmkill-windows-amd64.exe](https://github.com/lucaswangdev/nmkill/releases/latest/download/nmkill-windows-amd64.exe) |

After download, run:

```bash
chmod +x nmkill-* && sudo mv nmkill-* /usr/local/bin/nmkill
```

---

## Usage

### 1️⃣ Scan for node_modules

```bash
nmkill query ~/projects
```

### 2️⃣ Edit CSV file

Open `~/.nmkill/node_modules_*.csv` and change `delete` column:
- `yes` = will be deleted
- `no` = will be kept

### 3️⃣ Preview (dry-run)

```bash
nmkill execute --dry-run
```

### 4️⃣ Execute deletion

```bash
nmkill execute --yes
```

---

## Options

| Flag | Description |
|------|-------------|
| `nmkill query [path]` | Scan directory (default: current) |
| `nmkill execute [csv]` | Delete marked folders (uses latest CSV) |
| `--dry-run, -n` | Preview only, no deletion |
| `--yes, -y` | Skip confirmation |

---

## License

MIT
