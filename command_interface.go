package main

import "github.com/rivo/tview"

type command interface {
	callback(*config, ...string) (string, error)
	view(*config) tview.Primitive
	result(*config) tview.Primitive
}
