// +build wasm,electron

package eqlog

import (
	"bytes"
	"context"
	"github.com/gontikr99/chutzparse/internal/settings"
	"github.com/gontikr99/chutzparse/pkg/console"
	"github.com/gontikr99/chutzparse/pkg/multipattern"
	"github.com/gontikr99/chutzparse/pkg/nodejs/fs"
	"github.com/gontikr99/chutzparse/pkg/nodejs/path"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	filenameMatch = regexp.MustCompile("^eqlog_([A-Za-z]*)_([A-Za-z]*).txt$")
	loglineMatch  = regexp.MustCompile("^\\[([^\\]]*)] (.*)$")
)

type ListenerHandle int

var nextListenerHandler = ListenerHandle(1)
var logListeners = make(map[ListenerHandle]func([]*LogEntry))
var newListeners = make(map[ListenerHandle]func([]*LogEntry))

type LogEntry struct {
	Id        int
	Character string
	Server    string
	Timestamp time.Time
	Message   string
	Meaning   ParsedLog
}

func RegisterLogsListener(listener func([]*LogEntry)) ListenerHandle {
	curId := nextListenerHandler
	nextListenerHandler++
	newListeners[curId] = listener
	return curId
}

func (h ListenerHandle) Release() {
	delete(logListeners, h)
	delete(newListeners, h)
}

var cancelLogScans = func() {}

func RestartLogScans(baseCtx context.Context) {
	cancelLogScans()
	var ctx context.Context
	ctx, cancelLogScans = context.WithCancel(baseCtx)
	go readAllLogsLoop(ctx)
}

func readAllLogsLoop(ctx context.Context) {
	eqDir, _, _ := settings.LookupSetting(settings.EverQuestDirectory)
	eqLogDir := path.Join(eqDir, "Logs")
	seen := make(map[string]struct{})
	for {
		entries, err := fs.ReadDir(eqLogDir)
		if err == nil {
			for _, fi := range entries {
				if _, ok := seen[fi.Name()]; !ok {
					seen[fi.Name()] = struct{}{}
					if parts := filenameMatch.FindStringSubmatch(fi.Name()); parts != nil {
						go tailLog(ctx, path.Join(eqLogDir, fi.Name()), parts[1], parts[2])
					}
				}
			}
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(1000 * time.Millisecond):
		}
	}
}

// Provide unique identifers to log events
var logIdGen = 0

func tailLog(ctx context.Context, filename string, character string, server string) {
	fd, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		console.Log("Failed to open file")
		return
	}
	defer fd.Close()
	fd.Seek(-1, io.SeekEnd)
	rdbuf := make([]byte, 1024)
	buffer := new(bytes.Buffer)
	for {
		cnt, _ := fd.Read(rdbuf)
		if cnt <= 0 {
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
		buffer.Write(rdbuf[0:cnt])
		ib := bytes.IndexByte(buffer.Bytes(), '\n')
		if ib >= 0 {
			buffer.ReadBytes('\n')
			break
		}
	}

	for {
		for k, v := range newListeners {
			logListeners[k] = v
			delete(newListeners, k)
		}
		var entries []*LogEntry
		for ib := bytes.IndexByte(buffer.Bytes(), '\n'); ib >= 0; ib = bytes.IndexByte(buffer.Bytes(), '\n') {
			line, _ := buffer.ReadString('\n')
			line = strings.ReplaceAll(line, "\r", "")
			line = strings.ReplaceAll(line, "\n", "")
			if parts := loglineMatch.FindStringSubmatch(line); parts != nil {
				parsedTime, _ := time.Parse(time.ANSIC, parts[1])
				interpRaw := logterpreter.Dispatch(parts[2])
				var interp ParsedLog
				var ok bool
				if interp, ok = interpRaw.(ParsedLog); interp != nil && ok {
					interp = interp.Visit(substituteYouHandler{charName: character}).(ParsedLog)
				}
				entry := &LogEntry{
					Id:        logIdGen,
					Character: character,
					Server:    server,
					Timestamp: parsedTime,
					Message:   parts[2],
					Meaning:   interp,
				}
				logIdGen++
				entries = append(entries, entry)
			}
		}
		if len(entries) != 0 {
			for _, callback := range logListeners {
				callback(entries)
			}
		}
		cnt, _ := fd.Read(rdbuf)
		if cnt <= 0 {
			select {
			case <-ctx.Done():
				return
			case <-time.After(10 * time.Millisecond):
				continue
			}
		}
		buffer.Write(rdbuf[0:cnt])
	}
}

var logterpreter = handleChat(handleZone(handleHeal(handleDamage(handleDeath(multipattern.New())))))

type substituteYouHandler struct {
	charName string
}

func (s substituteYouHandler) OnZone(log *ZoneLog) interface{} {
	return log
}

func (s substituteYouHandler) OnDamage(log *DamageLog) interface{} {
	if log.Source == "You" {
		log.Source = s.charName
	}
	if log.Target == "You" {
		log.Target = s.charName
	}
	return log
}

func (s substituteYouHandler) OnHeal(log *HealLog) interface{} {
	if log.Source == "You" {
		log.Source = s.charName
	}
	if log.Target == "You" {
		log.Target = s.charName
	}
	return log
}

func (s substituteYouHandler) OnDeath(log *DeathLog) interface{} {
	if log.Source == "You" {
		log.Source = s.charName
	}
	if log.Target == "You" {
		log.Target = s.charName
	}
	return log
}

func (s substituteYouHandler) OnChat(log *ChatLog) interface{} {
	if log.Source == "You" {
		log.Source = s.charName
	}
	return log
}
