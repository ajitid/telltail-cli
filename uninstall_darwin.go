package main

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func uninstallSync() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return cli.Exit("Cannot determine your home folder", exitCannotDetermineUserHomeDir)

	}
	if cmdExists("systemctl") {
		cmd := exec.Command("launchctl", "unload", filepath.Join(homeDir, startupPath, scriptPrefix+"telltail-sync.plist"))
		cmd.Run()
	}

	err = removeFileIfPresent(filepath.Join(homeDir, startupPath, scriptPrefix+"telltail-sync.plist"))
	if err != nil {
		return err
	}

	baseBinLoc := filepath.Join(homeDir, binPath)
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

	return nil
}

func uninstallCenter() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return cli.Exit("Cannot determine your home folder", exitCannotDetermineUserHomeDir)
	}

	if cmdExists("systemctl") {
		cmd := exec.Command("launchctl", "unload", filepath.Join(homeDir, startupPath, scriptPrefix+"telltail-center.plist"))
		cmd.Run()
	}

	err = removeFileIfPresent(filepath.Join(homeDir, startupPath, scriptPrefix+"telltail-center.plist"))
	if err != nil {
		return err
	}

	baseBinLoc := filepath.Join(homeDir, binPath)
	err = removeFileIfPresent(filepath.Join(baseBinLoc, "telltail-center"))
	if err != nil {
		return err
	}

	err = removeDirIfEmpty(baseBinLoc)
	if err != nil {
		return err
	}

	return nil
}
