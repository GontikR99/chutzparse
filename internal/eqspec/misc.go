// +build wasm,electron

package eqspec

import (
	"github.com/gontikr99/chutzparse/pkg/multipattern"
	"regexp"
	"strconv"
	"strings"
)

func amount(text string) int64 {
	result, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		panic(err)
	}
	return result
}

var statusRE = regexp.MustCompile("^(.+) \\(.*\\)$")

func normalizeName(name string) string {
	if strings.EqualFold("you", name) || strings.EqualFold("your", name) || strings.EqualFold("yourself", name) {
		return "You"
	}
	if strings.HasSuffix(name, "'s corpse") {
		//return name[:len(name)-9]
		//FIXME: hack:  Attribute damage/healing from corpses to nobody
		return UnspecifiedName
	}
	// Handle names like "Mayong Mistmoore (Vulnerable)" or "Kessdona (Frozen Aura)"
	if m := statusRE.FindStringSubmatch(name); m != nil {
		name = m[1]
	}
	if name != "" {
		name = strings.ToUpper(name[:1]) + name[1:]
	}
	return name
}

type HitFlag int

const (
	CriticalFlag      = HitFlag(1 << 0)
	RiposteFlag       = HitFlag(1 << 1)
	RampageFlag       = HitFlag(1 << 2)
	WildRampageFlag   = HitFlag(1 << 3)
	FlurryFlag        = HitFlag(1 << 4)
	StrikethroughFlag = HitFlag(1 << 5)
	FinishingBlowFlag = HitFlag(1 << 6)
	DoubleBowFlag     = HitFlag(1 << 7)
	CrippleFlag       = HitFlag(1 << 8)
	SlayFlag          = HitFlag(1 << 9)
	DeadlyStrikeFlag  = HitFlag(1 << 10)
	AssassinateFlag   = HitFlag(1 << 11)
	HeadshotFlag      = HitFlag(1 << 12)
)

var hitFlagNames = []string{"Critical", "Riposte", "Rampage", "Wild Rampage", "Flurry", "Strikethrough", "Finishing Blow",
	"Double Bow Shot", "Crippling Blow", "Slay Undead", "Deadly Strike", "Assassinate", "Headshot"}

func (hf HitFlag) String() string {
	sb := &strings.Builder{}
	needSep := false
	for i, name := range hitFlagNames {
		if (hf & (1 << i)) != 0 {
			if needSep {
				sb.WriteByte('|')
			} else {
				needSep = true
			}
			sb.WriteString(name)
		}
	}
	return sb.String()
}

func hitFlags(text string) HitFlag {
	result := HitFlag(0)
	if strings.Contains(text, "Critical") {
		result |= CriticalFlag
	}
	if strings.Contains(text, "Riposte") {
		result |= RiposteFlag
	}
	if strings.Contains(text, "Flurry") {
		result |= FlurryFlag
	}
	if strings.Contains(text, "Strikethrough") {
		result |= StrikethroughFlag
	}
	if strings.Contains(text, "Wild Rampage") {
		result |= WildRampageFlag
	} else if strings.Contains(text, "Rampage") {
		result |= RampageFlag
	}
	if strings.Contains(text, "Finishing Blow") {
		result |= FinishingBlowFlag
	}
	if strings.Contains(text, "Double Bow Shot") {
		result |= DoubleBowFlag
	}
	if strings.Contains(text, "Crippling Blow") {
		result |= CrippleFlag
	}
	if strings.Contains(text, "Slay Undead") {
		result |= SlayFlag
	}
	if strings.Contains(text, "Deadly Strike") {
		result |= DeadlyStrikeFlag
	}
	if strings.Contains(text, "Assassinate") {
		result |= AssassinateFlag
	}
	if strings.Contains(text, "Headshot") {
		result |= HeadshotFlag
	}
	return result
}

func commonSubpatterns(mp *multipattern.Multipattern) *multipattern.Multipattern {
	return mp.Define("num", "0|[1-9][0-9]*").
		Define("hflag", "\\s\\((?:(?:"+HitFlag(-1).String()+")\\s?)+\\)")
}
