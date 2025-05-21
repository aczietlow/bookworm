package main

import (
	"fmt"

	"github.com/rivo/tview"
)

func commandDebug(conf *config, args ...string) (string, error) {

	authorData, err := conf.apiClient.DebugQuery()
	if err != nil {
		return "", err
	}

	output := "Debugging"
	output += fmt.Sprintf("%T", authorData)

	return output, nil
}

func viewDebug(conf *config) tview.Primitive {
	return tview.NewBox().SetTitle("Debug").SetBorder(true)
}
