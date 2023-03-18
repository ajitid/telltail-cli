package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

func isServiceAvailable(name string) bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	return fileOrFolderExists(filepath.Join(home, startupPath, "telltail-"+name+".ahk"))
}

func manageService(name string, action serviceAction) error {
	if !contains(validServices, name) {
		return cli.Exit("Invalid service name. Valid values are: "+strings.Join(validServices, ", "), exitServiceNotPresent)
	}
	if !isServiceAvailable(name) {
		// can also occur if user dir is not identifiable
		return cli.Exit("This service is unavailable. Install it first before you act upon it.", exitServiceNotPresent)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return cli.Exit("Cannot determine your home folder", exitCannotDetermineUserHomeDir)
	}
	dir := filepath.Join(home, startupPath)

	var cmd *exec.Cmd
	switch action {
	case startService:
		cmd = exec.Command("cmd.exe", "/C", "start", "/b", ".\\telltail-"+name+".ahk")
		cmd.Dir = dir
	case stopService:
		cmd = exec.Command("taskkill", "/im", "telltail-"+name+".exe")
	case restartService:
		cmd = exec.Command("taskkill", "/im", "telltail-"+name+".exe")
		cmd.Run()
		cmd = exec.Command("cmd.exe", "/C", "start", "/b", ".\\telltail-"+name+".ahk")
		cmd.Dir = dir
	}
	cmd.Run()
	return nil
}

func editCenterAuthKey() error {
	if !isServiceAvailable("center") {
		// can also occur if user dir is not identifiable
		return cli.Exit("This service is unavailable. Install it first before you act upon it.", exitServiceNotPresent)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return cli.Exit("Cannot determine your home folder", exitCannotDetermineUserHomeDir)
	}
	fullpath := filepath.Join(home, startupPath, "telltail-center.ahk")
	input, err := os.ReadFile(fullpath)
	if err != nil {
		return cli.Exit("Cannot find the config file", exitFileNotReadable)
	}

	ableToParseExistingAuthKey := false

	lines := strings.Split(string(input), "\n")
	// I thought to first use fmt.Sscanf but consistency is sexy
	// Teacher: The author mean to say that a consistent, repeatable code across platforms is easy to refactor and is maintainable
	// The reason it is copied over and not DRY-ed is because macOS uses launchd, which uses XML for its config
	startStr := "EnvSet, TS_AUTHKEY, "
	for i, line := range lines {
		if strings.HasPrefix(line, startStr) {
			ableToParseExistingAuthKey = true

			existingAuthKey := strings.TrimPrefix(line, startStr)
			if existingAuthKey == "" {
				fmt.Println("There doesn't seem to be an existing auth key, but we can add one.")
			} else {
				fmt.Println("Existing auth key is", existingAuthKey)
			}

			fmt.Print("Enter new key (or hit return to keep the existing one): ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			newKey := scanner.Text()
			if newKey == "" {
				return nil
			} else {
				lines[i] = startStr + newKey
			}
			break
		}
	}

	if !ableToParseExistingAuthKey {
		return cli.Exit("Unable to change auth key because we've found an invalid config", exitInvalidConfig)
	}

	output := strings.Join(lines, "\n")
	err = os.WriteFile(fullpath, []byte(output), 0644)
	if err != nil {
		return cli.Exit("Unable to write to config file", exitFileNotModifiable)
	}

	manageService("center", restartService)
	return nil
}
