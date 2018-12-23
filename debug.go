// +build wasm

package vdom

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

const DEBUG = true

func debug(format string, a ...interface{}) {
	if !DEBUG {
		return
	}
	fmt.Printf(format+"\n", a...)
}

func nodeInfo(node *html.Node) string {
	switch node.Type {
	case html.TextNode:
		return "text node \"" + strings.TrimSpace(node.Data) + "\""
	case html.ElementNode:
		return "element node <" + node.Data + ">"
	default:
		return "<unknown node>"
	}
}
