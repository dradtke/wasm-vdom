// +build wasm

package vdom

import (
	"bytes"
	"errors"
	"fmt"
	"syscall/js"

	"golang.org/x/net/html"
)

func createNodeValue(node *html.Node) (js.Value, error) {
	switch node.Type {
	case html.ElementNode:
		rendered, err := renderNode(node)
		if err != nil {
			return js.Null(), err
		}
		parser := js.Global().Get("DOMParser").New()
		doc := parser.Call("parseFromString", rendered, "text/html")
		return doc.Call("querySelector", "body").Get("firstChild"), nil
	case html.TextNode:
		return js.Global().Get("document").Call("createTextNode", node.Data), nil
	default:
		return js.Null(), errors.New("createNodeValue: unsupported node type")
	}
}

type Patcher interface {
	Patch(root js.Value) error
}

type Replace struct {
	Path []int
	Node *html.Node
}

func (r Replace) Patch(root js.Value) error {
	fmt.Printf("applying Replace %v %s\n", r.Path, r.Node.Data)
	newNode, err := createNodeValue(r.Node)
	if err != nil {
		return err
	}
	parent, oldNode := traverse(root, r.Path)
	parent.Call("replaceChild", newNode, oldNode)
	return nil
}

type Append struct {
	Path []int
	Node *html.Node
}

func (a Append) Patch(root js.Value) error {
	fmt.Printf("applying Append %v %s\n", a.Path, a.Node.Data)
	newNode, err := createNodeValue(a.Node)
	if err != nil {
		return err
	}
	parent, target := traverse(root, a.Path)
	if parent == js.Null() {
		root.Call("appendChild", newNode)
	} else {
		parent.Call("insertBefore", newNode, target.Get("nextSibling"))
	}
	return nil
}

type Remove struct {
	Path []int
	Node *html.Node
}

func (r Remove) Patch(root js.Value) error {
	return errors.New("remove not implemented")
}

type AddAttribute struct {
	Path       []int
	Key, Value string
}

func (a AddAttribute) Patch(root js.Value) error {
	return errors.New("add attribute not implemented")
}

type ModifyAttribute struct {
	Path       []int
	Key, Value string
}

func (m ModifyAttribute) Patch(root js.Value) error {
	return errors.New("modify attribute not implemented")
}

type DeleteAttribute struct {
	Path []int
	Key  string
}

func (d DeleteAttribute) Patch(root js.Value) error {
	return errors.New("delete attribute not implemented")
}

func traverse(start js.Value, path []int) (parent, target js.Value) {
	// If the very first child node is a text element, increment the first
	// path value in order to skip over it.
	if len(path) > 0 && start.Get("childNodes").Index(0).Get("nodeType").Int() == 3 {
		path[0]++
	}

	parent, target = js.Null(), start
	for _, i := range path {
		parent = target
		target = target.Get("childNodes").Index(i)
	}
	return parent, target
}

func renderNode(node *html.Node) (string, error) {
	var buf bytes.Buffer
	if err := html.Render(&buf, node); err != nil {
		return "", err
	}
	return buf.String(), nil
}
