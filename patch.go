// +build wasm

package vdom

import (
	"bytes"
	"errors"
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
	newNode, err := createNodeValue(a.Node)
	if err != nil {
		return err
	}
	_, target := traverse(root, a.Path)
	target.Call("appendChild", newNode)
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
	parent, target = js.Null(), start
	for _, i := range path {
		parent = target
		target = target.Get("children").Index(i)
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
