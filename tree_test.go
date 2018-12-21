// +build wasm

package vdom_test

import (
	"reflect"
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"wasm-vdom"
)

func TestNewTree(t *testing.T) {
	trees, err := vdom.NewTrees([]byte(`<div><span>Hello</span><span>World</span></div>`))
	if err != nil {
		t.Fatal(err)
	}
	tree := trees[0]

	if tree.Node.Type != html.ElementNode || tree.Node.DataAtom != atom.Div {
		t.Fatal("expected root to be div")
	}
	if len(tree.Children) != 2 {
		t.Fatal("expected root to have two children")
	}
	if !reflect.DeepEqual(tree.Path, []int{0}) {
		t.Fatal("expected root to have path of [0]")
	}

	helloTree := tree.Children[0]
	if helloTree.Node.Type != html.ElementNode || helloTree.Node.DataAtom != atom.Span {
		t.Fatal("expected hello to be span")
	}
	if len(helloTree.Children) != 1 {
		t.Fatal("expected hello to have one child")
	}
	if !reflect.DeepEqual(helloTree.Path, []int{0, 0}) {
		t.Fatalf("expected hello to have path of [0 0] (was: %v)", helloTree.Path)
	}
	helloText := helloTree.Children[0]
	if helloText.Node.Type != html.TextNode || helloText.Node.Data != "Hello" {
		t.Fatal("expected hello to say 'Hello'")
	}

	worldTree := tree.Children[1]
	if worldTree.Node.Type != html.ElementNode || worldTree.Node.DataAtom != atom.Span {
		t.Fatal("expected world to be span")
	}
	if len(worldTree.Children) != 1 {
		t.Fatal("expected world to have one child")
	}
	if !reflect.DeepEqual(worldTree.Path, []int{0, 1}) {
		t.Fatalf("expected hello to have path of [0 1] (was: %v)", worldTree.Path)
	}
	worldText := worldTree.Children[0]
	if worldText.Node.Type != html.TextNode || worldText.Node.Data != "World" {
		t.Fatal("expected world to say 'World'")
	}
}
