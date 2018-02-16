package domain

import (
	"fmt"
	"log"
	"time"
)

func (s *server) startScheduler() {

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

	// fetch end of day prices NOW
	s.dailyUpdate()

	// fetch end of day prices at end of day every 24 hours
	go func() {
		now := time.Now()
		next := now.Truncate(time.Hour*24).AddDate(0, 0, 1)
		delay := next.Sub(now)
		s.log(fmt.Sprintf("Next daily update in %s", delay.String()))

		timer := time.NewTimer(delay)

		<-timer.C

		// first update
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
		s.log(fmt.Sprintf("ERROR: updating prices - %s", err))
		return
	}

	// update portfolios
	if err := s.updatePortfolios(); err != nil {
		// log error
		s.log(fmt.Sprintf("ERROR: updating portfolios - %s", err))
		return
	}

	if err := s.saveMetrics(); err != nil {
		// log error
		s.log(fmt.Sprintf("ERROR: saving portfolios - %s", err))
		return
	}

}

// dailyUpdate - runs daily
func (s *server) dailyUpdate() {
	s.log("Daily update running")
	if err := DefaultArchive.UpdateDaySummaries(); err != nil {
		s.log(fmt.Sprintf("ERROR: updating closing prices - %s", err))
	}

}
