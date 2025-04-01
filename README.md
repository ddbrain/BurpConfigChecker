# BurpConfigChecker

BurpConfigChecker is a lightweight command-line application written in Go that reads Burp Suite user options files (`UserConfigPro.json` or `UserConfigCommunity.json`) and returns the value of a specified setting key.

## Features
- Supports both Burp Suite Pro and Community configuration files.
- Automatically detects the correct configuration file path based on your operating system.
- **Verbose Mode** provides detailed information, including:
  - Path to the JSON file
  - Line number in the JSON file where the key was found (if available)
  - The **Setting Key** and **Value**
  
## Default Paths
| OS      | Preferred File Location                        | Fallback Location                                   |
|---------|------------------------------------------------|---------------------------------------------------|
| Windows | `%USERPROFILE%\AppData\Roaming\BurpSuite\UserConfigPro.json` | `%USERPROFILE%\AppData\Roaming\BurpSuite\UserConfigCommunity.json` |
| Linux   | `$HOME/.BurpSuite/UserConfigPro.json`     | `$HOME/.BurpSuite/UserConfigCommunity.json`  |
| macOS   | `$HOME/Library/Application Support/BurpSuite/UserConfigPro.json` | `$HOME/Library/Application Support/BurpSuite/UserConfigCommunity.json` |

---

## Quickstart
### Compilation
To compile BurpConfigChecker for different operating systems, use the following commands. All binaries will be placed in the `build` directory.

### Windows Compilation (CMD)
```
mkdir build
set GOOS=windows
set GOARCH=amd64
go build -o build\BurpConfigChecker.exe BurpConfigChecker.go
```

### Windows Compilation (PowerShell)
```
mkdir build
$env:GOOS="windows"
$env:GOARCH="amd64"
go build -o build\BurpConfigChecker.exe BurpConfigChecker.go
```

### Linux Compilation
```
mkdir build
GOOS=linux GOARCH=amd64 go build -o build/BurpConfigChecker BurpConfigChecker.go
```

### macOS Compilation (Intel)
```
mkdir build
GOOS=darwin GOARCH=amd64 go build -o build/BurpConfigChecker BurpConfigChecker.go
```

### macOS Compilation (Apple Silicon)
```
mkdir build
GOOS=darwin GOARCH=arm64 go build -o build/BurpConfigChecker BurpConfigChecker.go
```

---

## Usage
```
BurpConfigChecker [-f <path_to_json_file>] [-v] [-h] <setting_key>
```

### Options:
- `-f` : Path to the Burp Suite user options JSON file (If not provided, the program uses the default path for your OS).
- `-v` : Enable verbose output (shows file location, line number, setting key, and value).
- `-h` : Show help message.

### Example:
#### Non-verbose mode
```
BurpConfigChecker -f /path/to/file.json project_options.connections.proxy.enabled
```
- **Output** (only the value if found, nothing if not found):
```
true
```

#### Verbose mode
```
BurpConfigChecker -f /path/to/file.json -v project_options.connections.proxy.enabled
```
- **Output** if found:
```
Reading from: /path/to/file.json
Found on line: 14
Setting Key: project_options.connections.proxy.enabled
Setting Value: true
```
- **Output** if not found:
```
Setting not found.
```

---

## License
MIT License
