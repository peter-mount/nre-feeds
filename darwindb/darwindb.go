package darwindb

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/peter-mount/golib/rabbitmq"
	"github.com/peter-mount/golib/statistics"
	"github.com/peter-mount/nre-feeds/bin"
	"log"
	"os"
)

type DarwinDB struct {
	db                          *sql.DB
	scheduleStatement           *sql.Stmt
	indexStatement              *sql.Stmt
	getServiceStatement         *sql.Stmt
	getStationServicesStatement *sql.Stmt
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

	stmt, err = db.Prepare("select darwin.indexservices() as processed")
	if err != nil {
		return err
	}
	d.indexStatement = stmt

	stmt, err = db.Prepare("select darwin.getservice( $1 )")
	if err != nil {
		return err
	}
	d.getServiceStatement = stmt

	stmt, err = db.Prepare("select * from darwin.getservices($1,date_trunc('hour',$2::timestamp with time zone))")
	if err != nil {
		return err
	}
	d.getStationServicesStatement = stmt

	return nil
}

func (d *DarwinDB) Stop() {
	if d.db != nil {
		_ = d.db.Close()
		d.db = nil
	}
}

func (d *DarwinDB) Subscribe(mq *rabbitmq.RabbitMQ, prefix, queueName string, f func([]byte), eventTypes ...string) error {

	// Queue prefix, try to use the local hostname (e.g. of the container)
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "error"
	} else {
		hostname = hostname
	}

	queueName = hostname + "." + prefix + "." + queueName

	if channel, err := mq.NewChannel(); err != nil {
		log.Println(err)
		return err
	} else {

		// Force prefetchCount to 1 so we don't get everything in one go
		_ = channel.Qos(1, 0, false)

		// Unlike the other services this is a durable queue
		_, _ = mq.QueueDeclare(channel, queueName, true, false, false, false, nil)

		for _, eventType := range eventTypes {
			routingKey := prefix + ".d3.event." + eventType
			_ = mq.QueueBind(channel, queueName, routingKey, "amq.topic", false, nil)
		}

		ch, _ := mq.Consume(channel, queueName, "DB Consumer "+queueName, true, true, false, false, nil)

		go func() {
			for {
				msg := <-ch
				f(msg.Body)
			}
		}()

		return nil
	}
}

// Store schedule updates in the db
func (d *DarwinDB) ScheduleUpdated(msg []byte) {
	statistics.Incr("darwin.db.schedule")
	_, err := d.scheduleStatement.Exec(string(msg))
	if err == nil {
		statistics.Incr("darwin.db.updated.success")
	} else {
		statistics.Incr("darwin.db.updated.error")
	}
}

// Invokes the indexing job
func (d *DarwinDB) IndexSchedules() {
	processedCount := 0

	err := d.indexStatement.QueryRow().Scan(&processedCount)
	if err != nil {
		log.Println(err)
	} else {
		statistics.Set("darwin.db.indexed", int64(processedCount))
	}

}
