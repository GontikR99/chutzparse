.PHONY: all clean full-clean start package npm-install

$(shell mkdir -p bin electron/src/bin build>/dev/null 2>&1 || true)

start: build/electron-built build/npm-installed
	cd electron && npm start --arch=ia32

package: build/electron-built build/npm-installed
	cd electron && npm run make --arch=ia32
	find electron/out/make -name \*.exe -exec cp \{\} bin \;

build/npm-installed:
	# FIXME: workaround until electron/windows-installer includes 7z-ia32.exe.  See issues:
	#    https://github.com/electron/windows-installer/issues/378
	#    https://github.com/electron/windows-installer/issues/386
	cd electron && npm install --arch=ia32 || true
	cd electron && npm install --arch=ia32 --ignore-scripts
	curl -L https://github.com/electron/windows-installer/raw/b2380345e8fe1ad7716108b10b552d75e6fad0b7/vendor/7z-ia32.dll -o electron/node_modules/electron-winstaller/vendor/7z.dll
	curl -L https://github.com/electron/windows-installer/raw/b2380345e8fe1ad7716108b10b552d75e6fad0b7/vendor/7z-ia32.exe -o electron/node_modules/electron-winstaller/vendor/7z.exe
	touch $@

build/electron-built: bin/main.wasm bin/window.wasm bin/overlay.wasm cmd/main/main.js cmd/main/preload.js
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
	rm -rf bin/* build/electron-built electron/.electron electron/src/* electron/out $(shell find . -name 0_components_vgen.go)

full-clean: clean
	rm -rf electron/node_modules/* build/npm-installed