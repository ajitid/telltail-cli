package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/urfave/cli/v2"
)

func installSync(params installSyncParams) error {
	////// Check basic necessities exist
	// FIXME TODO check if autohotkey exists
	// remove existing clipnotify folder

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return cli.Exit("Cannot determine your home folder", exitCannotDetermineUserHomeDir)
	}
	baseBinLoc := filepath.Join(homeDir, ".local/share/telltail")

	////// Download and store clipnotify
	{
		loc := filepath.Join(baseBinLoc, "clipnotify.zip")
		err, exitCode := downloadFile(
			"https://github.com/ajitid/clipnotify-for-desktop-os/releases/download/"+version+"/clipnotify-win-"+runtime.GOARCH+".zip",
			loc)
		if err != nil {
			return cli.Exit(err, exitCode)
		}
		// FIXME TODO extract it using tar.exe, see https://superuser.com/a/1473255
		// Also, delete existing folder if exists
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

	////// Put bootup configuration
	{
		// also change script such that triggering it restarts the service, right now it pops-up an alert
	}

	////// Start the service
	{
		// ahk can be started from commandline
		// cmd := exec.Command("systemctl", "--user", "daemon-reload")
		// cmd.Output()
		// cmd = exec.Command("systemctl", "--user", "enable", "telltail-sync", "--now")
		// cmd.Output()
	}

	// TODO handle failures:
	// systemctl status will give status code 3 if:
	// - service is stopped
	// - start the service fails
	// **Do note that** status code is 3 by telltail-center, not by telltail-sync (as I tested w/ telltail-center). It could be different for Sync.
	// so yeah, that ain't a way to distinguish. It also prints logs from journalctl, which we can use though:
	// On normal stop:
	// Mar 16 14:47:24 sd systemd[2235]: Stopped telltail.service - Telltail server.
	// Mar 16 14:47:24 sd systemd[2235]: telltail.service: Consumed 3min 39.217s CPU time.
	// On failure stop:
	// Mar 16 14:47:31 sd systemd[2235]: telltail.service: Main process exited, code=exited, status=203/EXEC
	// Mar 16 14:47:31 sd systemd[2235]: telltail.service: Failed with result 'exit-code'.
	//
	// We probably could also be able to pass flags and get the active statuses:
	//	Active: inactive (dead) // normal stop
	//	Active: failed (Result: exit-code) since Thu 2023-03-16 14:47:31 IST; 4s ago // failure stop

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
