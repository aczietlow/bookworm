package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func commandDebug(conf *config, args ...string) ([]byte, error) {

	authorData, err := conf.apiClient.DebugQuery()
	if err != nil {
		return nil, err
	}

	output := "Debugging\n"
	output += fmt.Sprintf("Author: %v\n", authorData)

	return []byte(output), nil
}

func debugView(conf *config) tview.Primitive {
	return tview.NewBox().SetTitle("Debug").SetBorder(true)
}

func resultTextDebug(conf *config) tview.Primitive {
	app := conf.tui.app
	results := tview.NewTextView().
		SetChangedFunc(func() {
			app.Draw()
		})
	results.SetTitle("Debug - Test").SetBorder(true)
	return results
}

func debugResultView(conf *config) tview.Primitive {
	// app := conf.tui.app

	rootDir := "."
	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	// A helper function which adds the files and directories of the given path
	// to the given target node.
	add := func(target *tview.TreeNode, path string) {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			node := tview.NewTreeNode(file.Name()).
				SetReference(filepath.Join(path, file.Name())).
				SetSelectable(file.IsDir())
			if file.IsDir() {
				node.SetColor(tcell.ColorGreen)
			}
			target.AddChild(node)
		}
	}

	// Add the current directory to the root node.
	add(root, rootDir)

	// If a directory was selected, open it.
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}
		children := node.GetChildren()
		if len(children) == 0 {
			// Load and show files in this directory.
			path := reference.(string)
			add(node, path)
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})

	return tree
}

func newDebugCommandView(conf *config) *commandView {
	return &commandView{
		view:             debugView(conf),
		updateView:       nil,
		resultView:       debugResultView(conf),
		updateResultView: nil,
	}
}
