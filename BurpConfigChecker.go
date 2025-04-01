package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var defaultPath string

func init() {
	var windowsProPath, windowsCommunityPath string

	if runtime.GOOS == "windows" {
		if os.Getenv("PSModulePath") != "" { // Detect PowerShell
			home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
			windowsProPath = filepath.Join(home, "AppData", "Roaming", "BurpSuite", "UserConfigPro.json")
			windowsCommunityPath = filepath.Join(home, "AppData", "Roaming", "BurpSuite", "UserConfigCommunity.json")
		} else if os.Getenv("ComSpec") != "" { // Detect CMD
			windowsProPath = os.ExpandEnv("%USERPROFILE%\\AppData\\Roaming\\BurpSuite\\UserConfigPro.json")
			windowsCommunityPath = os.ExpandEnv("%USERPROFILE%\\AppData\\Roaming\\BurpSuite\\UserConfigCommunity.json")
		}

		if fileExists(windowsProPath) {
			defaultPath = windowsProPath
		} else {
			defaultPath = windowsCommunityPath
		}
	} else {
		var linuxProPath = os.ExpandEnv("$HOME/.BurpSuite/UserConfigPro.json")
		var linuxCommunityPath = os.ExpandEnv("$HOME/.BurpSuite/UserConfigCommunity.json")
		var macProPath = os.ExpandEnv("$HOME/Library/Application Support/BurpSuite/UserConfigPro.json")
		var macCommunityPath = os.ExpandEnv("$HOME/Library/Application Support/BurpSuite/UserConfigCommunity.json")

		switch runtime.GOOS {
		case "linux":
			if fileExists(linuxProPath) {
				defaultPath = linuxProPath
			} else {
				defaultPath = linuxCommunityPath
			}
		case "darwin": // macOS
			if fileExists(macProPath) {
				defaultPath = macProPath
			} else {
				defaultPath = macCommunityPath
			}
		default:
			// Fallback if OS is unknown
			defaultPath = "./user-options.json"
		}
	}
}

func main() {
	var jsonFilePath string
	var showHelp bool
	var verbose bool

	flag.StringVar(&jsonFilePath, "f", defaultPath, "Path to the Burp Suite user options JSON file")
	flag.BoolVar(&showHelp, "h", false, "Show help message")
	flag.BoolVar(&verbose, "v", false, "Enable verbose output")
	flag.Parse()

	// If user asked for help or didn't provide a setting key
	if showHelp || flag.NArg() < 1 {
		printHelp()
		os.Exit(0)
	}

	settingKey := flag.Arg(0)

	fileContent, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		os.Exit(1)
	}

	var data map[string]interface{}
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		fmt.Printf("Failed to parse JSON: %v\n", err)
		os.Exit(1)
	}

	// Search for the setting value recursively
	result := findSetting(data, strings.Split(settingKey, "."))

	if result != nil {
		if verbose {
			// Attempt to find line number in the raw JSON for the final sub-key
			lines := strings.Split(string(fileContent), "\n")
			lastKey := settingKey[strings.LastIndex(settingKey, ".")+1:]
			lineNum := -1
			for i, line := range lines {
				if strings.Contains(line, lastKey) {
					lineNum = i + 1
					break
				}
			}

			fmt.Printf("Reading from: %s\n", jsonFilePath)
			if lineNum > 0 {
				fmt.Printf("Found on line: %d\n", lineNum)
			}
			fmt.Printf("Setting Key: %s\n", settingKey)
			fmt.Printf("Setting Value: %v\n", result)
		} else {
			// Non-verbose mode: print raw value only
			fmt.Printf("%v\n", result)
		}
	} else {
		if verbose {
			fmt.Println("Setting not found.")
		}
		// Default (non-verbose) mode shows no output if not found.
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func printHelp() {
	fmt.Println("Usage: BurpConfigChecker [-f <path_to_json_file>] [-v] [-h] <setting_key>")
	fmt.Println("Options:")
	fmt.Printf("  -f  Path to the Burp Suite user options JSON file (default: %s)\n", defaultPath)
	fmt.Println("  -v  Enable verbose output")
	fmt.Println("  -h  Show help message")
	fmt.Println("")
	fmt.Println("Example:")
	fmt.Println("  BurpConfigChecker -f /path/to/file.json -v project_options.connections.proxy.enabled")
}

// findSetting recursively searches a map/array structure for a dot-delimited key path
func findSetting(data interface{}, keys []string) interface{} {
	if len(keys) == 0 {
		return data
	}

	switch d := data.(type) {
	case map[string]interface{}:
		if val, ok := d[keys[0]]; ok {
			return findSetting(val, keys[1:])
		}
	case []interface{}:
		// If it's an array, try searching each element
		for _, item := range d {
			if res := findSetting(item, keys); res != nil {
				return res
			}
		}
	}
	return nil
}
