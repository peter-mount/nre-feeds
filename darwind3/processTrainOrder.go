package darwind3

import (
	"log"
	"time"
)

// Process processes an inbound set of TrainOrders and applies them to the
// relevant schedules
func (to *trainOrderWrapper) Process(tx *Transaction) error {

	// No order data then ignore
	if to.Set == nil {
		return nil
	}

	if to.Set.First != nil {
		if err := to.processOrder(tx, 1, to.Set.First); err != nil {
			return err
		}
	}

	if to.Set.Second != nil {
		if err := to.processOrder(tx, 2, to.Set.Second); err != nil {
			return err
		}
	}

	if to.Set.Third != nil {
		if err := to.processOrder(tx, 3, to.Set.Third); err != nil {
			return err
		}
	}

	return nil
}

// Processes a specific TrainOrderItem
func (to *trainOrderWrapper) processOrder(tx *Transaction, order int, tod *trainOrderItem) error {
	if tod.RID.RID != "" {
		return to.processOrderRID(tx, order, tod)
	}

	if tod.TrainId != "" {
		log.Printf("trainOrder with trainId \"%s\" currently unsupported\n", tod.TrainId)
	} else {
		log.Println("trainOrder with no rid or trainId received")
	}

	return nil
}

func (tod *trainOrderItem) apply(to *trainOrderWrapper, order int, ts time.Time, sched *Schedule) bool {
	for _, l := range sched.Locations {
		if l.Tiploc == to.Tiploc && l.Times.EqualInSchedule(&tod.RID.Times) {
			if to.Clear {
				l.Forecast.TrainOrder = nil
			} else {
				l.Forecast.TrainOrder = &TrainOrder{
					Order:    order,
					Platform: to.Platform,
					Date:     ts,
				}
				l.updated = true
			}
			return true
		}
	}
	return false
}

func (to *trainOrderWrapper) processOrderRID(tx *Transaction, order int, tod *trainOrderItem) error {

	// Retrieve the schedule to be updated
	sched := tx.d3.GetSchedule(tod.RID.RID)

	// No schedule then try to fetch it from the timetable
	if sched == nil {
		log.Println("TrainOrder: Resolving schedule", tod.RID.RID)
		sched = tx.ResolveSchedule(tod.RID.RID)
	}

	// Still no schedule then We've got a TS for a train with no known schedule so give up
	if sched == nil {
		log.Println("TrainOrder: Failed to resolve schedule", tod.RID.RID)
		return nil
	}

	sched.UpdateTime()

	if tod.apply(to, order, tx.pport.TS, sched) {
		sched.Date = tx.pport.TS
		if tx.d3.PutSchedule(sched) {
			tx.d3.EventManager.PostEvent(&DarwinEvent{
				Type:     Event_ScheduleUpdated,
				RID:      sched.RID,
				Schedule: sched,
			})
		}
	} else {
		/* Testing shows this is true, in some schedules we don't actually find the tiploc
		   log.Printf(
		     "Failed to apply TrainOrder %d %s %s %s %v %s\n",
		     order,
		     tod.RID.RID,
		     to.Tiploc,
		     to.CRS,
		     to.Clear,
		     tod.RID.Times.String(),
		   )
		*/
	}

	return nil
}
