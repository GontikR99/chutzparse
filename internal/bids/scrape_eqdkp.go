//go:build wasm && electron
// +build wasm,electron

package bids

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/gontikr99/chutzparse/internal/eqspec"
	"github.com/gontikr99/chutzparse/internal/settings"
	"github.com/gontikr99/chutzparse/pkg/console"
	"github.com/gontikr99/chutzparse/pkg/electron"
	"github.com/gontikr99/chutzparse/pkg/electron/browserwindow"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

var currentDKP = map[string]CharacterStat{}
var hasDKP = eqspec.NewItemTrie().Compress()

func init() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				console.Log(r)
			}
		}()
		for {
			_, err := refreshDKP()
			if err != nil {
				console.Log(err)
			}
			time.Sleep(5 * 60 * time.Second)
		}
	}()
}

func refreshDKP() (int32, error) {
	site, _, err := settings.LookupSetting(settings.EqDkpSite)
	if err != nil {
		return 0, err
	}
	if site == "" {
		return 0, errors.New("No site specified")
	}
	cd, err := scrapeEQDKP(site)
	if err != nil {
		return 0, err
	}
	newTrie := eqspec.NewItemTrie()
	for k, _ := range cd {
		newTrie = newTrie.With(strings.ToUpper(k))
	}
	currentDKP = cd
	hasDKP = newTrie.Compress()
	browserwindow.Broadcast(ChannelChange, []byte{})
	return int32(len(currentDKP)), nil
}

type EQDKPRaids struct {
	XMLName xml.Name     `xml:"response"`
	Raids   []*EQDKPRaid `xml:"raid"`
}

type EQDKPRaid struct {
	Timestamp int64       `xml:"date_timestamp"`
	Note      string      `xml:"note"`
	Value     int32       `xml:"value"`
	EventId   int32       `xml:"event_id"`
	PlayerIds PlayerIdSet `xml:"raid_attendees"`
}

type PlayerIdSet map[int32]struct{}

func (r *PlayerIdSet) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	result := map[int32]struct{}{}
	for {
		t, err := decoder.Token()
		if err != nil {
			return err
		}
		if ee, ok := t.(xml.EndElement); ok && ee.Name == start.Name {
			break
		}
		if se, ok := t.(xml.StartElement); ok && strings.HasPrefix(se.Name.Local, "i") {
			var playerId int32
			err = decoder.DecodeElement(&playerId, &se)
			if err != nil {
				return err
			}
			result[playerId] = struct{}{}
		}
	}
	*r = result
	time.Sleep(1 * time.Millisecond)
	return nil
}

type EQDKPPoints struct {
	XMLName xml.Name       `xml:"response"`
	Players []*EQDKPPlayer `xml:"players>player"`
	Now     int64          `xml:"info>timestamp"`
}

type EQDKPPlayer struct {
	PlayerId int32      `xml:"id"`
	Name     string     `xml:"name"`
	Active   int        `xml:"active"`
	Hidden   int        `xml:"hidden"`
	Points   int32      `xml:"points>multidkp_points>points_current"`
	Pause    EQDKPPause `xml:"items"`
}

// Arrange schedule yielding points periodically so interface doesn't freeze parsing players
type EQDKPPause struct{}

func (E *EQDKPPause) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}
		if ee, ok := t.(xml.EndElement); ok && ee.Name == start.Name {
			break
		}
	}
	time.Sleep(1 * time.Millisecond)
	return nil
}

type EQDKPEvents struct {
	XMLName xml.Name      `xml:"response"`
	Events  []*EQDKPEvent `xml:"event"`
}

type EQDKPEvent struct {
	Id         int32 `xml:"id"`
	Attendance int32 `xml:"multidkp_pools>multidkp_pool>attendance""`
}

type attendanceTally struct {
	raids30 int32
	raids60 int32
	raids90 int32
}

func scrapeEQDKP(site string) (map[string]CharacterStat, error) {
	result := map[string]CharacterStat{}

	if !strings.Contains(site, "://") {
		site = "https://" + site
	}
	loc, err := url.Parse(site)
	if err != nil {
		return nil, err
	}
	path := loc.Path
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	var port int16
	if strings.EqualFold("https", loc.Scheme) {
		port = 443
	} else if strings.EqualFold("http", loc.Scheme) {
		port = 80
	}
	portStr := loc.Port()
	if portStr != "" {
		pnum, _ := strconv.Atoi(loc.Port())
		port = int16(pnum)
	}

	var points EQDKPPoints
	var raids EQDKPRaids
	var events EQDKPEvents
	var fetcherr error

	sg := sync.WaitGroup{}
	sg.Add(1)
	go func() {
		defer sg.Done()
		pointsbytes, code, err := electron.HttpCall(loc.Scheme, "GET", loc.Host, port, path+"api.php?function=points", nil, nil)
		if err != nil {
			fetcherr = err
			return
		}
		if code >= 400 {
			fetcherr = fmt.Errorf("Failed to fetch DKP information (%d).  Bad path?", code)
			return
		}
		err = xml.Unmarshal(pointsbytes, &points)
		if err != nil {
			fetcherr = err
			return
		}
	}()
	sg.Add(1)
	go func() {
		defer sg.Done()
		raidsbytes, code, err := electron.HttpCall(loc.Scheme, "GET", loc.Host, port, path+"api.php?function=raids", nil, nil)
		if err != nil {
			fetcherr = err
			return
		}
		if code >= 400 {
			fetcherr = fmt.Errorf("Failed to fetch DKP information (%d).  Bad path?", code)
			return
		}
		err = xml.Unmarshal(raidsbytes, &raids)
		if err != nil {
			fetcherr = err
			return
		}
	}()
	sg.Add(1)
	go func() {
		defer sg.Done()
		eventsbytes, code, err := electron.HttpCall(loc.Scheme, "GET", loc.Host, port, path+"api.php?function=events", nil, nil)
		if err != nil {
			fetcherr = err
			return
		}
		if code >= 400 {
			fetcherr = fmt.Errorf("Failed to fetch DKP information (%d).  Bad path?", code)
			return
		}
		err = xml.Unmarshal(eventsbytes, &events)
		if err != nil {
			fetcherr = err
			return
		}
	}()
	sg.Wait()
	if fetcherr != nil {
		return nil, fetcherr
	}

	// Calculate attendance
	attendable := map[int32]struct{}{}
	for _, event := range events.Events {
		if event.Attendance != 0 {
			attendable[event.Id] = struct{}{}
		}
	}

	tallies := map[int32]*attendanceTally{}
	var raids30, raids60, raids90 int32
	var lastUpdate int64
	for _, raid := range raids.Raids {
		if _, ok := attendable[raid.EventId]; !ok {
			continue
		}
		if raid.Timestamp > lastUpdate {
			lastUpdate = raid.Timestamp
		}
	}
	for _, raid := range raids.Raids {
		if _, ok := attendable[raid.EventId]; !ok {
			continue
		}
		for pid, _ := range raid.PlayerIds {
			if _, ok := tallies[pid]; !ok {
				tallies[pid] = &attendanceTally{}
			}
		}
		if lastUpdate-raid.Timestamp < int64(30*24*60*60) {
			raids30 += 1
			for pid, _ := range raid.PlayerIds {
				tallies[pid].raids30 += 1
			}
		}
		if lastUpdate-raid.Timestamp < int64(60*24*60*60) {
			raids60 += 1
			for pid, _ := range raid.PlayerIds {
				tallies[pid].raids60 += 1
			}
		}
		if lastUpdate-raid.Timestamp < int64(90*24*60*60) {
			raids90 += 1
			for pid, _ := range raid.PlayerIds {
				tallies[pid].raids90 += 1
			}
		}
		time.Sleep(1 * time.Millisecond)
	}

	// Create characterstat map
	for _, player := range points.Players {
		if player.Active == 0 {
			continue
		}
		if player.Hidden != 0 {
			continue
		}
		cs := CharacterStat{
			Rank:       "",
			Balance:    player.Points,
			Attendance: nil,
		}
		if tally, ok := tallies[player.PlayerId]; ok {
			if raids30 != 0 {
				cs.Attendance = append(cs.Attendance, fmt.Sprintf("%d%% (%d/%d), 30 day", 100*tally.raids30/raids30, tally.raids30, raids30))
			}
			if raids60 != 0 {
				cs.Attendance = append(cs.Attendance, fmt.Sprintf("%d%% (%d/%d), 60 day", 100*tally.raids60/raids60, tally.raids60, raids60))
			}
			if raids90 != 0 {
				cs.Attendance = append(cs.Attendance, fmt.Sprintf("%d%% (%d/%d), 90 day", 100*tally.raids90/raids90, tally.raids90, raids90))
			}
		}
		result[player.Name] = cs
	}
	return result, nil
}
