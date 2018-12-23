// +build wasm

package vdom

import (
	"bytes"
	"errors"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Trees []Tree

type Tree struct {
	Node     *html.Node
	Children []Tree
	// Path is the list of indices needed to access this node from the root of the tree.
	Path []int
}

func NewTrees(src []byte) (Trees, error) {
	roots, err := parse(src)
	if err != nil {
		return nil, errors.New("vdom.NewTree: " + err.Error())
	}

	trees := make(Trees, 0, len(roots))
	var i int
	for _, root := range roots {
		if root.Type == html.ElementNode {
			trees = append(trees, NewTree(root, i))
			i++
		}
	}

	return trees, nil
}

func NewTree(root *html.Node, topIndex int) Tree {
	var f func(*html.Node, []int) Tree
	f = func(node *html.Node, path []int) Tree {
		var children []Tree
		var i int
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.ElementNode {
				children = append(children, f(child, addToPath(path, i)))
				i++
			}
		}
		return Tree{Node: node, Children: children, Path: path}
	}
	return f(root, []int{topIndex})
}

// addToPath creates a new copy of the path slice with i added to the end.
func addToPath(path []int, i int) []int {
	new := make([]int, len(path)+1)
	new[copy(new, path)] = i
	return new
}

func parse(src []byte) ([]*html.Node, error) {
	roots, err := html.ParseFragment(bytes.NewReader(bytes.TrimSpace(src)), &html.Node{
		Type:     html.ElementNode,
		Data:     "body",
		DataAtom: atom.Body,
	})
	if err != nil {
		return nil, err
	}
	return roots, nil
}
