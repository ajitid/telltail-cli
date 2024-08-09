package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/urfave/cli/v2"
	"golang.org/x/sys/windows/registry"
)

const startupPath = "AppData\\Roaming\\Microsoft\\Windows\\Start Menu\\Programs\\Startup"
const binPath = ".local\\share\\telltail"

func autohotkeyInstalled(rootKey registry.Key) (installed bool, version string) {
	k, err := registry.OpenKey(rootKey, `SOFTWARE\AutoHotkey`, registry.QUERY_VALUE)
	if err != nil {
		return false, ""
	}

	v, _, err := k.GetStringValue("Version")
	if err != nil {
		return false, ""
	}
	return true, v
}

func checkForAutohotkeyv2() error {
	// CURRENT_USER -> local install, LOCAL_MACHINE -> installed using admin privileges
	// We'll first look for a local install first because GUI installer of AHK v1 always require admin privileges
	installed, version := autohotkeyInstalled(registry.CURRENT_USER)
	if !installed {
		installed, version = autohotkeyInstalled(registry.LOCAL_MACHINE)
	}

	if installed {
		if !strings.HasPrefix(version, "2.") {
			return cli.Exit("Telltail needs AutoHotkey v2 for it to work while you have v"+version+" installed.\n"+
				"Installing AHK v2 will not break your existing scripts. After installing it, come back and rerun this command to continue the setup.", exitMissingDependency)
		}
		/*
			BUG: if you install AHK v1 followed by v2 both using GUI installers with admin privileges and then
			uninstall v2 (but keep v1), AHK does not update the registry (fully) and shows v2 there. This means
			that `telltail install` will falsely assume that the user has v2 installed and continue with the script.
			While we can't do much about it, the AHK scripts we run themselves have a minimum required version flag set,
			so we have a safety net there.
		*/
	} else {
		return cli.Exit("AutoHotkey is not present. We need that to run this service everytime you log in.\n"+
			"You can install v2 for free via https://www.autohotkey.com. Once installed, come back and rerun this command to continue the setup.", exitMissingDependency)
	}
	return nil
}

func installSync(params installSyncParams) error {
	////// Check basic necessities exist
	{
		err := checkForAutohotkeyv2()
		if err != nil {
			return err
		}
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return cli.Exit("Cannot determine your home folder", exitCannotDetermineUserHomeDir)
	}
	baseBinLoc := filepath.Join(homeDir, binPath)

	////// stop any running processes first, otherwise windows won't let us override them
	{
		cmd := exec.Command("taskkill", "/im", "telltail-sync.exe")
		cmd.Run()
	}

	////// Download and store clipnotify
	{
		loc := filepath.Join(baseBinLoc, "clipnotify.exe")
		exitCode, err := downloadFile(
			"https://github.com/ajitid/clipnotify-for-desktop-os/releases/download/"+version+"/clipnotify-win-"+runtime.GOARCH+".exe",
			loc)
		if err != nil {
			return cli.Exit(err, exitCode)
		}
	}

	////// Download and store the telltail-sync
	{
		loc := filepath.Join(baseBinLoc, "telltail-sync.exe")
		exitCode, err := downloadFile(
			"https://github.com/ajitid/telltail-sync/releases/download/"+version+"/telltail-sync-win-"+runtime.GOARCH+".exe",
			loc)
		if err != nil {
			return cli.Exit(err, exitCode)
		}
	}

	////// Put bootup configuration and start the service
	{
		dir := filepath.Join(homeDir, startupPath)
		loc := filepath.Join(dir, "telltail-sync.ahk")
		tmpl := getSyncAhkCfg()
		f, err := os.Create(loc)
		if err != nil {
			return cli.Exit("Cannot create AutoHotkey script", exitFileNotModifiable)
		}
		err = tmpl.Execute(f, syncAhkCfgAttrs{
			Tailnet:      params.tailnet,
			Device:       params.device,
			BinDirectory: baseBinLoc,
		})
		if err != nil {
			return cli.Exit("Cannot write to AutoHotkey script", exitFileNotModifiable)
		}
		f.Close()

		// from https://stackoverflow.com/a/50532038
		cmd := exec.Command("cmd.exe", "/C", "start", "/b", ".\\telltail-sync.ahk")
		cmd.Dir = dir
		if err := cmd.Run(); err != nil {
			log.Println("Error:", err)
		}
	}

	////// Success message
	fmt.Println("All done! You can read about the changes we've made on here: https://guide-on.gitbook.io/telltail/changes-done-by-install")
	return nil
}

func installCenter(authKey string) error {
	////// Check basic necessities exist
	{
		err := checkForAutohotkeyv2()
		if err != nil {
			return err
		}
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return cli.Exit("Cannot determine your home folder", exitCannotDetermineUserHomeDir)
	}
	baseBinLoc := filepath.Join(homeDir, binPath)

	////// stop any running processes first, otherwise windows won't let us override them
	{
		cmd := exec.Command("taskkill", "/im", "telltail-center.exe")
		cmd.Run()
	}

	////// Download and store the telltail-center
	{
		loc := filepath.Join(baseBinLoc, "telltail-center.exe")
		exitCode, err := downloadFile(
			"https://github.com/ajitid/telltail-center/releases/download/"+version+"/telltail-center-win-"+runtime.GOARCH+".exe",
			loc)
		if err != nil {
			return cli.Exit(err, exitCode)
		}
	}

	////// Put bootup configuration and start the service
	{
		dir := filepath.Join(homeDir, startupPath)
		loc := filepath.Join(dir, "telltail-center.ahk")
		tmpl := getCenterAhkCfg()
		f, err := os.Create(loc)
		if err != nil {
			return cli.Exit("Cannot create AutoHotkey script", exitFileNotModifiable)
		}
		err = tmpl.Execute(f, centerAhkCfgAttrs{
			BinDirectory: baseBinLoc,
			AuthKey:      authKey,
		})
		if err != nil {
			return cli.Exit("Cannot write to AutoHotkey script", exitFileNotModifiable)
		}
		f.Close()

		// from https://stackoverflow.com/a/50532038
		cmd := exec.Command("cmd.exe", "/C", "start", "/b", ".\\telltail-center.ahk")
		cmd.Dir = dir
		if err := cmd.Run(); err != nil {
			log.Println("Error:", err)
		}
	}

	////// Success message
	fmt.Println("All done! You can read about the changes we've made on here: https://guide-on.gitbook.io/telltail/changes-done-by-install")
	return nil
}
