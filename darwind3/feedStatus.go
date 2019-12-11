package darwind3

import (
	"log"
	"time"
)

// FeedStatus Manages the state of the current feed, detecting if the feed is missing messages.
// If missing messages have been detected then we:
// * clear their database;
// * download and process the snapshot;
// * download and process the pPort log file entries until they reach the timestamp of the first message they received
// through their topic after reconnection;
// * resume processing messages from their topic
type FeedStatus struct {
	SequenceNumber int32
	TS             time.Time
	d3             *DarwinD3
	initialized    bool
	entries        []logEntry
}

// process checks the message SequenceNumber and if we are missing any then holds the feed whilst it checks for a
// new Snapshot and retrieves the pending logs.
// There is a chance it still misses a recent message if that message hasn't yet been put into the remote log files
// however this will pickup missing messages if there was a longer outage, network disconnection etc.
//
// Also, on startup this will also force a retrieval so we are in a consistent state.
func (fs *FeedStatus) process(p *Pport) {
	// Do nothing for sR entries or entries with no headers
	if p.SnapshotUpdate || p.FeedHeaders.SequenceNumber < 0 {
		return
	}

	if fs.snapshotRequired(p) {
		log.Println("Sequence Mismatch", p.FeedHeaders.SequenceNumber, fs.SequenceNumber)
		err := fs.loadSnapshot(p.TS)
		if err != nil {
			log.Println("Error", err)
		}
	}

	// Update state to latest message
	fs.SequenceNumber = p.FeedHeaders.SequenceNumber
	fs.TS = p.TS
}

func (fs *FeedStatus) snapshotRequired(p *Pport) bool {
	// Ensures we always require a snapshot on startup
	if !fs.initialized {
		fs.initialized = true
		return true
	}

	lastId := fs.SequenceNumber
	nextId := p.FeedHeaders.SequenceNumber

	// SequenceNumber has rolled over
	if lastId == 9999999 && nextId == 0 {
		return false
	}

	// Difference must be 1
	return (nextId - lastId) != 1
}
