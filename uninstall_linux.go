package main

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func fileOrFolderExists(fullpath string) bool {
	_, err := os.Stat(fullpath)
	return !errors.Is(err, fs.ErrNotExist)
}

func removeFileIfPresent(fullpath string) error {
	if fileOrFolderExists(fullpath) {
		if err := os.Remove(fullpath); err != nil {
			return cli.Exit("Unable to remove "+filepath.Base(fullpath)+" from path "+filepath.Dir(fullpath), exitFileNotWriteable)
		}
	}
	return nil
}

func removeFolderIfPresent(fullpath string) error {
	if fileOrFolderExists(fullpath) {
		if err := os.RemoveAll(fullpath); err != nil {
			return cli.Exit("Unable to remove "+fullpath, exitDirNotCreatable)
		}
	}
	return nil
}

func uninstallSyncOnLinux() error {
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

	// fmt.Println("Uninstalled") << TODO do we need this
	return nil
}

func uninstallCenterOnLinux() error {
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

	// fmt.Println("Uninstalled")
	return nil
}
