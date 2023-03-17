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
	return fileOrFolderExists(filepath.Join(home, startupPath, "telltail-"+name+".service"))
}

func manageService(name string, action serviceAction) error {
	if !contains(validServices, name) {
		return cli.Exit("Invalid service name. Valid values are: "+strings.Join(validServices, ", "), exitServiceNotPresent)
	}
	if !isServiceAvailable(name) {
		// can also occur if user dir is not identifiable
		return cli.Exit("This service is unavailable. Install it first before you act upon it.", exitServiceNotPresent)
	}

	// We've intentionally used the word start/stop over enable/disable because the
	// latter in our mind means load up the service and enable it or disable the service and unload it.
	// This is not we want. We expect telltail to be intentionally disabled by the user for some period of time, not permanently.
	// That's why we'll make sure that telltail loads up on boot by `systemctl enable`-ing the service
	cmd := exec.Command("systemctl", "--user", "enable", "tellail-"+name)
	cmd.Run()

	// config might have changed by us or manually by the user, and this cmd loads that new config
	switch action {
	case startService:
	case restartService:
		cmd = exec.Command("systemctl", "--user", "daemon-reload")
		cmd.Run()
	}

	switch action {
	case startService:
		cmd = exec.Command("systemctl", "--user", "start", "telltail-"+name)
	case stopService:
		cmd = exec.Command("systemctl", "--user", "stop", "telltail-"+name)
	case restartService:
		cmd = exec.Command("systemctl", "--user", "restart", "telltail-"+name)
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
	fullpath := filepath.Join(home, startupPath, "telltail-center.service.d", "override.conf")
	input, err := os.ReadFile(fullpath)
	if err != nil {
		return cli.Exit("Cannot find the config file", exitFileNotReadable)
	}

	ableToParseExistingAuthKey := false

	lines := strings.Split(string(input), "\n")
	startStr := "Environment=TS_AUTHKEY="
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
