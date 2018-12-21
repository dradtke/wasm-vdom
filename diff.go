package vdom

import (
	"golang.org/x/net/html"
)

func Diff(old, new Tree) (patches []Patcher) {
	if areDifferentNodes(old.Node, new.Node) {
		return append(patches, Replace{
			Path: old.Path,
			Node: new.Node,
		})
	}

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

	if len(new.Children) > len(old.Children) {
		for _, newChild := range new.Children[len(old.Children):] {
			patches = append(patches, Append{
				Path: new.Path,
				Node: newChild.Node,
			})
		}
	} else if len(new.Children) < len(old.Children) {
		for _, oldChild := range old.Children[len(new.Children):] {
			patches = append(patches, Remove{
				Path: oldChild.Path,
				Node: oldChild.Node,
			})
		}
	}

	for i := 0; i < min(len(new.Children), len(old.Children)); i++ {
		patches = append(patches, Diff(old.Children[i], new.Children[i])...)
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

type Patcher interface{}

type Replace struct {
	Path []int
	Node *html.Node
}

type Append struct {
	Path []int
	Node *html.Node
}

type Remove struct {
	Path []int
	Node *html.Node
}

type AddAttribute struct {
	Path       []int
	Key, Value string
}

type ModifyAttribute struct {
	Path       []int
	Key, Value string
}

type DeleteAttribute struct {
	Path []int
	Key  string
}

func areDifferentNodes(n1, n2 *html.Node) bool {
	return n1.Type != n2.Type || n1.DataAtom != n2.DataAtom || n1.Data != n2.Data
}
