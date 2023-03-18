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
	return fileOrFolderExists(filepath.Join(home, startupPath, scriptPrefix+"telltail-"+name+".plist"))
}

func manageService(name string, action serviceAction) error {
	if !contains(validServices, name) {
		return cli.Exit("Invalid service name. Valid values are: "+strings.Join(validServices, ", "), exitServiceNotPresent)
	}
	if !isServiceAvailable(name) {
		// can also occur if user dir is not identifiable
		return cli.Exit("This service is unavailable. Install it first before you act upon it.", exitServiceNotPresent)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return cli.Exit("Cannot determine your home folder", exitCannotDetermineUserHomeDir)
	}
	dir := filepath.Join(homeDir, startupPath)

	// when user asks for start/restart we also make sure that the script loads up on next boot as well
	// Sidenote: oddly `load` might or might or not `start` the process (usually it does, but I've seen at times it doesn't)
	//	but `unload` fo' sure `stop` the process
	switch action {
	case startService:
	case restartService:
		cmd := exec.Command("launchctl", "load", filepath.Join(dir, scriptPrefix+"telltail-"+name+".plist"))
		cmd.Run()
	}

	// to observe the effect of these start/stop/load/unload commands, focus on disappearance and appearance of  % CPU
	// in Activity Monitor (and not on the disappearance and appearance of the row itself)
	// launch actions are not immediate — command itself is executed immediately but seems like launchd executes stuff in batches, and hence not immediate
	var cmd *exec.Cmd
	switch action {
	case startService:
		cmd = exec.Command("launchctl", "start", scriptPrefix+"telltail-"+name)
	case stopService:
		cmd = exec.Command("launchctl", "stop", scriptPrefix+"telltail-"+name)
	case restartService:
		cmd = exec.Command("launchctl", "stop", scriptPrefix+"telltail-"+name)
		cmd.Run()
		cmd = exec.Command("launchctl", "start", scriptPrefix+"telltail-"+name)
	}
	cmd.Run()

	return nil
}

func editCenterAuthKey() error {
	// tried a lot of stuff, but somehow env. var persists across stops, even across restarts
	// because of this even changing env var and unload -> load doesn't affect anything (1)
	// and the service runs even if I uninstall Center and reboot (2). `stop` works but only for that login session.
	// Load takes in a `-w` flag, but it doesn't seem to affect either. I need systemctl daemon-reload type of thing.
	// Need to investigate both (1) and (2) before enabling this feature.
	if true {
		return cli.Exit("This feature is unavailable for now", 2)
	}
	///////////////////////////

	if !isServiceAvailable("center") {
		// can also occur if user dir is not identifiable
		return cli.Exit("This service is unavailable. Install it first before you act upon it.", exitServiceNotPresent)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return cli.Exit("Cannot determine your home folder", exitCannotDetermineUserHomeDir)
	}
	fullpath := filepath.Join(home, startupPath, scriptPrefix+"telltail-center.plist")
	input, err := os.ReadFile(fullpath)
	if err != nil {
		return cli.Exit("Cannot find the config file", exitFileNotReadable)
	}

	ableToParseExistingAuthKey := false

	lines := strings.Split(string(input), "\n")
	for i, line := range lines {
		// we'll use a rudimentary way to parse this stuff for now
		if strings.TrimSpace(line) != "<key>TS_AUTHKEY</key>" {
			continue
		}
		if i+1 < len(lines) {
			authKeyLine := strings.TrimSpace(lines[i+1])
			if strings.HasPrefix(authKeyLine, "<string>") && strings.HasSuffix(authKeyLine, "</string>") {
				ableToParseExistingAuthKey = true
				existingAuthKey := strings.TrimPrefix(authKeyLine, "<string>")
				existingAuthKey = strings.TrimSuffix(existingAuthKey, "</string>")

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
					lines[i+1] = "    " + "<string>" + newKey + "</string>"
				}
				break
			}
		}
	}

	if !ableToParseExistingAuthKey {
		return cli.Exit("Unable to change auth key because we've found an invalid config", exitInvalidConfig)
	}

	output := strings.Join(lines, "\n")
	// https://stackoverflow.com/a/18415935/7683365
	err = os.WriteFile(fullpath, []byte(output), 0644)
	if err != nil {
		return cli.Exit("Unable to write to config file", exitFileNotModifiable)
	}

	manageService("center", restartService)
	return nil
}
