package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func cmdExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// taken from https://stackoverflow.com/questions/11692860/how-can-i-efficiently-download-a-large-file-using-go
func downloadFile(url, toLocation string) (error, int) {
	dirName := filepath.Dir(toLocation)
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		fileName := filepath.Base(toLocation)
		log.Println("Unable to create folder", dirName, "for file", fileName)
		return err, exitDirNotCreatable
	}

	// TODO check if it can successfully override existing file
	out, err := os.Create(toLocation)
	if err != nil {
		return err, exitFileNotWriteable
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err, exitUrlNotDownloadable
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s for %s", resp.Status, url), exitUrlNotDownloadable
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err, exitFileNotWriteable
	}

	return nil, 0
}

func markFileAsExecutableOnUnix(fullPath string) {
	cmd := exec.Command("chmod", "+x", fullPath)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println("cannnot mark file as executable:", fullPath)
		panic(err)
	}
}

type installSyncOnLinuxParams struct {
	tailnet, device string
}

func installSyncOnLinux(params installSyncOnLinuxParams) error {
	////// Check basic necessities exist
	// fmt.Println("Checking requirments...") <<< TODO this is useless and bearing for the user. Show them a progress bar and how much time is remaining instead
	// check if system is x11, https://github.com/atotto/clipboard has ways to indentify it
	env := os.Getenv("XDG_SESSION_TYPE")
	if env != "x11" {
		return cli.Exit("Sync cannot be installed on a non-X11 Linux", exitUnsupportedOsVariant)
	}

	{
		if !cmdExists("systemctl") {
			return cli.Exit("We use systemctl/systemd to run services on boot. We cannot proceed if that is not available.", exitMissingDependency)
		}

		// it'll fail if the systemd config is not present, which is fine as well, no need to panic
		// doing this + the fact writing a file overrides the existing one will make the `install` idempotent
		cmd := exec.Command("systemctl", "--user", "disable", "telltail-sync", "--now")
		cmd.Output()
	}

	if !cmdExists("xsel") && !cmdExists("xclip") {
		return cli.Exit("Either install `xsel` or `xclip` from your package manager first", exitMissingDependency)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return cli.Exit("Cannot determine your home folder", exitCannotDetermineUserHomeDir)
	}
	baseBinLoc := filepath.Join(homeDir, ".local/share/telltail")

	////// Download ans store clipnotify
	{
		loc := filepath.Join(baseBinLoc, "clipnotify")
		err, exitCode := downloadFile(
			"https://github.com/ajitid/clipnotify-for-desktop-os/releases/download/"+version+"/clipnotify-linux-x11",
			loc)
		if err != nil {
			return cli.Exit(err, exitCode)
		}
		markFileAsExecutableOnUnix(loc)
	}

	////// Download and store the telltail-sync
	// fmt.Println("Downloading files...")
	{
		loc := filepath.Join(baseBinLoc, "telltail-sync")
		err, exitCode := downloadFile(
			"https://github.com/ajitid/telltail-sync/releases/download/"+version+"/telltail-sync-linux",
			loc)
		if err != nil {
			return cli.Exit(err, exitCode)
		}
		markFileAsExecutableOnUnix(loc)
	}

	////// Put bootup configuration
	// fmt.Println("Configuring for it load on boot...")
	{
		dir := filepath.Join(homeDir, ".config/systemd/user")
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Println("Unable to create folder", dir)
			return cli.Exit(err, exitDirNotCreatable)
		}

		tmpl := getSyncSystemdCfgLinuxX11()
		f, err := os.Create(filepath.Join(dir, "telltail-sync.service"))
		if err != nil {
			return cli.Exit("Cannot create systemd file", exitFileNotWriteable)
		}
		err = tmpl.Execute(f, syncSystemdCfgLinuxX11Attrs{
			Tailnet:      params.tailnet,
			Device:       params.device,
			BinDirectory: baseBinLoc,
		})
		if err != nil {
			return cli.Exit("Cannot write to systemd file", exitFileNotWriteable)
		}
	}

	////// Start the service
	{
		cmd := exec.Command("systemctl", "--user", "daemon-reload")
		cmd.Output()
		cmd = exec.Command("systemctl", "--user", "enable", "telltail-sync", "--now")
		cmd.Output()
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
	fmt.Println("All done! You can read about the changes we've made on here https://guide-on.gitbook.io/telltail/changes-done-by-install")
	return nil
}

func installCenterOnLinux() {
	// tell them what they can use to change auth key if they need to
}
