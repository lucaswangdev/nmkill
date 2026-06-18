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

### 🪟 Windows

**PowerShell (Run as Administrator):**

```powershell
Invoke-WebRequest -Uri "https://github.com/lucaswangdev/nmkill/releases/download/v0.1.0/nmkill-windows-amd64.exe" -OutFile "C:\Windows\System32\nmkill.exe"
```

**验证:**
```powershell
nmkill --version
```

---

### 🍎 macOS (Intel)

**Terminal:**

```bash
curl -L https://github.com/lucaswangdev/nmkill/releases/download/v0.1.0/nmkill-darwin-amd64 -o /tmp/nmkill && chmod +x /tmp/nmkill && sudo mv /tmp/nmkill /usr/local/bin/nmkill
```

**验证:**
```bash
nmkill --version
```

---

### 🍎 macOS (Apple Silicon / M1/M2/M3)

**Terminal:**

```bash
curl -L https://github.com/lucaswangdev/nmkill/releases/download/v0.1.0/nmkill-darwin-arm64 -o /tmp/nmkill && chmod +x /tmp/nmkill && sudo mv /tmp/nmkill /usr/local/bin/nmkill
```

**验证:**
```bash
nmkill --version
```

---

### 🐧 Linux (amd64)

**Terminal:**

```bash
curl -L https://github.com/lucaswangdev/nmkill/releases/download/v0.1.0/nmkill-linux-amd64 -o /tmp/nmkill && chmod +x /tmp/nmkill && sudo mv /tmp/nmkill /usr/local/bin/nmkill
```

**验证:**
```bash
nmkill --version
```

---

### 🐧 Linux (arm64)

**Terminal:**

```bash
curl -L https://github.com/lucaswangdev/nmkill/releases/download/v0.1.0/nmkill-linux-arm64 -o /tmp/nmkill && chmod +x /tmp/nmkill && sudo mv /tmp/nmkill /usr/local/bin/nmkill
```

**验证:**
```bash
nmkill --version
```

---

### 📦 Alternative: Install via Go

```bash
go install github.com/lucaswangdev/nmkill@latest
```

> Requires [Go](https://go.dev/dl/) installed

---

## Usage

### 1️⃣ Scan for node_modules

```bash
nmkill query ~/projects
```

### 2️⃣ Edit CSV file

Open `~/.nmkill/node_modules_*.csv` and change `delete` column:
- `yes` = will be deleted (if size >= 500MB, defaults to yes)
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
| `nmkill execute [csv]` | Delete marked folders (uses latest CSV if not specified) |
| `--dry-run, -n` | Preview only, no deletion |
| `--yes, -y` | Skip confirmation |

---

## License

MIT
