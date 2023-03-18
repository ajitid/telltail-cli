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

const (
	binPath      = ".local/share/telltail"
	startupPath  = "Library/LaunchAgents"
	scriptPrefix = "com.hemarkable."
)

func installSync(params installSyncParams) error {
	////// Check basic necessities exist
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return cli.Exit("Cannot determine your home folder", exitCannotDetermineUserHomeDir)
	}
	baseBinLoc := filepath.Join(homeDir, binPath)

	{
		if !cmdExists("launchctl") {
			return cli.Exit("We use launchctl/launchd to run services on boot. We cannot proceed if that is not available.", exitMissingDependency)
		}

		// acts like: systemctl disable telltail-sync --now
		cmd := exec.Command("launchctl", "unload", filepath.Join(homeDir, startupPath, scriptPrefix+"telltail-sync.plist"))
		cmd.Run()
	}

	////// Download and store clipnotify
	{
		loc := filepath.Join(baseBinLoc, "clipnotify")
		err, exitCode := downloadFile(
			"https://github.com/ajitid/clipnotify-for-desktop-os/releases/download/"+version+"/clipnotify-mac-"+runtime.GOARCH,
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
			"https://github.com/ajitid/telltail-sync/releases/download/"+version+"/telltail-sync-mac-"+runtime.GOARCH,
			loc)
		if err != nil {
			return cli.Exit(err, exitCode)
		}
		markFileAsExecutableOnUnix(loc)
	}

	////// Put bootup configuration and start the service
	// fmt.Println("Configuring for it load on boot...")
	{
		dir := filepath.Join(homeDir, startupPath)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Println("Unable to create folder", dir)
			return cli.Exit(err, exitDirNotModifiable)
		}

		tmpl := getSyncLaunchdCfg()
		f, err := os.Create(filepath.Join(dir, scriptPrefix+"telltail-sync.plist"))
		if err != nil {
			return cli.Exit("Cannot create service file for systemd", exitFileNotModifiable)
		}
		defer f.Close()
		err = tmpl.Execute(f, syncLaunchdCfgAttrs{
			Tailnet:      params.tailnet,
			Device:       params.device,
			BinDirectory: baseBinLoc,
		})
		if err != nil {
			return cli.Exit("Cannot write to service file for systemd", exitFileNotModifiable)
		}

		cmd := exec.Command("launchctl", "load", filepath.Join(dir, scriptPrefix+"telltail-sync.plist"))
		cmd.Run()
		cmd = exec.Command("launchctl", "start", scriptPrefix+"telltail-sync")
		cmd.Run()
	}

	////// Success message
	fmt.Println("All done! You can read about the changes we've made on here: https://guide-on.gitbook.io/telltail/changes-done-by-install")
	return nil
}

func installCenter(authKey string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return cli.Exit("Cannot determine your home folder", exitCannotDetermineUserHomeDir)
	}
	baseBinLoc := filepath.Join(homeDir, binPath)

	{
		if !cmdExists("launchctl") {
			return cli.Exit("We use systemctl/systemd to run services on boot. We cannot proceed if that is not available.", exitMissingDependency)
		}

		cmd := exec.Command("launchctl", "unload", filepath.Join(homeDir, startupPath, scriptPrefix+"telltail-center.plist"))
		cmd.Run()
	}

	{
		loc := filepath.Join(baseBinLoc, "telltail-center")
		err, exitCode := downloadFile(
			"https://github.com/ajitid/telltail-center/releases/download/"+version+"/telltail-center-mac-"+runtime.GOARCH,
			loc)
		if err != nil {
			return cli.Exit(err, exitCode)
		}
		markFileAsExecutableOnUnix(loc)
	}

	{
		dir := filepath.Join(homeDir, startupPath)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Println("Unable to create folder", dir)
			return cli.Exit(err, exitDirNotModifiable)
		}

		tmpl := getCenterLaunchdCfg()
		f, err := os.Create(filepath.Join(dir, scriptPrefix+"telltail-center.plist"))
		if err != nil {
			return cli.Exit("Cannot create service file for systemd", exitFileNotModifiable)
		}
		defer f.Close()
		err = tmpl.Execute(f, centerLaunchdCfgAttrs{
			BinDirectory: baseBinLoc,
			AuthKey:      authKey,
		})
		if err != nil {
			return cli.Exit("Cannot write to service file for systemd", exitFileNotModifiable)
		}

		cmd := exec.Command("launchctl", "load", filepath.Join(dir, scriptPrefix+"telltail-center.plist"))
		cmd.Run()
		cmd = exec.Command("launchctl", "start", scriptPrefix+"telltail-center")
		cmd.Run()
	}

	////// Success message
	fmt.Println("All done! You can read about the changes we've made on here: https://guide-on.gitbook.io/telltail/changes-done-by-install")
	return nil
}
