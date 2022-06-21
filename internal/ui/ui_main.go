//go:build wasm && electron
// +build wasm,electron

package ui

import (
	"errors"
	"github.com/gontikr99/chutzparse/pkg/electron/dialog"
	"github.com/gontikr99/chutzparse/pkg/nodejs/clipboardy"
	"net/rpc"
)

type clipboardServer struct{}

func (c clipboardServer) Copy(text string) error {
	clipboardy.WriteSync(text)
	return nil
}

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

func HandleRPC() func(server *rpc.Server) {
	uiReg := handleClipboard(clipboardServer{})
	dirDlgReg := handleDirectoryDialog(dirDlgServer{})
	return func(server *rpc.Server) { uiReg(server); dirDlgReg(server) }
}
