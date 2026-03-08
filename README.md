# tudu

`tudu` is a cross-platform terminal todo application scaffold written in Go.

## Installation Guide

### Install from GitHub Releases (recommended)

1. Go to your project releases page:
   - `https://github.com/<your-user>/<your-repo>/releases`
2. Download the archive for your OS:
   - Linux: `tudu_<version>_linux_amd64.tar.gz` (or `linux_arm64`)
   - macOS: `tudu_<version>_darwin_amd64.tar.gz` (or `darwin_arm64`)
   - Windows: `tudu_<version>_windows_amd64.zip`
3. Extract the archive.
4. Move the `tudu` binary to a directory in your `PATH`.
5. Run:
   - Linux/macOS: `tudu`
   - Windows (PowerShell): `tudu.exe`

### Linux (step by step)

1. Download Linux release asset from GitHub Releases.
2. Extract:
   ```bash
   tar -xzf tudu_<version>_linux_amd64.tar.gz
   ```
3. Make executable:
   ```bash
   chmod +x tudu
   ```
4. Move binary:
   ```bash
   sudo mv tudu /usr/local/bin/tudu
   ```
5. Verify:
   ```bash
   tudu
   ```

### macOS (step by step)

1. Download macOS release asset from GitHub Releases.
2. Extract:
   ```bash
   tar -xzf tudu_<version>_darwin_arm64.tar.gz
   ```
   Use `darwin_amd64` for Intel Macs.
3. Make executable:
   ```bash
   chmod +x tudu
   ```
4. Move binary:
   ```bash
   sudo mv tudu /usr/local/bin/tudu
   ```
   For Apple Silicon + Homebrew path, you can also use:
   ```bash
   sudo mv tudu /opt/homebrew/bin/tudu
   ```
5. Verify:
   ```bash
   tudu
   ```

### Windows (step by step)

1. Download `tudu_<version>_windows_amd64.zip` from GitHub Releases.
2. Extract the zip.
3. Move `tudu.exe` to a folder, for example:
   - `C:\Tools\tudu\`
4. Add that folder to your `Path` environment variable:
   - Settings -> System -> About -> Advanced system settings -> Environment Variables -> `Path` -> Edit -> New.
5. Open a new PowerShell window and verify:
   ```powershell
   tudu.exe
   ```

### Install from source (all platforms)

1. Install Go (1.22+ recommended).
2. Clone the repository:
   ```bash
   git clone https://github.com/<your-user>/<your-repo>.git
   cd <your-repo>
   ```
3. Install binary:
   ```bash
   go install ./cmd/tudu
   ```
4. Run:
   ```bash
   tudu
   ```

## Goals

- Linux, macOS, Windows support
- Installable CLI binary
- JSON-backed persistence
- Clean package boundaries for future Bubble Tea + Lipgloss UI
- Vim-like workflow (`j`, `k`, `a`, `d`, `c`, `e`, `q`)

## Project Structure

- `cmd/tudu`: application entry point and dependency wiring
- `internal/models`: domain data models
- `internal/storage`: persistence interfaces + JSON file implementation
- `internal/commands`: application-level operations
- `internal/tui`: TUI model and keybinding map scaffolding
- `examples`: sample JSON data format

## Run

```bash
go run ./cmd/tudu
```

## Build

```bash
go build ./...
```

## Install (local)

```bash
go install ./cmd/tudu
```

## JSON Storage Path (planned)

Default target path is `~/.tudu/todos.json` (platform-aware home directory resolution).

## Sample JSON

See `examples/todos.sample.json`.
