package darwinkb

import (
	"github.com/peter-mount/go-kernel"
	"github.com/peter-mount/go-kernel/bolt"
	"github.com/peter-mount/go-kernel/cron"
	"github.com/peter-mount/nre-feeds/bin"
	"os"
)

type DarwinKB struct {
	db     map[string]*bolt.BoltService
	config *bin.Config
	cron   *cron.CronService

	token KBToken
}

func (r *DarwinKB) Name() string {
	return "DarwinKB"
}

func (a *DarwinKB) Init(k *kernel.Kernel) error {

	service, err := k.AddService(&bin.Config{})
	if err != nil {
		return err
	}
	a.config = (service).(*bin.Config)

	// We use a separate DB for each bucket for mainenance reasons.
	// see PostInit but use a dummy filename here
	a.db = make(map[string]*bolt.BoltService)
	for _, bucket := range []string{
		incidentsBucket,
		serviceIndicatorsBucket,
		stationsBucket,
		ticketTypesBucket,
		tocsBucket,
	} {
		service, err = k.AddService(&bolt.BoltService{FileName: bucket})
		if err != nil {
			return err
		}
		a.db[bucket] = (service).(*bolt.BoltService)
	}

	service, err = k.AddService(&cron.CronService{})
	if err != nil {
		return err
	}
	a.cron = (service).(*cron.CronService)

	return nil
}

func (a *DarwinKB) PostInit() error {
	if a.config.Database.KB == "" {
		a.config.Database.KB = "/database/"
	}

	err := os.MkdirAll(a.config.Database.KB+"static/", 0x755)
	if err != nil {
		return err
	}

	// Here we set the db filename for each bucket. This will work as the db isn't stated yet
	for k, v := range a.db {
		v.FileName = a.config.Database.KB + "dwkb_" + k + ".db"
	}

	return nil
}

func (a *DarwinKB) Start() error {

	// Ensure the buckets exist in each db
	for bucket, db := range a.db {
		err := db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(bucket)
			return err
		})
		if err != nil {
			return err
		}
	}

	// Now add the cron jobs for each feed.
	a.addJob(stationsSchedule, a.refreshStations)
	a.addJob(incidentsSchedule, a.refreshIncidents)
	a.addJob(tocsSchedule, a.refreshCompanies)
	a.addJob(serviceIndicatorsSchedule, a.refreshServiceIndicators)
	a.addJob(ticketTypesSchedule, a.refreshTicketTypes)

	return nil
}

// addJob adds the job to cron and then runs it so we always call the feed on startup
func (a *DarwinKB) addJob(schedule string, f func()) {
	_, _ = a.cron.AddFunc(schedule, f)
	f()
}
