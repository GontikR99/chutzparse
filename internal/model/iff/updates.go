package iff

import (
	"encoding/gob"
	"time"
)

func MakeFriend(name string) {
	update := &IffFriend{name}
	update.Apply()
	postUpdate(update)
}

func MakeFoe(name string) {
	update := &IffFoe{name}
	update.Apply()
	postUpdate(update)
}

func MakePet(pet string, owner string) {
	update := &IffPet{pet, owner}
	update.Apply()
	postUpdate(update)
}

const channelIffUpdate = "iffUpdate"

type IffUpdate interface {
	Apply()
}

type IffFriend struct{ Name string }
type IffFoe struct{ Name string }
type IffPet struct {
	Pet   string
	Owner string
}

func (i *IffFriend) Apply() {
	delete(foes, i.Name)
	friends[i.Name] = time.Now().Add(friendDuration)
}
func (i *IffFoe) Apply() {
	delete(friends, i.Name)
	foes[i.Name] = time.Now().Add(foeDuration)
}
func (i *IffPet) Apply() {
	pets[i.Pet] = i.Owner
}

type IffUpdateHolder struct {
	Update IffUpdate
}

func init() {
	gob.RegisterName("iff:friend", &IffFriend{})
	gob.RegisterName("iff:foe", &IffFoe{})
	gob.RegisterName("iff:pet", &IffPet{})
}
