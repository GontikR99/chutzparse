// +build native

package internal

type dict map[string]interface{}
type list []interface{}

var PackageJson = dict{
	"name": "chutzparse",
	"productName": "chutzparse",
	"version": Version,
	"description": "An EverQuest log parser and heads-up display",
	"main": "src/main.js",
	"scripts": dict{
		"start": "electron-forge start",
		"package": "electron-forge package",
		"make": "electron-forge make",
		"publish": "electron-forge publish",
	},
	"repository": dict{
		"type": "git",
		"url": "git+https://github.com/gontikr99/chutzparse.git",
	},
	"keywords": list{},
	"author": dict{
		"name": "GontikR99",
	},
	"license": "MIT",
	"config": dict{
		"forge": dict{
			"packagerConfig": dict{
				"icon": "./src/favicon.ico",
				"iconUrl": "https://raw.githubusercontent.com/GontikR99/chutzparse/master/web/static/data/favicon.ico",
			},
			"makers": list{
				dict{
					"name": "@electron-forge/maker-squirrel",
					"config": dict{
						"name": "chutzparse",
					},
				},
			},
		},
	},
	"iohook": dict{
		"targets": list{
			"electron-87",
			"node-88",
		},
		"platforms": list{
			"win32",
		},
		"arches": list{
			"ia32",
		},
	},
	"dependencies": dict{
		"electron-squirrel-startup": "^1.0.0",
		"iohook": "^0.9.3",
		"ref-napi": "^3.0.3",
		"win32-api": "^9.6.0",
		"clipboardy": "^2.3.0",
	},
	"devDependencies": dict{
		"@electron-forge/cli": "^6.0.0-beta.57",
		"@electron-forge/maker-squirrel": "^6.0.0-beta.57",
		"electron": "12.0.12",
	},
}
