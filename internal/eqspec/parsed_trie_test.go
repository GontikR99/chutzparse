package eqspec

import "testing"

func TestBuiltTrie(t *testing.T) {
	items := BuiltTrie.Scan("Al`Kabor's Pestle")
	if len(items) == 0 || items[0] != "Al`Kabor's Pestle" {
		t.Errorf("Match not found")
	}
}
