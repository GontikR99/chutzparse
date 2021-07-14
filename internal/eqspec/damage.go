// +build wasm,electron

package eqspec

import (
	"github.com/gontikr99/chutzparse/pkg/multipattern"
	"regexp"
	"strconv"
)

type DamageType int

const (
	PhysicalDamage = DamageType(iota)
	MagicDamage
	FireDamage
	ColdDamage
	PoisonDamage
	DiseaseDamage
	CorruptionDamage
	ChromaticDamage
	PrismaticDamage
	UnresistableDamage
	UnspecifiedDamage
)

var damageTypeNames = []string{"physical", "magic", "fire", "cold", "poison", "disease", "corruption", "chromatic", "prismatic", "unresistable", "other"}

func (dt DamageType) String() string { return damageTypeNames[dt] }

type HitType int

const (
	BackstabHit = HitType(iota)
	BashHit
	BiteHit
	BluntHit
	BowHit
	ClawHit
	FrenzyHit
	HandHit
	HitHit
	KickHit
	PierceHit
	SlashHit
	OtherMeleeHit

	DamageShieldHit
	DotHit
	EnvironmentalHit
	SpellHit
)

var hitTypeNames = []string{"backstab", "bash", "bite", "crush", "shoot", "claw", "frenzy", "hand to hand", "hit",
	"kick", "pierce", "slash", "other melee",
	"damage shield", "unspecified DoT", "environmental", "unspecified spell"}

func (ht HitType) String() string {
	return hitTypeNames[ht]
}

const UnspecifiedName = "Unspecified"

type DamageLog struct {
	Source string
	Target string
	Amount int64

	Type      HitType
	Element   DamageType
	Flag      HitFlag
	SpellName string
}

func (d *DamageLog) Visit(handler ParsedLogHandler) interface{} { return handler.OnDamage(d) }

func (d *DamageLog) String() string {
	return "Damage: " + d.Source + ": " + strconv.FormatInt(d.Amount, 10) + " (" + d.Type.String() + "/" + d.Element.String() + ") -> " +
		d.Target + " [" + d.Flag.String() + "] " + d.SpellName
}

func (d *DamageLog) DisplayCategory() string {
	if d.SpellName != "" {
		return d.SpellName
	} else {
		return d.Type.String()
	}
}

var splitDotRE = regexp.MustCompile("^(.*) by (.*)$")

func handleDamage(mp *multipattern.Multipattern) *multipattern.Multipattern {
	return commonSubpatterns(mp).
		Define("mhit", "bites|bite|backstabs|backstab|claws|claw|crushes|crush|slashes|slash|slams|slam|"+
			"strikes|strike|shoots|shoot|pierces|pierce|punch|punches|frenzies on|frenzy on|kicks|kick|bashes|bash|hits|hit|"+
			"gores|gore|slices|slice|smashes|smash|stab|stabs|maul|mauls|sting|stings|rend|rends").
		// Melee:
		// Earthmaster Grundag hits YOU for 440 points of damage.
		// Keker backstabs a fire elemental raider for 249 points of damage.
		On("(.*) (@mhit@) (.*) for (@num@) points? of damage.(@hflag@)?", func(parts []string, _ interface{}) interface{} {
			var hitType HitType
			switch parts[2] {
			case "backstab":
				hitType = BackstabHit
			case "backstabs":
				hitType = BackstabHit
			case "bash":
				hitType = BashHit
			case "bashes":
				hitType = BashHit
			case "bite":
				hitType = BiteHit
			case "bites":
				hitType = BiteHit
			case "claw":
				hitType = ClawHit
			case "claws":
				hitType = ClawHit
			case "crush":
				hitType = BluntHit
			case "crushes":
				hitType = BluntHit
			case "frenzies on":
				hitType = FrenzyHit
			case "frenzy on":
				hitType = FrenzyHit
			case "hit":
				hitType = HitHit
			case "hits":
				hitType = HitHit
			case "kick":
				hitType = KickHit
			case "kicks":
				hitType = KickHit
			case "pierce":
				hitType = PierceHit
			case "pierces":
				hitType = PierceHit
			case "punch":
				hitType = HandHit
			case "punches":
				hitType = HandHit
			case "shoot":
				hitType = BowHit
			case "shoots":
				hitType = BowHit
			case "slash":
				hitType = SlashHit
			case "slashes":
				hitType = SlashHit
			case "slam":
				hitType = BashHit
			case "slams":
				hitType = BashHit
			case "strike":
				hitType = HandHit
			case "strikes":
				hitType = HandHit
			default:
				hitType = OtherMeleeHit
			}
			return &DamageLog{
				Source:  normalizeName(parts[1]),
				Target:  normalizeName(parts[3]),
				Amount:  amount(parts[4]),
				Type:    hitType,
				Element: PhysicalDamage,
				Flag:    hitFlags(parts[5]),
			}
		}).
		// Spell/proc:
		// a fire elemental blazemaster hit you for 304 points of cold damage by Draught of Ice.
		// You hit a girplan geomancer for 350 points of magic damage by Touch of the Cursed III.
		On("(.*) hit (.*) for (@num@) points? of (physical|magic|cold|fire|poison|disease|corruption|chromatic|prismatic|unresistable) damage by (.+)[.!](@hflag@)?",
			func(parts []string, _ interface{}) interface{} {
				var dmgType DamageType
				switch parts[4] {
				case "physical":
					dmgType = PhysicalDamage
				case "magic":
					dmgType = MagicDamage
				case "cold":
					dmgType = ColdDamage
				case "fire":
					dmgType = FireDamage
				case "poison":
					dmgType = PoisonDamage
				case "disease":
					dmgType = DiseaseDamage
				case "corruption":
					dmgType = CorruptionDamage
				case "chromatic":
					dmgType = ChromaticDamage
				case "prismatic":
					dmgType = PrismaticDamage
				case "unresistable":
					dmgType = UnresistableDamage
				default:
					dmgType = UnspecifiedDamage
				}
				return &DamageLog{
					Source:    normalizeName(parts[1]),
					Target:    normalizeName(parts[2]),
					Amount:    amount(parts[3]),
					Type:      SpellHit,
					Element:   dmgType,
					Flag:      hitFlags(parts[6]),
					SpellName: parts[5],
				}
			}).
		// Damage shield:
		// Zordak Ragefire is tormented by YOUR frost for 57 points of non-melee damage.
		// Zordak Ragefire is burned by Bardarsed's flames for 62 points of non-melee damage.
		On("(.*) (?:is|are) (burned|tormented|pierced) by (?:([^']+)'s|YOUR) (flames|thorns|frost|light) for (@num@) points? of non-melee damage[.!]",
			func(parts []string, _ interface{}) interface{} {
				var dmgType DamageType
				switch parts[4] {
				case "flames":
					dmgType = FireDamage
				case "thorns":
					dmgType = PhysicalDamage
				case "frost":
					dmgType = ColdDamage
				default:
					dmgType = UnspecifiedDamage
				}
				src := parts[3]
				if src == "" {
					src = "you"
				}
				return &DamageLog{
					Source:  normalizeName(src),
					Target:  normalizeName(parts[1]),
					Amount:  amount(parts[5]),
					Type:    DamageShieldHit,
					Element: dmgType,
					Flag:    0,
				}
			}).
		// Damage Shield (self)
		// You were hit by non-melee for 3 damage.
		On("You were hit by non-melee for (@num@) damage[.!](@hflag@)?", func(parts []string, _ interface{}) interface{} {
			return &DamageLog{
				Source:  normalizeName(UnspecifiedName),
				Target:  normalizeName("you"),
				Amount:  amount(parts[1]),
				Type:    DamageShieldHit,
				Element: UnspecifiedDamage,
				Flag:    hitFlags(parts[2]),
			}
		}).
		// Damage Shield (other):
		// Va Xi Aten Ha Ra was chilled to the bone for 28 points of non-melee damage.
		On("(.*) was (.*) for (@num@) points? of non-melee damage[.!](@hflag@)?", func(parts []string, _ interface{}) interface{} {
			return &DamageLog{
				Source:  normalizeName(UnspecifiedName),
				Target:  normalizeName(parts[1]),
				Amount:  amount(parts[3]),
				Type:    DamageShieldHit,
				Element: UnspecifiedDamage,
				Flag:    hitFlags(parts[4]),
			}
		}).
		// DOT (other):
		// An earth elemental intruder has taken 469 damage from Locust Swarm by Tavaren.
		// An earth elemental intruder has taken 1264 damage from Pyre of Mori by Grumpo. (Critical)
		On("(.*) has taken (@num@) damage from (.*) by (.*)[.!](@hflag@)?", func(parts []string, _ interface{}) interface{} {
			return &DamageLog{
				Source:    normalizeName(parts[4]),
				Target:    normalizeName(parts[1]),
				Amount:    amount(parts[2]),
				Type:      DotHit,
				Element:   UnspecifiedDamage,
				Flag:      hitFlags(parts[5]),
				SpellName: parts[3],
			}
		}).
		// DOT (self):
		// A water mephit has taken 20 damage from your Dooming Darkness.
		On("(.*) has taken (@num@) damage from your (.*)[.!](@hflag@)?", func(parts []string, _ interface{}) interface{} {
			return &DamageLog{
				Source:    normalizeName("you"),
				Target:    normalizeName(parts[1]),
				Amount:    amount(parts[2]),
				Type:      DotHit,
				Element:   UnspecifiedDamage,
				Flag:      hitFlags(parts[4]),
				SpellName: parts[3],
			}
		}).
		// DOT:
		// You have taken 57 damage from Chaos Claws by an elite ukun boneretriever.
		// You have taken 57 damage from Chaos Claws.
		// You have taken 294 damage from Bond of Inruku. (Critical)
		// Enchanter has taken 102 damage by Deathly Chants.
		On("(.*) (?:have|has) taken (@num@) damage (?:by|from) (.*)[.!](@hflag@)?", func(parts []string, _ interface{}) interface{} {
			var source string
			var spell string
			if srcSplit := splitDotRE.FindStringSubmatch(parts[3]); srcSplit == nil {
				source = UnspecifiedName
				spell = parts[3]
			} else {
				source = srcSplit[2]
				spell = srcSplit[1]
			}
			return &DamageLog{
				Source:    normalizeName(source),
				Target:    normalizeName(parts[1]),
				Amount:    amount(parts[2]),
				Type:      DotHit,
				Element:   UnspecifiedDamage,
				Flag:      hitFlags(parts[4]),
				SpellName: spell,
			}
		}).
		// Feedback:
		// Jowwn's mind burns from Vishimtar the Fallen's feedback for 250 points of non-melee damage.
		// Your mind burns from Vishimtar the Fallen's feedback for 250 points of non-melee damage!
		On("(?:(.*)'s|Your) mind burns from (.*)'s feedback for (@num@) points of non-melee damage[.!](@hflag@)?", func(parts []string, _ interface{}) interface{} {
			target := parts[1]
			if target == "" {
				target = "you"
			}
			return &DamageLog{
				Source:  normalizeName(parts[2]),
				Target:  normalizeName(target),
				Amount:  amount(parts[3]),
				Type:    EnvironmentalHit,
				Element: UnspecifiedDamage,
				Flag:    hitFlags(parts[4]),
			}
		}).
		// Misc:
		// Genarn is poisoned by Neimon of Air's venom for 125 points of non-melee damage.
		// YOU are poisoned by Neimon of Air's venom for 111 points of non-melee damage!
		On("(.*) (?:is|are) [^ ]+ by (.+)'s .* for (@num@) points of non-melee damage[.!](@hflag@)?", func(parts []string, _ interface{}) interface{} {
			return &DamageLog{
				Source:  normalizeName(parts[2]),
				Target:  normalizeName(parts[1]),
				Amount:  amount(parts[3]),
				Type:    EnvironmentalHit,
				Element: UnspecifiedDamage,
				Flag:    hitFlags(parts[4]),
			}
		})
}
