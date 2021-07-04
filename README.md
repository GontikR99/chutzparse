# ChutzParse, an EverQuest log parser and Heads Up Display

## Building
While my usual development environment is Linux, I developed this software in Windows.  To start, you
need to install:

* MSYS 2 (Recommended install: `C:\msys64`)
* Go 1.16+ (make sure the install path doesn't have any spaces in it.  Recommended: `C:\msys64\go`)
* Node v14.17.2 **32-bit** (make sure the install path doesn't have any spaces in it.  Recommended:
`C:\msys64\nodejs`)
* Python 2.7.17 (Default installation location: `C:\Python27`)

In MSYS 2, edit your `.bashrc` to add node and go to your path.  In the recommended scenario this would
be something like:

```bash
export PATH=$PATH:/go/bin:/nodejs
```

Once that's done, one can run one of the following commands in the project directory:

* `make start`: Build ChutzParse and run in developer mode
* `make package`: Build a ChutzParse package, `chutzparse-x.x.x Setup.exe` in `bin/`
* `make clean`: Remove ChutzParse intermediate binaries
* `make full-clean`: Remove ChutzParse intermediate binaries and all downloaded/built node/electron modules.
