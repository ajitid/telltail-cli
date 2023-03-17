package main

import (
	"strings"
	"text/template"
)

/*
➜ which pnpm                                                                                                                                                                            (base)
/home/ajitid/.local/share/pnpm/pnpm
so store your executables in there too
*/

type syncSystemdCfgLinuxX11Attrs struct {
	BinDirectory, Tailnet, Device string
}

// read it as: get Telltail Sync's systemd config for "Linux on X11"
func getSyncSystemdCfgLinuxX11() *template.Template {
	tmpl, err := template.New("sync-systemd-cfg-linux").Parse(strings.TrimSpace(`
[Unit]
Description=Telltail Sync
Wants=network-online.target
After=network-online.target
After=graphical.target

[Service]
Type=simple
WorkingDirectory={{.BinDirectory}}
ExecStart={{.BinDirectory}}/telltail-sync --url https://telltail.{{.Tailnet}} --device {{.Device}}
Environment=DISPLAY=:0

[Install]
WantedBy=default.target
`))

	if err != nil {
		// panicking is fine as here, as user cannot report anything about it
		panic(err)
	}
	return tmpl
}

type centerSystemdCfgLinuxAttrs struct {
	BinDirectory string
}

func getCenterSystemdCfgLinux() *template.Template {
	tmpl, err := template.New("center-systemd-cfg-linux").Parse(strings.TrimSpace(`
[Unit]
Description=Telltail Center
Wants=network-online.target
After=network-online.target

[Service]
Type=simple
ExecStart={{.BinDirectory}}/telltail-center

[Install]
WantedBy=default.target
`))

	if err != nil {
		// panicking is fine as here, as user cannot report anything about it
		panic(err)
	}
	return tmpl
}

type centerSystemdOverrideCfgLinuxAttrs struct {
	AuthKey string
}

func getCenterSystemdOverrideCfgLinux() *template.Template {
	tmpl, err := template.New("center-systemd-cfg-linux").Parse(strings.TrimSpace(`
[Service]
Environment=TS_AUTHKEY={{.AuthKey}}
`))

	if err != nil {
		// panicking is fine as here, as user cannot report anything about it
		panic(err)
	}
	return tmpl
}
