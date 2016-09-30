package floodcontrol

import (
	"fmt"
	"github.com/Nyubis/mibot/core"
	"time"
	"sync"
)

// Keep a record of events by type, nick and channel, and store the timestamp
// When a new event comes in, remove all records older than x seconds
// Count how many of the same type, nick, and channel are left
// Return true if this count exceeds a certain value (so the module will know to ignore it)

type record struct {
	nick      string
	channel   string
	timestamp time.Time
}

type mutexmap struct {
	sync.RWMutex
	// The string is the event type (e.g. "link", "invite")
	m map[string][]record
}

var recent mutexmap

// The time in seconds it remembers groups of events (the string is the type of event)
var memtime map[string]int

// The amount of allowed events in the above timespan
var maxcount map[string]int

func Init(_ *core.Bot) {
	recent = mutexmap{m: make(map[string][]record)}
	memtime = make(map[string]int)
	maxcount = make(map[string]int)
	// These should probably be read from the config...
	memtime["invite"] = 15
	maxcount["invite"] = 2
	memtime["link"] = 12
	maxcount["link"] = 3
}

func FloodCheck(event, nick, channel string) bool {
	now := time.Now()
	recent.Lock()
	recent.m[event] = append(recent.m[event], record{nick, channel, now})
	cutofftime := now.Add(time.Duration(-1*memtime[event]) * time.Second)
	slice := removeBefore(event, cutofftime)
	recent.Unlock()

	if len(slice) > maxcount[event] {
		fmt.Printf("%s in %s has used %s %d times in the past %d seconds\n", nick, channel, event, len(slice), maxcount[event])
	}

	return len(slice) > maxcount[event]
}

// Remove the values before the given timestamp
func removeBefore(event string, timestamp time.Time) []record {
	slice, ok := recent.m[event]
	if !ok {
		// There aren't any to remove
		return []record{}
	}

	var i int
	var entry record
	for i, entry = range slice {
		if timestamp.Before(entry.timestamp) {
			break
		}
	}

	recent.m[event] = slice[i:]
	return slice[i:]
}
