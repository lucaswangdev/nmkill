# nmkill Project Documentation

## Project Overview

**Name:** nmkill  
**Repository:** https://github.com/lucaswangdev/nmkill  
**Purpose:** CLI tool to find and delete `node_modules` directories  
**Language:** Go 1.21  
**License:** MIT

## Architecture

```
nmkill/
├── main.go                 # Entry point
├── cmd/
│   ├── root.go            # Root cobra command
│   ├── query.go           # Scan command - finds node_modules
│   └── execute.go         # Delete command - deletes marked folders
├── scanner/
│   └── scanner.go         # Directory traversal logic
├── csv/
│   └── csv.go             # CSV read/write operations
├── delete/
│   └── delete.go          # File deletion logic
├── config/
│   └── config.go          # ~/.nmkill config directory
└── .github/
    └── workflows/
        └── release.yml    # Auto-release on git tag
```

## Key Design Decisions

### 1. CSV-based Workflow
- User edits CSV to choose what to delete (not a CLI flag)
- Reason: Safe, auditable, reversible
- CSV stored in `~/.nmkill/` with timestamps

### 2. Auto-detection of Latest CSV
- `execute` command without args uses latest CSV
- Reason: Convenient for repeated use
- No explicit CSV path needed after first query

### 3. 500MB Threshold
- `>= 500MB` → `delete=yes`
- `< 500MB` → `delete=no`
- Reason: Sensible default for typical node_modules

### 4. Cross-Platform Release
- GitHub Actions builds for: linux/{amd64,arm64}, darwin/{amd64,arm64}, windows/amd64
- CGO_ENABLED=0 for static binaries
- `-ldflags="-s -w"` to strip debug info

## CLI Commands

### `nmkill query [path]`
- Scans directory recursively for `node_modules`
- Outputs CSV to `~/.nmkill/node_modules_YYYYMMDD_HHMMSS.csv`
- Default path: current directory

### `nmkill execute [csv]`
- Deletes paths marked `yes` in CSV
- Without args: uses latest CSV in `~/.nmkill/`
- Dry-run: `--dry-run`
- Skip confirm: `--yes`

## Release Process

```bash
# 1. Make changes
git add . && git commit -m "description"

# 2. Push to main
git push origin main

# 3. Create and push tag
git tag v0.x.x
git push origin v0.x.x
```

GitHub Actions automatically:
1. Builds for all platforms
2. Creates GitHub Release
3. Uploads binaries as release assets

## Release File Naming

| Platform | Filename |
|----------|----------|
| Linux amd64 | nmkill-linux-amd64 |
| Linux arm64 | nmkill-linux-arm64 |
| macOS amd64 | nmkill-darwin-amd64 |
| macOS arm64 | nmkill-darwin-arm64 |
| Windows | nmkill-windows-amd64.exe |

## Testing

```bash
go test ./...
```

Test coverage:
- `scanner/` - Directory scanning, nested node_modules skipping
- `csv/` - CSV read/write, filtering
- `delete/` - Dry-run, real deletion, error handling

## Future Enhancements

- [ ] Add `--size-threshold` flag to customize 500MB default
- [ ] Add `--exclude` flag to skip certain directories
- [ ] Add progress bar for large directory scans
- [ ] Add JSON output format option
- [ ] Parallel scanning for faster results

## Config Directory

```
~/.nmkill/
├── node_modules_20250618_143000.csv
├── node_modules_20250618_144500.csv
└── ...
```

All data stored locally, no external services.

## Dependencies

- `github.com/spf13/cobra` - CLI framework
- Standard library only for rest

## Troubleshooting

### "command not found" after install
- Ensure `/usr/local/bin` is in PATH
- Run `echo $PATH`

### CSV file not found for execute
- Run `nmkill query` first
- Check files in `~/.nmkill/`

### Permission denied
- Use `sudo` on macOS/Linux for `/usr/local/bin`
- Run PowerShell as Administrator on Windows
