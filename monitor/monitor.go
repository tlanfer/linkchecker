package monitor

import (
	"fmt"
	"github.com/go-ping/ping"
	"log"
	"pinglog/logger"
	"pinglog/trayicon"
	"time"
)

func New(host string, interval time.Duration, logger *logger.Logger, ui trayicon.Ui) *Monitor {
	return &Monitor{
		host:     host,
		interval: interval,
		log:      logger,
		stop:     make(chan any, 1),
		ui:       ui,
	}
}

type Monitor struct {
	host     string
	stop     chan any
	log      *logger.Logger
	interval time.Duration
	ui       trayicon.Ui
}

func (m *Monitor) Start() {
	go func() {
		ticker := time.NewTicker(m.interval)
		for {
			select {
			case <-ticker.C:
				stats, err := getStats(m.host)
				if err != nil {
					log.Println("Pinger failed:", err)
				}
				m.log.Log(m.host, stats.PacketLoss, stats.AvgRtt.Milliseconds())
				if stats.PacketLoss > 0 {
					m.ui.Broken()
				}
			case <-m.stop:
				ticker.Stop()
				return
			}
		}
	}()
}

func (m *Monitor) Stop() {
	m.stop <- struct{}{}
}

func getStats(host string) (*ping.Statistics, error) {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		return nil, fmt.Errorf("failed create pinger: %w", err)
	}

	pinger.SetPrivileged(true)
	pinger.Count = 3
	pinger.Timeout = 1 * time.Second

	err = pinger.Run()
	if err != nil {
		return nil, err
	}
	stats := pinger.Statistics()
	return stats, nil
}
