// +build wasm,electron

package mainrpc

import (
	"errors"
	"github.com/gontikr99/chutzparse/internal/rpc"
	"github.com/gontikr99/chutzparse/pkg/electron/dialog"
)

type dirDlgServer struct{}

func (d dirDlgServer) Choose(initial string) (chosenDirectory string, err error) {
	filePaths, err := dialog.ShowOpenDialog(&dialog.OpenOptions{
		Title:       "Select a directory",
		DefaultPath: initial,
		Properties:  &[]string{dialog.OpenDirectory, dialog.DontAddToRecent},
	})

	if err != nil {
		return "", err
	}
	if len(filePaths) != 1 {
		return "", errors.New("expected single path")
	}
	return filePaths[0], nil
}

func init() {
	register(rpc.HandleDirectoryDialog(dirDlgServer{}))
}
