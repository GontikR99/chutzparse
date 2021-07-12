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
	update := &IffNewPet{pet, owner}
	update.Apply()
	postUpdate(update)
}

func UnlinkPet(pet string) {
	update := &IffUnlinkPet{pet}
	update.Apply()
	postUpdate(update)
}

const channelIffUpdate = "iffUpdate"

type IffUpdate interface {
	Apply()
}

type IffFriend struct{ Name string }
type IffFoe struct{ Name string }
type IffNewPet struct {
	Pet   string
	Owner string
}
type IffUnlinkPet struct {
	Pet string
}

func (i *IffFriend) Apply() {
	delete(foes, i.Name)
	friends[i.Name] = time.Now().Add(friendDuration)
}
func (i *IffFoe) Apply() {
	delete(friends, i.Name)
	foes[i.Name] = time.Now().Add(foeDuration)
}
func (i *IffNewPet) Apply() {
	if oldOwner, present := pets[i.Pet]; !present || oldOwner!=i.Owner {
		pets[i.Pet] = i.Owner
		petChange(i)
	}
}
func (i *IffUnlinkPet) Apply() {
	if _, present := pets[i.Pet]; present {
		delete(pets, i.Pet)
		petChange(i)
	}
}

type IffUpdateHolder struct {
	Update IffUpdate
}

func init() {
	gob.RegisterName("iff:friend", &IffFriend{})
	gob.RegisterName("iff:foe", &IffFoe{})
	gob.RegisterName("iff:pet", &IffNewPet{})
	gob.RegisterName("iff:unlinkpet", &IffUnlinkPet{})
}

var petListenerGen=0
var petListeners=map[int]chan<- struct{}{}

// ListenPets lets us listen for changes to the pet mapping
func ListenPets() (<-chan struct{}, func()) {
	id := petListenerGen
	petListenerGen++
	updateChan := make(chan struct{})
	petListeners[id]=updateChan
	doneFunc := func() {delete(petListeners, id)}
	return updateChan, doneFunc
}

func petChange(update IffUpdate) {
	for _, outChan := range petListeners {
		func(){
			defer func() {recover()}()
			outChan <- struct{}{}
		}()
	}
}