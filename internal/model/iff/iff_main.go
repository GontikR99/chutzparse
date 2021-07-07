// +build wasm,electron

package iff

import (
	"bytes"
	"encoding/gob"
	"github.com/gontikr99/chutzparse/internal/eqlog"
	"github.com/gontikr99/chutzparse/pkg/console"
	"github.com/gontikr99/chutzparse/pkg/electron/browserwindow"
	"regexp"
	"strings"
	"unicode"
)

func Update(entry *eqlog.LogEntry) {
	MakeFriend(entry.Character)
	if entry.Meaning != nil {
		entry.Meaning.Visit(iffAction{entry})
	}
}

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
	if name == eqlog.UnspecifiedName {

	} else if strings.HasSuffix(name, "`s pet") {
		MakePet(name, name[:len(name)-6])
	} else if strings.HasSuffix(name, "`s warder") {
		MakePet(name, name[:len(name)-9])
	} else if strings.HasSuffix(name, "`s ward") {
		MakePet(name, name[:len(name)-7])
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
	Entry *eqlog.LogEntry
}

func (i iffAction) OnDamage(log *eqlog.DamageLog) interface{} {
	if log.Source == log.Target {
		return nil
	}
	if log.Source == eqlog.UnspecifiedName || log.Target == eqlog.UnspecifiedName {
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

func (i iffAction) OnHeal(log *eqlog.HealLog) interface{} {
	if log.Source == log.Target {
		return nil
	}
	if log.Source == eqlog.UnspecifiedName || log.Target == eqlog.UnspecifiedName {
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

func (i iffAction) OnChat(log *eqlog.ChatLog) interface{} {
	if match := leaderRE.FindStringSubmatch(log.Text); match != nil {
		MakePet(log.Source, match[1])
	}
	return nil
}

func (i iffAction) OnDeath(log *eqlog.DeathLog) interface{} { return nil }
func (i iffAction) OnZone(log *eqlog.ZoneLog) interface{}   { return nil }
