package vdom

import (
	"errors"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Tree struct {
	Node     *html.Node
	Children []Tree
	// Path is the list of indices needed to access this node from the root of the tree.
	Path []int
}

func NewTree(src string) (Tree, error) {
	root, err := parse(src)
	if err != nil {
		return Tree{}, errors.New("vdom.NewTree: " + err.Error())
	}

	var f func(*html.Node, []int) Tree
	f = func(node *html.Node, path []int) Tree {
		var children []Tree
		var i int
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			children = append(children, f(child, addToPath(path, i)))
			i++
		}
		return Tree{Node: node, Children: children, Path: path}
	}
	return f(root, nil), nil
}

func addToPath(path []int, i int) []int {
	new := make([]int, len(path)+1)
	new[copy(new, path)] = i
	return new
}

func parse(src string) (*html.Node, error) {
	roots, err := html.ParseFragment(strings.NewReader(src), &html.Node{
		Type:     html.ElementNode,
		Data:     "body",
		DataAtom: atom.Body,
	})
	if err != nil {
		return nil, err
	}
	if len(roots) != 1 {
		return nil, errors.New("expected only one root")
	}
	return roots[0], nil
}
