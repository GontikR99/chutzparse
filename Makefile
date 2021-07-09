.PHONY: all clean full-clean start package npm-install test help

help:
	@echo "Supported targets are:"
	@echo "    help        Display this help message"
	@echo "    clean       Remove built go artifacts (wasm and vugu stuff)"
	@echo "    fullclean  Remove all built artifacts, including node_modules"
	@echo "    package     Create a redistributable package in bin/"
	@echo "    start       Build and start in development mode"
	@echo "    test        Run go unit tests"

# Build a package for distribution
package: build/electron-built build/npm-installed
	cd electron && npm run make --arch=ia32
	find electron/out/make -name \*.exe -exec cp \{\} bin \;

# Build and run the program in development mode
start: build/electron-built build/npm-installed
	cd electron && npm start --arch=ia32

# Clean up
clean:
	rm -rf bin/* build/electron-built build/internal-generated electron/src/* electron/out $(shell find . -name 0_components_vgen.go)

# Clean up everything, including electron/node modules
fullclean: clean
	rm -rf electron/node_modules/* build/npm-installed electron/package-lock.json electron/package.json build/packagejsongen.exe

test:
	go test $$(dirname $$(find . -name \*_test.go))

##
## ChutzParse main build steps
##

bin/main.wasm: $(shell find cmd/main -name \*.go) $(shell find cmd/main -name \*.js) build/internal-generated $(shell find pkg -type f)
	GOOS=js GOARCH=wasm go build -tags electron -o $@ ./cmd/main

bin/window.wasm: $(shell find cmd/window -type f) build/internal-generated $(shell find pkg -type f)
	go run -mod=vendor github.com/vugu/vugu/cmd/vugugen -s -r -skip-go-mod -skip-main cmd/window
	GOOS=js GOARCH=wasm go build -tags web -o $@ ./cmd/window

bin/overlay.wasm: $(shell find cmd/overlay -type f) build/internal-generated $(shell find pkg -type f)
	go run -mod=vendor github.com/vugu/vugu/cmd/vugugen -s -r -skip-go-mod -skip-main cmd/overlay
	GOOS=js GOARCH=wasm go build -tags web -o $@ ./cmd/overlay

##
## Supporting infrastructure
##

$(shell mkdir -p bin electron/src/bin build>/dev/null 2>&1 || true)

# Download and install electron and other node modules
build/npm-installed:
	# FIXME: workaround until electron/windows-installer includes 7z-ia32.exe.  See issues:
	#    https://github.com/electron/windows-installer/issues/378
	#    https://github.com/electron/windows-installer/issues/386
	cd electron && npm install --arch=ia32 || true
	cd electron && npm install --arch=ia32 --ignore-scripts
	curl -L https://github.com/electron/windows-installer/raw/b2380345e8fe1ad7716108b10b552d75e6fad0b7/vendor/7z-ia32.dll -o electron/node_modules/electron-winstaller/vendor/7z.dll
	curl -L https://github.com/electron/windows-installer/raw/b2380345e8fe1ad7716108b10b552d75e6fad0b7/vendor/7z-ia32.exe -o electron/node_modules/electron-winstaller/vendor/7z.exe
	touch $@

build/internal-generated: $(shell find internal -type f)
	go run -mod=vendor github.com/vugu/vugu/cmd/vugugen -s -r -skip-go-mod -skip-main internal
	touch $@

# Populate the electron directory with our code and assets
build/electron-built: bin/main.wasm bin/window.wasm bin/overlay.wasm cmd/main/main.js cmd/main/preload.js $(shell find web/static/data -type f) electron/package.json
	cp -r web/static/data/* electron/src
	cp cmd/main/main.js cmd/main/preload.js electron/src
	cp bin/main.wasm bin/window.wasm bin/overlay.wasm electron/src/bin
	touch $@

electron/package.json: build/packagejsongen.exe
	build/packagejsongen.exe > $@

build/packagejsongen.exe: $(shell find cmd/packagejsongen -name \*.go) internal/version.go internal/package.json.go
	go build -tags native -o $@ ./cmd/packagejsongen
