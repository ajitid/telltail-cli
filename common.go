package main

import "github.com/urfave/cli/v2"

func contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

var validServices = []string{"sync", "center"}

type serviceAction int

const (
	startService serviceAction = iota
	stopService
	restartService
)

func guardArgsNonEmpty(cc *cli.Context, args ...string) error {
	for _, arg := range args {
		if len(cc.String(arg)) == 0 {
			return cli.Exit(arg+" cannot be empty", 1)
		}
	}
	return nil
}
