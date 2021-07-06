// +build wasm,electron

// iff: Identification, friend or foe.
package iff

import (
	"github.com/gontikr99/chutzparse/internal/eqlog"
	"regexp"
	"strings"
	"time"
	"unicode"
)

const foeDuration = 2 * time.Minute
const friendDuration = 15 * time.Minute

var foes = map[string]time.Time{}
var friends = map[string]time.Time{}
var pets = map[string]string{}

func init() {
	// Periodically erase friends/foes which we haven't seen for a while
	go func() {
		for {
			<-time.After(1 * time.Second)
			now := time.Now()
			for name, expiration := range foes {
				if now.After(expiration) {
					delete(foes, name)
				}
			}
			for name, expiration := range friends {
				if now.After(expiration) {
					delete(foes, name)
				}
			}
		}
	}()
}

func Update(entry *eqlog.LogEntry) {
	MakeFriend(entry.Character)
	if entry.Meaning != nil {
		entry.Meaning.Visit(iffAction{entry})
	}
}

func MakeFriend(name string) {
	delete(foes, name)
	friends[name] = time.Now().Add(friendDuration)
}

func MakeFoe(name string) {
	delete(friends, name)
	foes[name] = time.Now().Add(foeDuration)
}

func MakePet(pet string, owner string) {
	pets[pet] = owner
}

func IsFriend(name string) bool {
	_, present := friends[name]
	return present
}

func IsFoe(name string) bool {
	_, present := foes[name]
	return present
}

func GetOwner(name string) string {
	owner, _ := pets[name]
	return owner
}

func heuristicId(name string) {
	if strings.HasSuffix(name, "`s pet") {
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
