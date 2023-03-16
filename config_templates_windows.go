package main

import (
	"strings"
	"text/template"
)

type syncAhkCfgAttrs struct {
	BinDirectory, Tailnet, Device string
}

func getSyncAhkCfg() *template.Template {
	tmpl, err := template.New("sync-ahk-cfg-windows").Parse(strings.TrimSpace(`
#SingleInstance Force
#NoEnv  ; suggested by AHK
SendMode Input  ; suggested by AHK

SetWorkingDir {{.BinDirectory}}
RunWait telltail-sync.exe --url https://telltail.{{.Tailnet}} --device {{.Device}},,Hide
`))

	if err != nil {
		// panicking is fine as here, as user cannot report anything about it
		panic(err)
	}
	return tmpl
}

type centerAhkCfgAttrs struct {
	BinDirectory string
}

func getCenterAhkCfg() *template.Template {
	tmpl, err := template.New("sync-ahk-cfg-windows").Parse(strings.TrimSpace(`
#SingleInstance Force
#NoEnv  ; suggested by AHK
SendMode Input  ; suggested by AHK

SetWorkingDir {{.BinDirectory}}
RunWait telltail-center.exe,,Hide
`))

	if err != nil {
		// panicking is fine as here, as user cannot report anything about it
		panic(err)
	}
	return tmpl
}
