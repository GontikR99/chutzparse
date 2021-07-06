(function () {
    const {app} = require('electron');

    app.commandLine.appendSwitch("high-dpi-support", 1)
    // app.commandLine.appendSwitch("force-device-scale-factor", 1)

    const path = require('path');
    const fs = require('fs');
    require(path.join(app.getAppPath(), 'src/external/wasm_exec.js'));

    // Handle creating/removing shortcuts on Windows when installing/uninstalling.
    if (require('electron-squirrel-startup')) { // eslint-disable-line global-require
        app.quit();
        return;
    }

    const EventBarrier = function (signal) {
        const self = this;
        this.onSignalList = []
        this.onSignalEvent = null
        this.onSignal = function (newFunc) {
            if (self.onSignalEvent != null) {
                newFunc(self.onSignalEvent)
            } else {
                self.onSignalList.push(newFunc)
            }
        }
        app.on(signal, function (evt) {
            self.onSignalEvent = evt
            self.onSignalList.forEach(function (item) {
                item(evt)
            })
        })
    }

    global.eventBarriers = {
        "ready": new EventBarrier("ready"),
    }

    async function run() {
        const go = new Go();
        const mod = await WebAssembly.compile(fs.readFileSync(path.join(app.getAppPath(), 'src/bin/main.wasm')));
        let inst = await WebAssembly.instantiate(mod, go.importObject);
        go.run(inst);
    }

    run();
})()