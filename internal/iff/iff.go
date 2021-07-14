// iff: Identification, friend or foe.
package iff

import (
	"time"
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

func GetPets() map[string]string {
	result := map[string]string{}
	for k, v := range pets {
		result[k] = v
	}
	return result
}
