//go:build wasm && electron
// +build wasm,electron

package bids

import (
	"net/rpc"
	"sort"
)

type bidsServer struct{}

func (b bidsServer) RefreshDKP() (int32, error) {
	chars, err := refreshDKP()
	return chars, err
}

func (b bidsServer) AuctionActive() (bool, error) {
	return auctionActive, nil
}

type byCostRev []*AnnotatedBid

func (b byCostRev) Len() int { return len(b) }
func (b byCostRev) Less(i, j int) bool {
	if b[i].Bid.CalculatedBid == b[j].Bid.CalculatedBid {
		return b[i].Character < b[j].Character
	}
	return b[i].Bid.CalculatedBid > b[j].Bid.CalculatedBid
}
func (b byCostRev) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

type byName []*ItemBids

func (b byName) Len() int { return len(b) }
func (b byName) Less(i, j int) bool {
	if b[i].Item == UnspecifiedItem && b[j].Item == UnspecifiedItem {
		return false
	}
	if b[i].Item == UnspecifiedItem && b[j].Item != UnspecifiedItem {
		return false
	}
	if b[i].Item != UnspecifiedItem && b[j].Item == UnspecifiedItem {
		return true
	}
	return b[i].Item < b[j].Item
}
func (b byName) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

func (b bidsServer) FetchBids() ([]*ItemBids, error) {
	result := []*ItemBids{}
	for itemName, bidders := range activeBids {
		charBids := []*AnnotatedBid{}
		for bidder, bid := range bidders {
			annBid := &AnnotatedBid{
				Character: bidder,
				Bid:       bid,
				Stat: CharacterStat{
					Rank:       "???",
					Balance:    -1,
					Attendance: nil,
				},
			}
			if stat, ok := currentDKP[bidder]; ok {
				annBid.Stat = stat
			}
			charBids = append(charBids, annBid)
		}
		sort.Sort(byCostRev(charBids))
		result = append(result, &ItemBids{
			Item: itemName,
			Bids: charBids,
		})
	}
	sort.Sort(byName(result))
	return result, nil
}

func HandleRPC() func(server *rpc.Server) {
	return handleBids(bidsServer{})
}
