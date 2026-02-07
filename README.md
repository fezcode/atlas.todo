# Atlas Todo

![Banner Image](./banner-image.png)

**atlas.todo** is a fast, keyboard-centric terminal user interface (TUI) for task management. Part of the **Atlas Suite**, it helps you organize your life with local-first storage, smart grouping, and Vim-like navigation.

![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey)

## ✨ Features

- 📊 **Smart Grouping:** Cycle views by Category, Day, or Priority with a single key.
- ⌨️ **Vim Bindings:** Navigate, edit, and move tasks without leaving the keyboard.
- 🏷️ **Metadata Parsing:** Add `@category` or `!priority` (!high, !med, !low) directly in the task title.
- 🔍 **Real-time Search:** Filter tasks instantly as you type.
- 💾 **Local First:** Your data lives in `~/.atlas/todo.json`—no cloud, no latency.
- 📦 **Cross-Platform:** Binaries available for Windows, Linux, and macOS.

## 🚀 Installation

### From Source
```bash
git clone https://github.com/fezcode/atlas.todo
cd atlas.todo
go build -o atlas.todo .
```

## ⌨️ Usage

Simply run the binary to enter the TUI:
```bash
./atlas.todo
```

### CLI Quick Add
Add a task without opening the UI (ensure you use quotes for tasks with metadata):
```bash
./atlas.todo add "Finish the report @work !high"
```

### CLI List Mode (MOTD)
Display your top tasks and exit. Perfect for your shell's startup script.

```bash
# Show 5 tasks (default)
./atlas.todo list

# Show 3 tasks
./atlas.todo list 3

# Show top 5 high-priority tasks
./atlas.todo list desc 5
```

### Shell Integration
You can add `atlas.todo list` to your shell profile to see your tasks every time you open a terminal.

**Bash (`~/.bashrc`) or Zsh (`~/.zshrc`):**
```bash
# Add to the end of your config file
/path/to/atlas.todo list desc 5
```

**PowerShell (`$PROFILE`):**
```powershell
# Add to the end of your profile script
& "C:\path\to\atlas.todo.exe" list desc 5
```

**Fish (`~/.config/fish/config.fish`):**
```fish
# Add to the end of your config file
/path/to/atlas.todo list desc 5
```

### ✨ Pro Tip: Advanced PowerShell Dashboard
Want a more "alive" terminal? Copy this into your `$PROFILE` for a personalized greeting, a custom ASCII banner, and color-coded task highlighting:

```powershell
# --- 1. Setup ---
$hour = (Get-Date).Hour
$greeting = "Good Evening"
if ($hour -lt 12) { $greeting = "Good Morning" }
elseif ($hour -lt 18) { $greeting = "Good Afternoon" }

# --- 2. Banner ---
Write-Host "    ____                     _       " -ForegroundColor Magenta
Write-Host "   / __/___ ___  _________  / /__    " -ForegroundColor Magenta
Write-Host "  / /_ / _ \_  \/ ___/ __ \/ _  /___ " -ForegroundColor Magenta
Write-Host " / __/  __// /_/ /__/ /_/ / /_/ / -_)" -ForegroundColor Magenta
Write-Host "/_/  \___/____/\___/\____/\____/\__/ " -ForegroundColor Magenta

Write-Host "`n  $greeting, $($env:USERNAME)! Focus on these items:`n" -ForegroundColor White

# --- 3. Tasks ---
if (Get-Command atlas.todo -ErrorAction SilentlyContinue) {
    $taskList = & atlas.todo list desc 5

    if ($null -eq $taskList -or $taskList.Count -eq 0 -or ($taskList -join " ") -match "No pending tasks") {
        Write-Host "  ✨ Your board is clear! Ready for something new?" -ForegroundColor Green
    } else {
        foreach ($line in $taskList) {
            if ($line -like "[!]*") { 
                Write-Host "  $line" -ForegroundColor Red -NoNewline
                Write-Host "  🔥" -ForegroundColor Red
            }
            elseif ($line -like "[.]*") { 
                Write-Host "  $line" -ForegroundColor DarkGray 
            }
            else { 
                Write-Host "  $line" -ForegroundColor Yellow 
            }
        }
    }
}

# --- 4. Footer & Alias ---
Write-Host "`n  --------------------------------------------------" -ForegroundColor DarkGray
Write-Host "  Run 't' or 'atlas.todo' to manage tasks." -ForegroundColor Gray
Write-Host ""

if (-not (Get-Alias t -ErrorAction SilentlyContinue)) {
    New-Alias -Name t -Value atlas.todo
}
```

## 🕹️ Controls

| Key | Action |
|-----|--------|
| `↑/↓` or `k/j` | Navigate tasks |
| `Space` | Toggle task completion |
| `n` | Create a new task |
| `/` | Search/Filter tasks |
| `g` | Cycle grouping (None, Category, Day, Priority) |
| `s` | Cycle sorting (Default, Asc ↑, Desc ↓) |
| `c` | Toggle showing completed tasks |
| `d` | Delete task (requires confirmation) |
| `q` or `Esc` | Quit |

## 🏗️ Building for all platforms

The project includes a Go-based build script to generate binaries for all platforms:

```bash
go run scripts/build.go
```
Binaries will be placed in the `build/` directory.

## 📄 License
MIT License - see [LICENSE](LICENSE) for details.
