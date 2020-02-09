package floodcontrol

import (
	"fmt"
	"sync"
	"time"

	"github.com/Nyubis/mibot/core"
)

// Keep a record of events by type and channel; store the timestamp
// When a new event comes in, remove all records older than x seconds
// Count how many of the same type and channel are left
// Return true if this count exceeds a certain value (so the module will know to ignore it)

type record struct {
	channel   string
	timestamp time.Time
}

type lockedList struct {
	sync.RWMutex
	list []record
}

var recent map[string]*lockedList

// The time in seconds it remembers groups of events (the string is the type of event)
var memtime map[string]int

// The amount of allowed events in the above timespan
var maxcount map[string]int

func Init(_ *core.Bot) {
	recent = make(map[string]*lockedList)
	recent["invite"] = &lockedList{list: make([]record, 0)}
	recent["link"] = &lockedList{list: make([]record, 0)}
	recent["reply"] = &lockedList{list: make([]record, 0)}
	memtime = make(map[string]int)
	maxcount = make(map[string]int)
	LoadCfg()
}

func LoadCfg() {
	memtime["invite"] = core.Config.FloodControl.Invite.Time
	maxcount["invite"] = core.Config.FloodControl.Invite.Max
	memtime["link"] = core.Config.FloodControl.Link.Time
	maxcount["link"] = core.Config.FloodControl.Link.Max
	memtime["reply"] = core.Config.FloodControl.Reply.Time
	maxcount["reply"] = core.Config.FloodControl.Reply.Max
}

func FloodCheck(event, channel string) bool {
	now := time.Now()
	recent[event].Lock()
	defer recent[event].Unlock()
	list := recent[event].list
	cutofftime := now.Add(time.Duration(-1*memtime[event]) * time.Second)
	slice := filterBefore(list, cutofftime)
	recent[event].list = append(slice, record{channel, now})

	if len(slice) > maxcount[event] {
		fmt.Printf("Users in %s have used %s %d times in the past %d seconds\n", channel, event, len(slice), memtime[event])
		return true
	}

	return false
}

// Remove the values before the given timestamp
func filterBefore(list []record, timestamp time.Time) []record {
	var i int
	var entry record
	for i, entry = range list {
		if timestamp.Before(entry.timestamp) {
			break
		}
	}

	return list[i:]
}
