package logger

import (
	"fmt"
	"os"
)

func New(filename string) *Logger {
	return &Logger{
		filename: filename,
		stop:     make(chan any),
		events:   make(chan event, 1),
	}
}

type Logger struct {
	filename string
	stop     chan any
	events   chan event
}

type event struct {
	Host       string
	PacketLoss float64
	Rtt        int64
}

func (l *Logger) Start() error {
	file, err := os.OpenFile(l.filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	go func() {
		for {
			select {
			case e := <-l.events:
				fmt.Fprintf(file, "%v,%v,%v\n", e.Host, e.PacketLoss, e.Rtt)
			case <-l.stop:
				file.Close()
				return
			}
		}
	}()

	return nil
}

func (l *Logger) Stop() {
	l.stop <- struct{}{}
}

func (l *Logger) Log(host string, packetLoss float64, rtt int64) {
	l.events <- event{
		Host:       host,
		PacketLoss: packetLoss,
		Rtt:        rtt,
	}
}
