package main

import (
	"fmt"
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
		removeFolderIfPresent(filepath.Join(dir, "clipnotify"))
		removeFileIfPresent(filepath.Join(dir, "telltail-sync.exe"))
		removeDirIfEmpty(dir) // Doesn't remove the dir because windows says the access to this dir is held by some other program. Works if I rerun the uninstall though.
	}
	return nil
}

func uninstallCenter() error {
	{
		cmd := exec.Command("taskkill", "/im", "telltail-center.exe")
		cmd.Output()
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return cli.Exit("Cannot determine your home folder", exitCannotDetermineUserHomeDir)
	}

	{
		loc := filepath.Join(homeDir, startupPath, "telltail-center.ahk")
		removeFileIfPresent(loc)
	}

	{
		dir := filepath.Join(homeDir, binPath)
		removeFileIfPresent(filepath.Join(dir, "telltail-center.exe"))
		removeDirIfEmpty(dir) // Doesn't remove the dir because windows says the access to this dir is held by some other program. Works if I rerun the uninstall though.
		fmt.Println("If you are not planning to use Telltail Center on this device anytime soon and have installed AutoHotkey only for using Telltail Center, you can remove AutoHotkey as well.")
	}
	return nil
}
