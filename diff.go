// +build wasm

package vdom

import (
	"golang.org/x/net/html"
)

func Diff(old, new Trees) []Patcher {
	debug("== DIFFING ==")
	return diff(Tree{Children: old}, Tree{Children: new})
}

func diff(old, new Tree) (patches []Patcher) {
	if areDifferentNodes(old.Node, new.Node) {
		debug("Patch %v: replace %s with %s", new.Path, nodeInfo(old.Node), nodeInfo(new.Node))
		return append(patches, Replace{
			Path: new.Path,
			Node: new.Node,
		})
	}

	if old.Node != nil && new.Node != nil {
		oldAttrs, newAttrs := attrMap(old.Node.Attr), attrMap(new.Node.Attr)
		for key, newValue := range newAttrs {
			oldValue, ok := oldAttrs[key]
			if !ok {
				// add attribute
				patches = append(patches, AddAttribute{
					Path:  new.Path,
					Key:   key,
					Value: newValue,
				})
			} else if newValue != oldValue {
				// modify attribute
				patches = append(patches, ModifyAttribute{
					Path:  new.Path,
					Key:   key,
					Value: newValue,
				})
			}
		}
		for key, _ := range oldAttrs {
			if _, ok := newAttrs[key]; !ok {
				// delete attribute
				patches = append(patches, DeleteAttribute{
					Path: old.Path,
					Key:  key,
				})
			}
		}
	}

	if len(new.Children) > len(old.Children) {
		for _, newChild := range new.Children[len(old.Children):] {
			debug("Patch %v: append %s", new.Path, nodeInfo(newChild.Node))
			patches = append(patches, Append{
				Path: new.Path,
				Node: newChild.Node,
			})
		}
	} else if len(new.Children) < len(old.Children) {
		for _, oldChild := range old.Children[len(new.Children):] {
			debug("Patch %v: remove %s", oldChild.Path, nodeInfo(oldChild.Node))
			patches = append(patches, Remove{
				Path: oldChild.Path,
				Node: oldChild.Node,
			})
		}
	}

	for i := 0; i < min(len(new.Children), len(old.Children)); i++ {
		patches = append(patches, diff(old.Children[i], new.Children[i])...)
	}

	return patches
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func attrMap(attrs []html.Attribute) map[string]string {
	m := make(map[string]string, len(attrs))
	for _, attr := range attrs {
		k := attr.Key
		if attr.Namespace != "" {
			k = attr.Namespace + ":" + k
		}
		m[k] = attr.Val
	}
	return m
}

func areDifferentNodes(n1, n2 *html.Node) bool {
	if n1 == nil && n2 == nil {
		return false
	}
	if n1 == nil || n2 == nil {
		return true
	}
	if n1.Type != n2.Type {
		return true
	}
	if n1.DataAtom != n2.DataAtom {
		return true
	}
	if n1.Data != n2.Data {
		return true
	}
	n1, n2 = n1.FirstChild, n2.FirstChild
	for n1 != nil && n2 != nil {
		if areDifferentNodes(n1, n2) {
			return true
		}
		n1, n2 = n1.NextSibling, n2.NextSibling
	}
	return false
}
