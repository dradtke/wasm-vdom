package main

import (
	"bytes"
	"html/template"
	"syscall/js"

	"wasm-vdom"
)

type State struct {
	Items []string
}

var (
	body  = js.Global().Get("document").Call("querySelector", "body")
	trees vdom.Trees
	state = State{}
	t     = template.Must(template.New("").Parse(`
		<p>To-Do List:</p>
		<ul>
		  {{range .Items}}
		    <li>{{.}}</li>
		  {{else}}
		    <em>Nothing here yet.</em>
		  {{end}}
		</ul>
		<br>
		<p>
		  <input type="text" id="newItem">&nbsp;<button type="button" onclick="addNewItem()">Add</button>
		</p>
	`))
)

func render() {
	var buf bytes.Buffer
	if err := t.Execute(&buf, state); err != nil {
		panic(err)
	}
	newTrees, err := vdom.NewTrees(buf.Bytes())
	if err != nil {
		panic("error parsing new tree: " + err.Error())
	}
	/*
		for _, tree := range newTrees {
			fmt.Printf("%v %v: %s\n", tree.Path, tree.Node.Type, tree.Node.Data)
		}
	*/
	patches := vdom.Diff(trees, newTrees)
	for _, patch := range patches {
		if err := patch.Patch(body); err != nil {
			panic("error applying patch: " + err.Error())
		}
	}
	trees = newTrees
}

func addNewItem(args []js.Value) {
	input := js.Global().Get("document").Call("getElementById", "newItem")
	value := input.Get("value").String()

	state.Items = append(state.Items, value)
	render()
}

func main() {
	js.Global().Set("addNewItem", js.NewCallback(addNewItem))
	render()

	select {}
}
