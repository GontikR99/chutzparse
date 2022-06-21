//go:build wasm && electron
// +build wasm,electron

package randoms

import (
	"github.com/gontikr99/chutzparse/internal/eqspec"
	"github.com/gontikr99/chutzparse/pkg/electron/browserwindow"
	"net/rpc"
	"sort"
)

type randomsServer struct{}

type byBounds []*RollGroup

func (b byBounds) Len() int      { return len(b) }
func (b byBounds) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b byBounds) Less(i, j int) bool {
	if b[i].Max != b[j].Max {
		return b[i].Max < b[j].Max
	}
	return b[i].Min < b[j].Min
}

type byValueRev []*CharacterRoll

func (b byValueRev) Len() int           { return len(b) }
func (b byValueRev) Less(i, j int) bool { return b[i].Value > b[j].Value }
func (b byValueRev) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

func (r randomsServer) Reset() error {
	currentRolls = []*eqspec.RandomLog{}
	browserwindow.Broadcast(ChannelChange, []byte{})
	return nil
}

func (r randomsServer) FetchRandoms() ([]*RollGroup, error) {
	resultMap := map[int32]map[int32][]*CharacterRoll{}
	for _, roll := range currentRolls {
		if _, ok := resultMap[roll.Lbound]; !ok {
			resultMap[roll.Lbound] = map[int32][]*CharacterRoll{}
		}
		umap := resultMap[roll.Lbound]
		if _, ok := umap[roll.Ubound]; !ok {
			umap[roll.Ubound] = []*CharacterRoll{}
		}
		umap[roll.Ubound] = append(umap[roll.Ubound], &CharacterRoll{
			Character: roll.Source,
			Value:     roll.Value,
		})
	}
	result := []*RollGroup{}
	for lb, umap := range resultMap {
		for ub, rolllist := range umap {
			sort.Sort(byValueRev(rolllist))
			result = append(result, &RollGroup{
				Min:   lb,
				Max:   ub,
				Rolls: rolllist,
			})
		}
	}
	sort.Sort(byBounds(result))
	return result, nil
}

func HandleRPC() func(server *rpc.Server) {
	return handleRandoms(randomsServer{})
}
