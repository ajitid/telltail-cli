package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/urfave/cli/v2"
)

const startupPath = "AppData\\Roaming\\Microsoft\\Windows\\Start Menu\\Programs\\Startup"

func installSync(params installSyncParams) error {
	////// Check basic necessities exist
	if !cmdExists("autohotkey.exe") {
		return cli.Exit("AutoHotkey is not present. We need that to run this program everytime you log in.\n"+
			"You install it for free via https://www.autohotkey.com. Once installed, come back and rerun this command to continue the setup.", exitMissingDependency)
	}
	// TODO FIXME stop if there's existing telltail-sync running first. Otherwise we won't be able to override it.
	// use AHK and something to stop it

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return cli.Exit("Cannot determine your home folder", exitCannotDetermineUserHomeDir)
	}
	baseBinLoc := filepath.Join(homeDir, ".local\\share\\telltail")

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
			return cli.Exit("Couldn't delete existing clipnotify folder", exitDirNotCreatable)
		}
		extract := exec.Command("tar.exe", "-xf", "clipnotify.zip")
		extract.Dir = baseBinLoc
		_, err = extract.Output()
		if err != nil {
			return cli.Exit("Couldn't extract clipnotify.zip", exitFileNotWriteable)
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
			return cli.Exit("Cannot create AutoHotkey script", exitFileNotWriteable)
		}
		err = tmpl.Execute(f, syncAhkCfgAttrs{
			Tailnet:      params.tailnet,
			Device:       params.device,
			BinDirectory: baseBinLoc,
		})
		if err != nil {
			return cli.Exit("Cannot write to AutoHotkey script", exitFileNotWriteable)
		}
		f.Close()

		cmd := exec.Command("autohotkey.exe", loc)
		_, err = cmd.Output()
		if err != nil {
			return cli.Exit("Couldn't start telltail-sync using AutoHotkey startup script", exitInvokingStartupScriptFailed)
		}
	}

	////// Success message
	fmt.Println("All done! You can read about the changes we've made on here: https://guide-on.gitbook.io/telltail/changes-done-by-install")
	return nil
}

func installCenter(authKey string) error {
	// TODO FIXME check for autohotkey
	// kill running and delete existing stuff if needed

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return cli.Exit("Cannot determine your home folder", exitCannotDetermineUserHomeDir)
	}
	baseBinLoc := filepath.Join(homeDir, ".local/share/telltail")

	{
		loc := filepath.Join(baseBinLoc, "telltail-center.exe")
		err, exitCode := downloadFile(
			"https://github.com/ajitid/telltail-center/releases/download/"+version+"/telltail-center-win-"+runtime.GOARCH+".exe",
			loc)
		if err != nil {
			return cli.Exit(err, exitCode)
		}
		markFileAsExecutableOnUnix(loc)
	}

	// write to local override file and tell user to open it and manually enter key there to avoid
	// and because they'll have the familiarity, they'll be able to update it as well. Revocation and expiration of key is quite common to happen
	// tell them what they can use to change auth key if they need to

	////// Success message
	fmt.Println("All done! You can read about the changes we've made on here: https://guide-on.gitbook.io/telltail/changes-done-by-install")
	return nil
}
