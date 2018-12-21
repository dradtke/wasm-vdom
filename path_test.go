// +build wasm

package vdom

import (
	"reflect"
	"testing"
)

func TestAddToPath(t *testing.T) {
	path := []int{1, 2, 3}
	if !reflect.DeepEqual(addToPath(path, 4), []int{1, 2, 3, 4}) {
		t.Error("addToPath failed")
	}
}
