package main

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func uninstallSync() error {
	{
		cmd := exec.Command("taskkill", "/im", "telltail-sync.exe")
		cmd.Output()
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return cli.Exit("Cannot determine your home folder", exitCannotDetermineUserHomeDir)
	}

	{
		loc := filepath.Join(homeDir, startupPath, "telltail-sync.ahk")
		removeFileIfPresent(loc)
	}

	{
		dir := filepath.Join(homeDir, binPath)
		removeFileIfPresent(filepath.Join(dir, "clipnotify.exe"))
		removeFileIfPresent(filepath.Join(dir, "telltail-sync.exe"))
		removeDirIfEmpty(dir) // FIXME this doesn't removes the dir even if its empty, check why
	}
	return nil
}

func uninstallCenter() error {
	return nil
}
