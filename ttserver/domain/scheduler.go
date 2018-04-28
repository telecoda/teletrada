package domain

import (
	"fmt"
	"log"
	"time"
)

func (s *server) startScheduler() {

	// fetch end of day prices NOW
	log.Printf("Initialising prices at server start")
	s.dailyUpdate()
	log.Printf("Initialisation complete")

	// this update runs every x seconds
	go func(frequency time.Duration) {
		updateTicker := time.NewTicker(frequency)
		defer updateTicker.Stop()

		for {
			select {
			case <-s.stopUpdate:
				log.Printf("Scheduled price update stopping.")
				return
			case <-updateTicker.C:
				s.scheduledUpdate()
			}
		}
	}(s.config.UpdateFreq)

	// fetch end of day prices at end of day every 24 hours
	go func() {
		now := ServerTime()
		next := now.Truncate(time.Hour*24).AddDate(0, 0, 1)
		delay := next.Sub(now)
		DefaultLogger.log(fmt.Sprintf("Next daily update in %s", delay.String()))

		timer := time.NewTimer(delay)

		<-timer.C

		// update at end of first day
		s.dailyUpdate()

		frequency := time.Duration(24 * time.Hour)

		updateTicker := time.NewTicker(frequency)
		defer updateTicker.Stop()

		for {
			select {
			case <-s.stopUpdate:
				log.Printf("Scheduled price update stopping.")
				return
			case <-updateTicker.C:
				// update every 24 hours
				s.dailyUpdate()
			}
		}
	}()

}

func (s *server) stopScheduler() {
	s.Lock()
	s.stopUpdate <- true
	s.Unlock()
}

// scheduledUpdate - runs every x seconds
func (s *server) scheduledUpdate() {

	// Update latest prices
	if err := DefaultArchive.UpdatePrices(); err != nil {
		// log error
		DefaultLogger.log(fmt.Sprintf("ERROR: updating prices - %s", err))
	}

	// update portfolios
	if err := s.updatePortfolios(); err != nil {
		// log error
		DefaultLogger.log(fmt.Sprintf("ERROR: updating portfolios - %s", err))
	}

	if err := s.saveMetrics(); err != nil {
		// log error
		DefaultLogger.log(fmt.Sprintf("ERROR: saving portfolios - %s", err))
	}

}

// dailyUpdate - runs daily
func (s *server) dailyUpdate() {
	DefaultLogger.log("started Daily update")
	if err := DefaultArchive.UpdatePrices(); err != nil {
		// log error
		DefaultLogger.log(fmt.Sprintf("ERROR: updating prices - %s", err))
		return
	}
	if err := DefaultArchive.UpdateDaySummaries(); err != nil {
		DefaultLogger.log(fmt.Sprintf("ERROR: updating closing prices - %s", err))
	}
	DefaultLogger.log("ended Daily update")

}
