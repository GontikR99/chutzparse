//go:build wasm && electron
// +build wasm,electron

package bids

import (
	"github.com/gontikr99/chutzparse/internal/eqspec"
	"github.com/gontikr99/chutzparse/internal/settings"
	"github.com/gontikr99/chutzparse/pkg/console"
	"github.com/gontikr99/chutzparse/pkg/electron/browserwindow"
	"github.com/gontikr99/chutzparse/pkg/nodejs/path"
	"io/ioutil"
	"os"
	"strings"
)

var guildMembers = map[string]MemberInfo{}

func init() {
	eqspec.RegisterLogsListener(func(entries []*eqspec.LogEntry) {
		for _, entry := range entries {
			oc, ok := entry.Meaning.(*eqspec.OutputfileComplete)
			if !ok {
				continue
			}
			eqdir, _, _ := settings.LookupSetting(settings.EverQuestDirectory)
			fullPath := path.Join(eqdir, oc.Filename)
			newMembers := map[string]MemberInfo{}
			contents, err := ioutil.ReadFile(fullPath)
			if err != nil {
				console.Log(err)
			}
			for _, lineRaw := range strings.Split(string(contents), "\n") {
				line := strings.Trim(lineRaw, "\r\n\t ")
				if line == "" {
					continue
				}
				parts := strings.Split(line, "\t")
				if len(parts) != 14 {
					newMembers = map[string]MemberInfo{}
					break
				}
				newMembers[strings.Trim(parts[0], " ")] = MemberInfo{
					Rank:    strings.Trim(parts[3], " "),
					Comment: strings.Trim(parts[7], " "),
				}
			}
			if len(newMembers) != 0 {
				guildMembers = newMembers
				os.Remove(fullPath)
				browserwindow.Broadcast(ChannelChange, []byte{})
			}
		}
	})
}
