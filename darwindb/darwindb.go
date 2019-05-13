package darwindb

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/peter-mount/golib/statistics"
	"github.com/peter-mount/nre-feeds/bin"
)

type DarwinDB struct {
	db                *sql.DB
	scheduleStatement *sql.Stmt

	updateChannel chan []byte
}

func (d *DarwinDB) Init(cfg *bin.Config) error {
	db, err := sql.Open("postgres", cfg.DB.PostgresUri)
	if err != nil {
		return err
	}
	d.db = db

	stmt, err := db.Prepare("select darwin.updateschedule( $1 )")
	if err != nil {
		return err
	}
	d.scheduleStatement = stmt

	// The update channel so all calls to scheduleStatement are done in sequence
	d.updateChannel = make(chan []byte, 1)
	go func() {
		for {
			msg := <-d.updateChannel
			_, err := d.scheduleStatement.Exec(string(msg))
			if err == nil {
				statistics.Incr("darwin.db.updated.success")
			} else {
				statistics.Incr("darwin.db.updated.error")
			}
		}
	}()

	return nil
}

func (d *DarwinDB) Stop() {
	if d.db != nil {
		_ = d.db.Close()
		d.db = nil
	}
}

// Deactivated simply updates the schedule.
// It's a separate hook in case we need to do anything else
func (d *DarwinDB) Deactivated(msg []byte) {
	statistics.Incr("darwin.db.deactivated")
	d.updateChannel <- msg
}

// Store schedule updates in the db
func (d *DarwinDB) ScheduleUpdated(msg []byte) {
	statistics.Incr("darwin.db.schedule")
	d.updateChannel <- msg
}
