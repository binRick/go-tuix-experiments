package main

import (
	"io/ioutil"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	rootDir = "."
	root    = tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed)
	tree = tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)
)

func add(target *tview.TreeNode, path string) {
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

func tree_main() {
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
	if err := tview.NewApplication().SetRoot(tree, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
