package main

import (
	"os"
	"strings"
	"text/template"
)

/*
➜ which pnpm                                                                                                                                                                            (base)
/home/ajitid/.local/share/pnpm/pnpm
so store your executables in there too
*/

type SyncSystemdCfgLinuxX11Attrs struct {
	Tailnet string
	Device  string
}

// read it as: get Telltail Sync's systemd config for "Linux on X11"
func getSyncSystemdCfgLinuxX11(attrs SyncSystemdCfgLinuxX11Attrs) {
	tmpl, err := template.New("sync-systemd-cfg-linux").Parse(strings.TrimSpace(`
[Unit]
Description=Telltail Sync
Wants=network-online.target
After=network-online.target
After=graphical.target

[Service]
Type=simple
ExecStart=/home/ajitid/playground/telltail-sync/telltail-sync --url https://telltail.{{.Tailnet}} --device {{.Device}}
Environment=DISPLAY=:0

[Install]
WantedBy=default.target
`))
	// ^ TODO in ExecStart, decide on a location

	if err != nil {
		panic(err)
	}
	// TODO don't std out but return the output so we can write it to a file
	err = tmpl.Execute(os.Stdout, attrs)
	if err != nil {
		panic(err)
	}
}

type CenterSystemdCfgLinuxX11Attrs struct {
	AuthKey string
}

func getCenterSystemdCfgLinuxX11(attrs SyncSystemdCfgLinuxX11Attrs) {
	tmpl, err := template.New("center-systemd-cfg-linux").Parse(strings.TrimSpace(`
[Unit]
Description=Telltail Center
Wants=network-online.target
After=network-online.target

[Service]
Type=simple
ExecStart=/home/ajitid/playground/telltail/telltail
; Environment=TS_AUTHKEY={{.AuthKey}}
; ^ here's how you can safely add it ajit w/o pushing it to your dotfiles which are public: https://serverfault.com/a/413408
; systemctl --user edit telltail << prefer this instead

[Install]
WantedBy=default.target
`))
	// ^ TODO in ExecStart, decide on a location

	if err != nil {
		panic(err)
	}
	// TODO don't std out but return the output so we can write it to a file
	err = tmpl.Execute(os.Stdout, attrs)
	if err != nil {
		panic(err)
	}
}
