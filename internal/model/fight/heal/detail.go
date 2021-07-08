// +build wasm,web

package heal

import (
	"fmt"
	"github.com/vugu/vugu"
	"strconv"
	"time"
)

func (r *Report) Detail() vugu.Builder {
	result := &Detail{
		report:        r,
		openedSources: map[string]struct{}{},
	}
	return result
}

// Detail is the healing report's detailed view
type Detail struct {
	report        *Report
	openedSources map[string]struct{}
}

func (c *Detail) toggle(event vugu.DOMEvent, source string) {
	event.PreventDefault()
	if _, present := c.openedSources[source]; present {
		delete(c.openedSources, source)
	} else {
		c.openedSources[source] = struct{}{}
	}
}

type toggleStateT int

const (
	toggleAbsent = toggleStateT(iota)
	toggleClosed
	toggleOpen
)

type detailRow struct {
	ToggleState toggleStateT
	BgColor     string
	FontSize    string

	Rank             string
	AttributedSource string
	Source           string
	Category         string
	Percent          string
	Amount           string
	HPS              string
}

func (c *Detail) rows() []*detailRow {
	duration := c.report.ActivitySet.TotalDuration()
	if duration < time.Second {
		duration = time.Second
	}
	durationSec := float64(duration) / float64(time.Second)

	var rows []*detailRow
	agRep := c.report.Aggregate()
	for idx, contrib := range agRep.SortedContributors() {
		toggleState := toggleClosed
		if _, present := c.openedSources[contrib.AttributedSource]; present {
			toggleState = toggleOpen
		}
		bgColor := "rgba(255,255,255,0.05)"
		if idx&1 == 0 {
			bgColor = "rgb(128, 128, 128)"
		}
		rows = append(rows, &detailRow{
			ToggleState:      toggleState,
			BgColor:          bgColor,
			FontSize:         "medium",
			Rank:             strconv.FormatInt(int64(1+idx), 10),
			AttributedSource: contrib.AttributedSource,
			Source:           contrib.DisplayName(),
			Percent:          fmt.Sprintf("%.3g%%", float64(100*contrib.TotalHealed)/float64(agRep.TotalHealed)),
			Amount:           strconv.FormatInt(contrib.TotalHealed, 10),
			HPS:              fmt.Sprintf("%.0f", float64(contrib.TotalHealed)/durationSec),
		})
		if _, present := c.openedSources[contrib.AttributedSource]; present {
			for _, cat := range contrib.SortedCategories() {
				rows = append(rows, &detailRow{
					ToggleState: toggleAbsent,
					BgColor:     bgColor,
					FontSize:    "small",

					Category: cat.DisplayName,
					Percent:  fmt.Sprintf("%.3g%%", float64(100*cat.TotalHealed)/float64(contrib.TotalHealed)),
					Amount:   strconv.FormatInt(cat.TotalHealed, 10),
					HPS:      fmt.Sprintf("%.0f", float64(cat.TotalHealed)/durationSec),
				})
			}
		}
	}
	return rows
}
