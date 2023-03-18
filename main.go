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
	// TODO think about if I really need this
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	app := &cli.App{
		Name:    "telltail",
		Version: version,
		Usage:   "Telltail is a universal clipboard for text. And this CLI lets you configure it for your device.",
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
							err := guardArgsNonEmpty(cc, "auth-key")
							if err != nil {
								return err
							}
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
							err := guardArgsNonEmpty(cc, "tailnet", "device")
							if err != nil {
								return err
							}
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
				Usage:    "Start a program",
				Action: func(cc *cli.Context) error {
					return manageService(cc.Args().Get(0), startService)
				},
			},
			{
				Name:     "stop",
				Category: "Manage",
				Usage:    "Stop a program. (This will not stop it from running on system restart. Use uninstall for that.)",
				Action: func(cc *cli.Context) error {
					return manageService(cc.Args().Get(0), stopService)
				},
			},
			{
				Name:     "restart",
				Category: "Manage",
				Usage:    "Restart a program",
				Action: func(cc *cli.Context) error {
					return manageService(cc.Args().Get(0), restartService)
				},
			},
			{
				Name:     "edit",
				Category: "Manage",
				Usage:    "Edit config of a program",
				Subcommands: []*cli.Command{
					{
						Name: "center-auth-key",
						Action: func(cc *cli.Context) error {
							return editCenterAuthKey()
						},
					},
				},
			},
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
