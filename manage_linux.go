package main

import (
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
