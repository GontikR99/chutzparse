//go:build wasm && electron
// +build wasm,electron

package bids

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gontikr99/chutzparse/internal/eqspec"
	"github.com/gontikr99/chutzparse/internal/settings"
	"github.com/gontikr99/chutzparse/pkg/console"
	"github.com/gontikr99/chutzparse/pkg/electron"
	"github.com/gontikr99/chutzparse/pkg/electron/browserwindow"
	"net/url"
	"regexp"
	"strconv"
	"strings"
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
			time.Sleep(60 * time.Second)
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
	page, code, err := electron.HttpCall(loc.Scheme, "GET", loc.Host, port, path+"Points/", nil, nil)
	if err != nil {
		return nil, err
	}
	if code >= 400 {
		return nil, fmt.Errorf("Failed to fetch DKP information (%d).  Bad path?", code)
	}
	for _, row := range scrapeRows(page) {
		cols := scrapeCols(row)
		if len(cols) != 8 || !strings.Contains(string(cols[0]), "Character/") {
			continue
		}
		name := htmlTrim(cols[0])
		balance, _ := strconv.Atoi(htmlTrim(cols[2]))
		result[name] = CharacterStat{
			Rank:       "",
			Balance:    int32(balance),
			Attendance: []string{htmlTrim(cols[3]), htmlTrim(cols[4]), htmlTrim(cols[5])},
		}
	}
	return result, nil
}

var rowStartRE = regexp.MustCompile("<\\s*tr")
var rowEndRE = regexp.MustCompile("<\\s*/\\s*tr")

func scrapeRows(rawPage []byte) [][]byte {
	rows := [][]byte{}

	for {
		idxs := rowStartRE.FindIndex(rawPage)
		if idxs == nil {
			return rows
		}
		rawPage = rawPage[idxs[0]:]
		closeIdx := strings.Index(string(rawPage), ">")
		if closeIdx < 0 {
			return rows
		}
		rawPage = rawPage[closeIdx+1:]
		idxs = rowEndRE.FindIndex(rawPage)
		if idxs == nil {
			return rows
		}
		row := rawPage[:idxs[0]]
		rows = append(rows, row)
		rawPage = rawPage[idxs[0]:]
		closeIdx = strings.Index(string(rawPage), ">")
		if closeIdx < 0 {
			return rows
		}
		rawPage = rawPage[closeIdx+1:]
	}
}

var colStartRE = regexp.MustCompile("<\\s*td")
var colEndRE = regexp.MustCompile("<\\s*/\\s*td")

func scrapeCols(rawRow []byte) [][]byte {
	cols := [][]byte{}

	for {
		idxs := colStartRE.FindIndex(rawRow)
		if idxs == nil {
			return cols
		}
		rawRow = rawRow[idxs[0]:]
		closeIdx := strings.Index(string(rawRow), ">")
		if closeIdx < 0 {
			return cols
		}
		rawRow = rawRow[closeIdx+1:]
		idxs = colEndRE.FindIndex(rawRow)
		if idxs == nil {
			return cols
		}
		col := rawRow[:idxs[0]]
		cols = append(cols, col)
		rawRow = rawRow[idxs[0]:]
		closeIdx = strings.Index(string(rawRow), ">")
		if closeIdx < 0 {
			return cols
		}
		rawRow = rawRow[closeIdx+1:]
	}
}

func htmlTrim(cell []byte) string {
	var result []byte
	idx := 0
	for idx < len(cell) {
		if cell[idx] != '<' {
			result = append(result, cell[idx])
			idx++
			continue
		}
		closeIdx := bytes.IndexByte(cell[idx:], '>')
		if closeIdx < 0 {
			break
		}
		idx = idx + closeIdx + 1
	}
	return strings.Trim(string(result), " \t\r\n")
}
