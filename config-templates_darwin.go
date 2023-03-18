package main

import (
	"strings"
	"text/template"
)

type syncLaunchdCfgAttrs struct {
	BinDirectory, Tailnet, Device string
}

func getSyncLaunchdCfg() *template.Template {
	tmpl, err := template.New("sync-launchd-cfg-mac").Parse(strings.TrimSpace(`
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>

    <key>Label</key>
    <string>com.hemarkable.telltail-sync</string>

    <key>RunAtLoad</key>
    <true/>

    <key>WorkingDirectory</key>
    <string>{{.BinDirectory}}</string>

    <key>ProgramArguments</key>
    <array>
      <string>./telltail-sync</string>
      <string>--url</string>
      <string>https://telltail.{{.Tailnet}}</string>
      <string>--device</string>
      <string>{{.Device}}</string>
    </array>

  </dict>
</plist>
`))

	if err != nil {
		// panicking is fine as here, as user cannot report anything about it
		panic(err)
	}
	return tmpl
}

type centerLaunchdCfgAttrs struct {
	BinDirectory, AuthKey string
}

func getCenterLaunchdCfg() *template.Template {
	tmpl, err := template.New("center-systemd-cfg-linux").Parse(strings.TrimSpace(`
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>

    <key>Label</key>
    <string>com.hemarkable.telltail-center</string>

    <key>RunAtLoad</key>
    <true/>

    <key>EnvironmentVariables</key>
    <dict>
      <key>TS_AUTHKEY</key>
      <string>{{.AuthKey}}</string>
    </dict>

    <key>WorkingDirectory</key>
    <string>{{.BinDirectory}}</string>

    <key>ProgramArguments</key>
    <array>
      <string>./telltail-center</string>
    </array>

  </dict>
</plist>
`))

	if err != nil {
		// panicking is fine as here, as user cannot report anything about it
		panic(err)
	}
	return tmpl
}
