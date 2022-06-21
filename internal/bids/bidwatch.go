//go:build wasm && electron
// +build wasm,electron

package bids

import (
	"github.com/gontikr99/chutzparse/internal/eqspec"
	"github.com/gontikr99/chutzparse/internal/settings"
	"github.com/gontikr99/chutzparse/pkg/console"
	"github.com/gontikr99/chutzparse/pkg/electron/browserwindow"
	"strconv"
	"strings"
)

var auctionActive = false
var activeBids = map[string]map[string]ItemBid{}

func init() {
	settings.DefaultSetting(settings.BidStartCmd, "!start")
	settings.DefaultSetting(settings.BidEndCmd, "!end")
	eqspec.RegisterLogsListener(func(entries []*eqspec.LogEntry) {
		didWork := false
		for _, entry := range entries {
			bidStart, _, err := settings.LookupSetting(settings.BidStartCmd)
			if err != nil || bidStart == "" {
				console.Log(err)
				return
			}
			bidEnd, _, err := settings.LookupSetting(settings.BidEndCmd)
			if err != nil || bidEnd == "" {
				console.Log(err)
				return
			}
			var msg *eqspec.ChatLog
			var ok bool
			if msg, ok = entry.Meaning.(*eqspec.ChatLog); !ok {
				continue
			}
			if msg.Method != eqspec.MethodTell {
				continue
			}
			if strings.EqualFold(entry.Character, msg.Source) {
				console.Log("Starting auction")
				if strings.EqualFold(bidStart, msg.Text) {
					auctionActive = true
					activeBids = map[string]map[string]ItemBid{}
					didWork = true
					continue
				} else if strings.EqualFold(bidEnd, msg.Text) {
					console.Log("Ending auction")
					auctionActive = false
					didWork = true
					continue
				}
			}
			if !auctionActive {
				continue
			}
			didWork = true
			var itemNames []string = eqspec.BuiltTrie.Scan(msg.Text)
			var item string
			if len(itemNames) == 0 || len(itemNames[0]) < 3 {
				item = UnspecifiedItem
			} else {
				item = itemNames[0]
			}
			var remText string
			nameIdx := strings.Index(msg.Text, item)
			if nameIdx < 0 {
				remText = msg.Text
			} else {
				remText = msg.Text[:nameIdx] + "|" + msg.Text[nameIdx+len(item):]
			}
			bidValue := extractBid(remText)
			if bidValue >= 0 {
				if _, ok := activeBids[item]; !ok {
					activeBids[item] = map[string]ItemBid{}
				}
				if _, ok := activeBids[item][msg.Source]; !ok {
					activeBids[item][msg.Source] = ItemBid{}
				}
				newBid := ItemBid{
					CalculatedBid: bidValue,
					BidMessages:   append(activeBids[item][msg.Source].BidMessages, msg.Text),
				}
				activeBids[item][msg.Source] = newBid
			}
		}
		if didWork {
			browserwindow.Broadcast(ChannelChange, []byte{})
		}
	})
}

func extractBid(msgText string) int32 {
	bestBid := int32(-1)
	var bidBuf []byte
	for _, c := range []byte(msgText) {
		if '0' <= c && c <= '9' {
			bidBuf = append(bidBuf, c)
		} else {
			bid, err := strconv.Atoi(string(bidBuf))
			if err == nil && int32(bid) > bestBid {
				bestBid = int32(bid)
			}
			bidBuf = []byte{}
		}
	}
	bid, err := strconv.Atoi(string(bidBuf))
	if err == nil && int32(bid) > bestBid {
		bestBid = int32(bid)
	}
	return bestBid
}
