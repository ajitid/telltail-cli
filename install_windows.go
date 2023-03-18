package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/urfave/cli/v2"
)

const startupPath = "AppData\\Roaming\\Microsoft\\Windows\\Start Menu\\Programs\\Startup"
const binPath = ".local\\share\\telltail"

func installSync(params installSyncParams) error {
	////// Check basic necessities exist
	if !cmdExists("autohotkey.exe") {
		return cli.Exit("AutoHotkey is not present. We need that to run this service everytime you log in.\n"+
			"You install it for free via https://www.autohotkey.com. Once installed, come back and rerun this command to continue the setup.", exitMissingDependency)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return cli.Exit("Cannot determine your home folder", exitCannotDetermineUserHomeDir)
	}
	baseBinLoc := filepath.Join(homeDir, binPath)

	////// stop any running processes first, otherwise windows won't let us override them
	{
		cmd := exec.Command("taskkill", "/im", "telltail-sync.exe")
		cmd.Output()
	}

	////// Download and store clipnotify
	{
		if !cmdExists("tar.exe") {
			return cli.Exit("Please upgrade to Windows 10 build 17063 or higher as we need `tar.exe` to extract a zip file.", exitMissingDependency)
		}

		zipLoc := filepath.Join(baseBinLoc, "clipnotify.zip")
		err, exitCode := downloadFile(
			"https://github.com/ajitid/clipnotify-for-desktop-os/releases/download/"+version+"/clipnotify-win-"+runtime.GOARCH+".zip",
			zipLoc)
		if err != nil {
			return cli.Exit(err, exitCode)
		}

		err = removeFolderIfPresent(filepath.Join(baseBinLoc, "clipnotify"))
		if err != nil {
			return cli.Exit("Couldn't delete existing clipnotify folder", exitDirNotModifiable)
		}
		extract := exec.Command("tar.exe", "-xf", "clipnotify.zip")
		extract.Dir = baseBinLoc
		_, err = extract.Output()
		if err != nil {
			return cli.Exit("Couldn't extract clipnotify.zip", exitFileNotModifiable)
		}
		err = removeFileIfPresent(zipLoc)
		if err != nil {
			fmt.Println("Couldn't delete the zip, please do it by yourself. It is at:\n", zipLoc)
		}
	}

	////// Download and store the telltail-sync
	{
		loc := filepath.Join(baseBinLoc, "telltail-sync.exe")
		err, exitCode := downloadFile(
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
	if !cmdExists("autohotkey.exe") {
		return cli.Exit("AutoHotkey is not present. We need that to run this service everytime you log in.\n"+
			"You install it for free via https://www.autohotkey.com. Once installed, come back and rerun this command to continue the setup.", exitMissingDependency)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return cli.Exit("Cannot determine your home folder", exitCannotDetermineUserHomeDir)
	}
	baseBinLoc := filepath.Join(homeDir, binPath)

	////// stop any running processes first, otherwise windows won't let us override them
	{
		cmd := exec.Command("taskkill", "/im", "telltail-center.exe")
		cmd.Output()
	}

	////// Download and store the telltail-center
	{
		loc := filepath.Join(baseBinLoc, "telltail-center.exe")
		err, exitCode := downloadFile(
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
