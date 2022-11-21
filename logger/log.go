package logger

import (
	"fmt"
	"os"
	"time"
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

	files := map[string]*os.File{}

	go func() {
		for {
			select {
			case e := <-l.events:
				l, exists := files[e.Host]
				if !exists {
					f, err := open(e.Host)
					if err != nil {
						panic(err)
					}
					l = f
				}
				timestamp := time.Now().Format("2006-01-02 15:04:05")
				fmt.Fprintf(l, "%v,%v,%v,%v\n", timestamp, e.Host, e.PacketLoss, e.Rtt)
			case <-l.stop:
				for _, file := range files {
					file.Close()
				}
				return
			}
		}
	}()

	return nil
}

func open(hostname string) (*os.File, error) {
	file, err := os.OpenFile(fmt.Sprintf("monitor_%v.log", hostname), os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	return file, nil
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
