// +build tools

package tools

import (
	_ "github.com/shurcooL/vfsgen"
	_ "github.com/shurcooL/vfsgen/cmd/vfsgendev"
	_ "github.com/timshannon/bolthold"
	_ "github.com/vugu/html"
	_ "github.com/vugu/vjson"
	_ "github.com/vugu/vugu"
	_ "github.com/vugu/vugu/cmd/vugugen"
	_ "github.com/vugu/vugu/domrender"
	_ "github.com/vugu/vugu/js"
	_ "github.com/vugu/vugu/vgform"
)
