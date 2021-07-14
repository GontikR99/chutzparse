// +build wasm,electron

package iff

import (
	"bytes"
	"encoding/gob"
	"github.com/gontikr99/chutzparse/internal/eqspec"
	"github.com/gontikr99/chutzparse/internal/settings"
	"github.com/gontikr99/chutzparse/pkg/console"
	"github.com/gontikr99/chutzparse/pkg/electron/browserwindow"
	"net/rpc"
	"regexp"
	"strings"
	"unicode"
)

func Update(entry *eqspec.LogEntry) {
	MakeFriend(entry.Character)
	if entry.Meaning != nil {
		entry.Meaning.Visit(iffAction{entry})
	}
}

// postUpdate sends an IFF change made in the main process to the renderer(s)
func postUpdate(update IffUpdate) {
	buf := &bytes.Buffer{}
	err := gob.NewEncoder(buf).Encode(&IffUpdateHolder{Update: update})
	if err == nil {
		browserwindow.Broadcast(channelIffUpdate, buf.Bytes())
	} else {
		console.Log(err)
	}
}

func heuristicId(name string) {
	linkText, present, err := settings.LookupSetting(settings.LinkObviousPets)
	linkObviousPets := false
	if present && err == nil && linkText == "true" {
		linkObviousPets = true
	}
	if name == eqspec.UnspecifiedName {

	} else if strings.HasSuffix(name, "`s pet") {
		if linkObviousPets {
			MakePet(name, name[:len(name)-6])
		}
	} else if strings.HasSuffix(name, "`s warder") {
		if linkObviousPets {
			MakePet(name, name[:len(name)-9])
		}
	} else if strings.HasSuffix(name, "`s ward") {
		if linkObviousPets {
			MakePet(name, name[:len(name)-7])
		}
	} else {
		for _, c := range name {
			if !unicode.IsLetter(rune(c)) {
				MakeFoe(name)
				return
			}
		}
	}
}

type iffAction struct {
	Entry *eqspec.LogEntry
}

func (i iffAction) OnDamage(log *eqspec.DamageLog) interface{} {
	if log.Source == log.Target {
		return nil
	}
	if log.Source == eqspec.UnspecifiedName || log.Target == eqspec.UnspecifiedName {
		return nil
	}
	heuristicId(log.Source)
	heuristicId(log.Target)

	if IsFriend(log.Source) {
		MakeFoe(log.Target)
		return nil
	}
	if IsFriend(log.Target) {
		MakeFoe(log.Source)
		return nil
	}
	return nil
}

func (i iffAction) OnHeal(log *eqspec.HealLog) interface{} {
	if log.Source == log.Target {
		return nil
	}
	if log.Source == eqspec.UnspecifiedName || log.Target == eqspec.UnspecifiedName {
		return nil
	}
	heuristicId(log.Source)
	heuristicId(log.Target)

	if IsFriend(log.Source) {
		MakeFriend(log.Target)
		return nil
	}
	if IsFriend(log.Target) {
		MakeFriend(log.Source)
		return nil
	}

	return nil
}

var leaderRE = regexp.MustCompile("^My leader is ([A-Z][a-z]+)[.]$")

func (i iffAction) OnChat(log *eqspec.ChatLog) interface{} {
	linkText, present, err := settings.LookupSetting(settings.LinkObviousPets)
	linkObviousPets := false
	if present && err == nil && linkText == "true" {
		linkObviousPets = true
	}
	if match := leaderRE.FindStringSubmatch(log.Text); match != nil {
		if linkObviousPets {
			MakePet(log.Source, match[1])
		}
	}
	return nil
}

func (i iffAction) OnDeath(log *eqspec.DeathLog) interface{} { return nil }
func (i iffAction) OnZone(log *eqspec.ZoneLog) interface{}   { return nil }

func init() {
	settings.DefaultSetting(settings.LinkObviousPets, "true")
}

type iffStub struct{}

func (i iffStub) Unlink(pet string) error {
	UnlinkPet(pet)
	return nil
}

func (i iffStub) Link(pet string, owner string) error {
	MakePet(pet, owner)
	return nil
}

func HandleRPC() func(server *rpc.Server) {
	return handleControl(iffStub{})
}
