//go:build (windows && amd64) || (linux && amd64)

package main

import (
	"log"
	"os"

	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
)

func main() {
	// removes timestamp from `log` https://stackoverflow.com/a/48630122
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	app := &cli.App{
		Name:    "telltail-cli",
		Version: version,
		Usage: "Telltail is a universal clipboard for text." +
			" To make it work, we configure it for your device, and this CLI helps you to do that.",
		Commands: []*cli.Command{
			{
				Name:     "install",
				Category: "Install",
				Usage:    "Install one of Telltail programs",
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
							return installCenter(cc.String("auth-key"))
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
							return installSync(installSyncParams{tailnet: cc.String("tailnet"), device: cc.String("device")})
						},
					},
				},
			},
			{
				Name:     "uninstall",
				Category: "Install",
				Usage:    "Uninstall one of Telltail programs",
				Subcommands: []*cli.Command{
					{
						Name:  "center",
						Usage: "Uninstalls Center",
						Action: func(cc *cli.Context) error {
							return uninstallCenter()
						},
					}, {
						Name:  "sync",
						Usage: "Uninstalls Sync",
						Action: func(cc *cli.Context) error {
							return uninstallSync()
						},
					},
				},
			},
			{
				Name:     "start",
				Category: "Manage",
				Usage:    "Start a service",
				Action: func(cc *cli.Context) error {
					return manageService(cc.Args().Get(0), startService)
				},
			},
			{
				Name:     "stop",
				Category: "Manage",
				Usage:    "Stop a service. (This will not stop it from running on system restart. Use uninstall for that.)",
				Action: func(cc *cli.Context) error {
					return manageService(cc.Args().Get(0), stopService)
				},
			},
			{
				Name:     "restart",
				Category: "Manage",
				Usage:    "Restart a service",
				Action: func(cc *cli.Context) error {
					return manageService(cc.Args().Get(0), restartService)
				},
			},
			{
				Name:     "edit",
				Category: "Manage",
				Subcommands: []*cli.Command{
					{
						Name: "center-auth-key",
					},
				},
			},
			// TODO telltail edit center-auth-key >> makes a cheap call to systemctl --user edit telltail-center
			//
			// telltail sync stop/start/restart
			{
				Name:     "guide",
				Category: "Help",
				Usage:    "Opens documentation guide in your browser",
				Action: func(cc *cli.Context) error {
					return browser.OpenURL("https://guide-on.gitbook.io/telltail")
				},
			},
			{
				Name:     "help",
				Category: "Help",
				Usage:    "Opens this help",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
