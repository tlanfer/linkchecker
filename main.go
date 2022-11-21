package main

import (
	"pinglog/config"
	"pinglog/logger"
	"pinglog/monitor"
	"pinglog/trayicon"
)

func main() {
	c := config.Load()
	ui := trayicon.New()
	l := logger.New(c.LogFile)
	l.Start()

	var ms []*monitor.Monitor
	for _, host := range c.Hosts {
		m := monitor.New(host, c.Interval, l, ui)
		m.Start()
		ms = append(ms, m)
	}

	ui.Wait()

	for _, m := range ms {
		m.Stop()
	}

	l.Stop()
}
