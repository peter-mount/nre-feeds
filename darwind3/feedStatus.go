package darwind3

import (
	"log"
	"os"
	"sync"
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
	mutex          sync.Mutex
	d3             *DarwinD3
	snapshotTime   time.Time
}

func (fs *FeedStatus) process(p *Pport) {
	// Do nothing for sR entries
	if p.SnapshotUpdate {
		return
	}

	//fs.mutex.Lock()
	//defer fs.mutex.Unlock()

	if fs.snapshotRequired(p) {
		log.Println("Sequence Mismatch", p.FeedHeaders.SequenceNumber, fs.SequenceNumber)
		err := fs.loadSnapshot(p.TS, &fs.mutex)
		if err != nil {
			// Treat this as terminal
			log.Println("Error", err)
			os.Exit(1)
		}
	}

	// Update state to latest message
	fs.SequenceNumber = p.FeedHeaders.SequenceNumber
	fs.TS = p.TS
}

func (fs *FeedStatus) snapshotRequired(p *Pport) bool {
	lastId := fs.SequenceNumber
	nextId := p.FeedHeaders.SequenceNumber

	// SequenceNumber has rolled over
	if lastId == 9999999 && nextId == 0 {
		return false
	}

	// Difference must be 1
	return (nextId - lastId) != 1
}
