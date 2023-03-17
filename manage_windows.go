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
