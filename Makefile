.PHONY: all clean start package

$(shell mkdir -p bin electron/src/bin >/dev/null 2>&1 || true)

start: electron/.electron
	cd electron && npm start

electron/.electron: bin/main.wasm bin/window.wasm bin/overlay.wasm cmd/main/main.js cmd/main/preload.js
	cp -r web/static/data/* electron/src
	cp cmd/main/main.js cmd/main/preload.js electron/src
	cp bin/main.wasm bin/window.wasm bin/overlay.wasm electron/src/bin
	touch $@

bin/main.wasm: $(shell find cmd/main -name \*.go) $(shell find cmd/main -name \*.js) $(shell find internal -type f) $(shell find pkg -type f)
	GOOS=js GOARCH=wasm go build -tags electron -o $@ ./cmd/main


bin/window.wasm: $(shell find cmd/window -type f) $(shell find internal -type f) $(shell find pkg -type f)
	go run -mod=vendor github.com/vugu/vugu/cmd/vugugen -s -r -skip-go-mod -skip-main cmd/window
	GOOS=js GOARCH=wasm go build -tags web -o $@ ./cmd/window

bin/overlay.wasm: $(shell find cmd/overlay -type f) $(shell find internal -type f) $(shell find pkg -type f)
	go run -mod=vendor github.com/vugu/vugu/cmd/vugugen -s -r -skip-go-mod -skip-main cmd/overlay
	GOOS=js GOARCH=wasm go build -tags web -o $@ ./cmd/overlay

clean:
	rm -rf bin/* build/* electron/.electron electron/src/* electron/out $(shell find . -name 0_components_vgen.go)