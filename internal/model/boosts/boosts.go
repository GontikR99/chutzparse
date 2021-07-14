// +build wasm

package boosts

type BoostType int

const (
	BardEpic1_5 = iota
	BardEpic2
	ShamanEpic1_5
	ShamanEpic2
	IntensityResolute
	GlyphDestruction
)
const boostChars = "bBsS7$"

func (bt BoostType) String() string {
	return boostChars[bt : bt+1]
}

type BoostSet map[BoostType]map[int]struct{}

func (bs BoostSet) Add(bt BoostType, logId int) {
	update, ok := bs[bt]
	if !ok {
		update = map[int]struct{}{}
		bs[bt] = update
	}
	update[logId] = struct{}{}
}

func (bs BoostSet) AddAll(other BoostSet) {
	for bt, idMap := range other {
		for id := range idMap {
			bs.Add(bt, id)
		}
	}
}
