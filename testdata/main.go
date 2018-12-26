package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
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
	t     *template.Template
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
	for _, patch := range vdom.Diff(trees, newTrees) {
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

func deleteItem(args []js.Value) {
	i := args[0].Int()
	state.Items = append(state.Items[:i], state.Items[i+1:]...)
	fmt.Printf("items = %v\n", state.Items)
	render()
}

func loadTemplate() *template.Template {
	resp, err := http.Get("http://localhost:8080/index.tmpl")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return template.Must(template.New("").Parse(string(data)))
}

func registerGlobals() {
	for _, f := range []struct {
		name string
		impl func([]js.Value)
	}{
		// Add callbacks here.
		{name: "addNewItem", impl: addNewItem},
		{name: "deleteItem", impl: deleteItem},
	} {
		js.Global().Set(f.name, js.NewCallback(f.impl))
	}
}

func main() {
	registerGlobals()
	t = loadTemplate()

	render()

	select {}
}
