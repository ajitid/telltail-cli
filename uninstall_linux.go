package main

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func uninstallSync() error {
	if cmdExists("systemctl") {
		cmd := exec.Command("systemctl", "--user", "disable", "telltail-sync", "--now")
		cmd.Output()
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return cli.Exit("Cannot determine your home folder", exitCannotDetermineUserHomeDir)
	}

	err = removeFileIfPresent(filepath.Join(homeDir, ".config/systemd/user/telltail-sync.service"))
	if err != nil {
		return err
	}

	baseBinLoc := filepath.Join(homeDir, ".local/share/telltail")
	err = removeFileIfPresent(filepath.Join(baseBinLoc, "clipnotify"))
	if err != nil {
		return err
	}
	err = removeFileIfPresent(filepath.Join(baseBinLoc, "telltail-sync"))
	if err != nil {
		return err
	}

	err = removeDirIfEmpty(baseBinLoc)
	if err != nil {
		return err
	}

	// fmt.Println("Uninstalled") << TODO do we need this
	return nil
}

func uninstallCenter() error {
	if cmdExists("systemctl") {
		cmd := exec.Command("systemctl", "--user", "disable", "telltail-center", "--now")
		cmd.Output()
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return cli.Exit("Cannot determine your home folder", exitCannotDetermineUserHomeDir)
	}

	err = removeFileIfPresent(filepath.Join(homeDir, ".config/systemd/user/telltail-center.service"))
	if err != nil {
		return err
	}
	err = removeFolderIfPresent(filepath.Join(homeDir, ".config/systemd/user/telltail-center.service.d/"))
	if err != nil {
		return err
	}

	baseBinLoc := filepath.Join(homeDir, ".local/share/telltail")
	err = removeFileIfPresent(filepath.Join(baseBinLoc, "telltail-center"))
	if err != nil {
		return err
	}

	err = removeDirIfEmpty(baseBinLoc)
	if err != nil {
		return err
	}

	// fmt.Println("Uninstalled")
	return nil
}
