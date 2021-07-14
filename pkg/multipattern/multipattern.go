package multipattern

import (
	"regexp"
	"strings"
)

type Multipattern struct {
	needSep       bool
	rexText       *strings.Builder
	rex           *regexp.Regexp
	offsets       []int
	callbacks     []func(parts []string, context interface{}) interface{}
	substitutions map[string]string
}

func New() *Multipattern {
	result := &Multipattern{
		needSep:       false,
		rexText:       &strings.Builder{},
		offsets:       []int{1},
		substitutions: map[string]string{},
	}
	result.rexText.WriteString("^(?:")
	return result
}

func (mp *Multipattern) Define(name string, subexpr string) *Multipattern {
	name = "@" + name + "@"
	for needle, value := range mp.substitutions {
		subexpr = strings.ReplaceAll(subexpr, needle, value)
	}
	subexpr = "(?:" + subexpr + ")"
	if oldVal, present := mp.substitutions[name]; present && oldVal != subexpr {
		panic("Conflicting definitions for " + name + ": " + oldVal + " / " + subexpr)
	}
	mp.substitutions[name] = subexpr
	return mp
}

func (mp *Multipattern) On(pattern string, callback func(parts []string, context interface{}) interface{}) *Multipattern {
	mp.rex = nil

	for needle, value := range mp.substitutions {
		pattern = strings.ReplaceAll(pattern, needle, value)
	}

	rex := regexp.MustCompile(pattern)
	if mp.needSep {
		mp.rexText.WriteString("|")
	} else {
		mp.needSep = true
	}
	mp.rexText.WriteString("(x")
	mp.rexText.WriteString(pattern)
	mp.rexText.WriteString(")")

	mp.callbacks = append(mp.callbacks, callback)
	mp.offsets = append(mp.offsets, 1+rex.NumSubexp()+mp.offsets[len(mp.offsets)-1])
	return mp
}

func (mp *Multipattern) Dispatch(needle string, context interface{}) interface{} {
	if mp.rex == nil {
		mp.rex = regexp.MustCompile(mp.rexText.String() + ")$")
	}
	submatch := mp.rex.FindStringSubmatch("x" + needle)
	if submatch != nil {
		for i := 0; i < len(mp.callbacks); i++ {
			if submatch[mp.offsets[i]] != "" {
				parts := []string{submatch[mp.offsets[i]][1:]}
				parts = append(parts, submatch[1+mp.offsets[i]:mp.offsets[i+1]]...)
				return mp.callbacks[i](parts, context)
			}
		}
	}
	return nil
}
