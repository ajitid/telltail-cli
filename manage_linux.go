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
		return cli.Exit("This service is unavailable. Install it first before you act upon it.", exitServiceNotPresent)
	}

	var cmd *exec.Cmd
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
