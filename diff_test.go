// +build wasm

package vdom_test

import (
	"reflect"
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"wasm-vdom"
)

func TestDiff(t *testing.T) {
	t.Run("NoPatches", func(t *testing.T) {
		old := newTrees(t, `<div><span>Hello</span><span>World</span></div>`)
		new := newTrees(t, `<div><span>Hello</span><span>World</span></div>`)
		if vdom.Diff(old, new) != nil {
			t.Error("expected no patches")
		}
	})

	t.Run("ReplaceInnerText", func(t *testing.T) {
		old := newTrees(t, `<div><span>Hello</span><span>World</span></div>`)
		new := newTrees(t, `<div><span>Hello</span><span>Universe</span></div>`)
		patches := vdom.Diff(old, new)
		if len(patches) != 1 {
			t.Fatal("expected one patch")
		}
		replace, ok := patches[0].(vdom.Replace)
		if !ok {
			t.Fatal("expected Replace patch")
		}
		if !reflect.DeepEqual(replace.Path, []int{0, 1, 0}) {
			t.Fatalf("unexpected path: %v", replace.Path)
		}
		if replace.Node.Type != html.TextNode {
			t.Fatal("unexpected replacement node type")
		}
		if replace.Node.Data != "Universe" {
			t.Fatal("unexpected replacement node data")
		}
	})

	t.Run("ReplaceNode", func(t *testing.T) {
		old := newTrees(t, `<div><span>Hello</span><span>World</span></div>`)
		new := newTrees(t, `<div><strong>Hello</strong><span>World</span></div>`)
		patches := vdom.Diff(old, new)
		if len(patches) != 1 {
			t.Fatal("expected one patch")
		}
		replace, ok := patches[0].(vdom.Replace)
		if !ok {
			t.Fatal("expected Replace patch")
		}
		if !reflect.DeepEqual(replace.Path, []int{0, 0}) {
			t.Fatalf("unexpected path: %v", replace.Path)
		}
		if replace.Node.Type != html.ElementNode {
			t.Fatal("unexpected replacement node type")
		}
		if replace.Node.DataAtom != atom.Strong {
			t.Fatal("unexpected replacement node data")
		}
	})

	t.Run("AddNewSpan", func(t *testing.T) {
		old := newTrees(t, `<div><span>Hello</span><span>World</span></div>`)
		new := newTrees(t, `<div><span>Hello</span><span>World</span><span>!</span></div>`)
		patches := vdom.Diff(old, new)
		if len(patches) != 1 {
			t.Fatal("expected one patch")
		}
		patch, ok := patches[0].(vdom.Append)
		if !ok {
			t.Fatal("expected Append patch")
		}
		if !reflect.DeepEqual(patch.Path, []int{0}) {
			t.Fatalf("unexpected path: %v", patch.Path)
		}
		if patch.Node.Type != html.ElementNode || patch.Node.DataAtom != atom.Span {
			t.Fatal("unexpected node type")
		}
		if patch.Node.FirstChild.Type != html.TextNode || patch.Node.FirstChild.Data != "!" {
			t.Fatal("unexpected node contents")
		}
	})

	t.Run("DeleteSpan", func(t *testing.T) {
		old := newTrees(t, `<div><span>Hello</span><span>World</span></div>`)
		new := newTrees(t, `<div><span>Hello</span></div>`)
		patches := vdom.Diff(old, new)
		if len(patches) != 1 {
			t.Fatal("expected one patch")
		}
		patch, ok := patches[0].(vdom.Remove)
		if !ok {
			t.Fatal("expected Append patch")
		}
		if !reflect.DeepEqual(patch.Path, []int{0, 1}) {
			t.Fatalf("unexpected path: %v", patch.Path)
		}
		if patch.Node.Type != html.ElementNode || patch.Node.DataAtom != atom.Span {
			t.Fatal("unexpected node type")
		}
		if patch.Node.FirstChild.Type != html.TextNode || patch.Node.FirstChild.Data != "World" {
			t.Fatal("unexpected node contents")
		}
	})

	t.Run("AddAttribute", func(t *testing.T) {
		old := newTrees(t, `<div><span>Hello</span><span>World</span></div>`)
		new := newTrees(t, `<div><span class="hello">Hello</span><span>World</span></div>`)
		patches := vdom.Diff(old, new)
		if len(patches) != 1 {
			t.Fatal("expected one patch")
		}
		patch, ok := patches[0].(vdom.AddAttribute)
		if !ok {
			t.Fatal("expected AddAttribute patch")
		}
		if !reflect.DeepEqual(patch.Path, []int{0, 0}) {
			t.Fatalf("unexpected path: %v", patch.Path)
		}
		if patch.Key != "class" || patch.Value != "hello" {
			t.Fatal("unexpected attribute")
		}
	})

	t.Run("ModifyAttribute", func(t *testing.T) {
		old := newTrees(t, `<div><span class="hola">Hello</span><span>World</span></div>`)
		new := newTrees(t, `<div><span class="hello">Hello</span><span>World</span></div>`)
		patches := vdom.Diff(old, new)
		if len(patches) != 1 {
			t.Fatal("expected one patch")
		}
		patch, ok := patches[0].(vdom.ModifyAttribute)
		if !ok {
			t.Fatal("expected ModifyAttribute patch")
		}
		if !reflect.DeepEqual(patch.Path, []int{0, 0}) {
			t.Fatalf("unexpected path: %v", patch.Path)
		}
		if patch.Key != "class" || patch.Value != "hello" {
			t.Fatal("unexpected attribute")
		}
	})

	t.Run("DeleteAttribute", func(t *testing.T) {
		old := newTrees(t, `<div><span class="hola">Hello</span><span>World</span></div>`)
		new := newTrees(t, `<div><span>Hello</span><span>World</span></div>`)
		patches := vdom.Diff(old, new)
		if len(patches) != 1 {
			t.Fatal("expected one patch")
		}
		patch, ok := patches[0].(vdom.DeleteAttribute)
		if !ok {
			t.Fatal("expected DeleteAttribute patch")
		}
		if !reflect.DeepEqual(patch.Path, []int{0, 0}) {
			t.Fatalf("unexpected path: %v", patch.Path)
		}
		if patch.Key != "class" {
			t.Fatal("unexpected attribute")
		}
	})
}

func newTrees(t *testing.T, src string) vdom.Trees {
	t.Helper()
	trees, err := vdom.NewTrees([]byte(src))
	if err != nil {
		t.Fatal(err)
	}
	return trees
}
