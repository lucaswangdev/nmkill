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

### Uninstall Old Version (if exists)

**Linux / macOS:**
```bash
sudo rm /usr/local/bin/nmkill
```

**Windows:**
```powershell
Remove-Item C:\Windows\System32\nmkill.exe
```

---

### 🪟 Windows

**PowerShell (Run as Administrator):**

```powershell
# Check latest version at: https://github.com/lucaswangdev/nmkill/releases
Invoke-WebRequest -Uri "https://github.com/lucaswangdev/nmkill/releases/download/v0.1.4/nmkill-windows-amd64.exe" -OutFile "C:\Windows\System32\nmkill.exe"
```

---

### 🍎 macOS (Intel)

**Terminal:**

```bash
# Check latest version at: https://github.com/lucaswangdev/nmkill/releases
curl -L https://github.com/lucaswangdev/nmkill/releases/download/v0.1.4/nmkill-darwin-amd64 -o /tmp/nmkill && chmod +x /tmp/nmkill && sudo mv /tmp/nmkill /usr/local/bin/nmkill
```

---

### 🍎 macOS (Apple Silicon / M1/M2/M3)

**Terminal:**

```bash
# Check latest version at: https://github.com/lucaswangdev/nmkill/releases
curl -L https://github.com/lucaswangdev/nmkill/releases/download/v0.1.4/nmkill-darwin-arm64 -o /tmp/nmkill && chmod +x /tmp/nmkill && sudo mv /tmp/nmkill /usr/local/bin/nmkill
```

---

### 🐧 Linux (amd64)

**Terminal:**

```bash
# Check latest version at: https://github.com/lucaswangdev/nmkill/releases
curl -L https://github.com/lucaswangdev/nmkill/releases/download/v0.1.4/nmkill-linux-amd64 -o /tmp/nmkill && chmod +x /tmp/nmkill && sudo mv /tmp/nmkill /usr/local/bin/nmkill
```

---

### 🐧 Linux (arm64)

**Terminal:**

```bash
# Check latest version at: https://github.com/lucaswangdev/nmkill/releases
curl -L https://github.com/lucaswangdev/nmkill/releases/download/v0.1.4/nmkill-linux-arm64 -o /tmp/nmkill && chmod +x /tmp/nmkill && sudo mv /tmp/nmkill /usr/local/bin/nmkill
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
