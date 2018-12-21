build:
	GOOS=js GOARCH=wasm go build

test:
	GOOS=js GOARCH=wasm go test -exec="$$(go env GOROOT)/misc/wasm/go_js_wasm_exec"
