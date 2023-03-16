package main

import (
	"log"
	"os"
	"runtime"

	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
)

// TODO
// rangeSupported =

func main() {
	// removes timestamp from `log` https://stackoverflow.com/a/48630122
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	app := &cli.App{
		Name:    "telltail",
		Version: version,
		Usage: "Telltail is a universal clipboard for text." +
			" To make it work, we need to configure it for your device, and this CLI helps you manage that.",
		Commands: []*cli.Command{
			{
				Name:  "install",
				Usage: "Install one of Telltail programs",
				Subcommands: []*cli.Command{
					{
						Name:  "center",
						Usage: "Installs Center. Only one device in your Tailscale network needs to install this, and the device should mostly be running when you work.",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "auth-key",
								Usage:    "A reusable, non-ephermal key you obtained from https://login.tailscale.com/admin/settings/keys",
								Required: true,
							},
						},
						Action: func(cc *cli.Context) error {
							switch runtime.GOOS {
							case "linux":
								return installCenterOnLinux(cc.String("auth-key"))
							default:
								return cli.Exit("This OS is not supported yet. If this is a desktop OS, please file an issue using `telltail file-issue`", exitUnsupportedOs)
							}
						},
					}, {
						Name:  "sync",
						Usage: "Installs Sync. If you want Ctrl+C and Ctrl+V (or Command+C and Command+V in macOS) to use universal clipboard, then you should install it on this device.",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "tailnet",
								Usage:    "You can find your tailnet name in here: https://login.tailscale.com/admin/dns",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "device",
								Usage:    "A unique name for this device. Tailscale would've assinged a name and an IP to this device. You could use that for example.",
								Required: true,
							},
						},
						Action: func(cc *cli.Context) error {
							switch runtime.GOOS {
							case "linux":
								return installSyncOnLinux(installSyncOnLinuxParams{tailnet: cc.String("tailnet"), device: cc.String("device")})
							default:
								return cli.Exit("This OS is not supported yet. If this is a desktop OS, please file an issue using `telltail file-issue`", exitUnsupportedOs)
							}
						},
					},
				},
			},
			{
				Name:  "uninstall",
				Usage: "Uninstall one of Telltail programs",
				Subcommands: []*cli.Command{
					{
						Name:  "center",
						Usage: "Uninstalls Center",
					}, {
						Name:  "sync",
						Usage: "Uninstalls Sync",
					},
				},
			},
			// TODO telltail edit center-auth-key >> makes a cheap call to systemctl --user edit telltail-center
			//
			// telltail sync stop/start/restart
			{
				Name:  "healthcheck",
				Usage: "TODO",
				// TODO also open filepath location for easy access for the user, even better if you highlight the files in the explorer
			},
			{
				Name:  "check-update",
				Usage: "TODO",
			},
			{
				Name:  "tldr",
				Usage: "TODO",
				// TODO also open filepath location for easy access for the user, even better if you highlight the files in the explorer
				// publish to TLDR npm as well
				// check `tldr screenkey` and `tldr tailscale`
			},
			{
				Name:  "guide",
				Usage: "Opens documentation guide in your browser",
				Action: func(cc *cli.Context) error {
					return browser.OpenURL("https://guide-on.gitbook.io/telltail")
				},
			},
			{
				Name:  "file-issue",
				Usage: "TODO",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
