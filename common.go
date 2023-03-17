package main

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
