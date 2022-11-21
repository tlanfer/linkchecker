package logger

import (
	"fmt"
	"os"
	"pinglog/config"
	"time"
)

func New(prefix string, threshold config.Threshold) *Logger {
	return &Logger{
		prefix:    prefix,
		threshold: threshold,
		stop:      make(chan any),
		events:    make(chan event, 1),
	}
}

type Logger struct {
	prefix    string
	stop      chan any
	events    chan event
	threshold config.Threshold
}

type event struct {
	Host       string
	PacketLoss float64
	Rtt        time.Duration
}

func (l *Logger) Start() error {

	files := map[string]*os.File{}

	go func() {
		for {
			select {
			case e := <-l.events:
				file, exists := files[e.Host]
				if !exists {
					f, err := l.open(e.Host)
					if err != nil {
						panic(err)
					}
					file = f
				}

				if e.PacketLoss >= l.threshold.PacketLoss || e.Rtt >= l.threshold.Rtt {
					timestamp := time.Now().Format("2006-01-02 15:04:05")
					fmt.Fprintf(file, "%v,%v,%v,%v\n", timestamp, e.Host, e.PacketLoss, e.Rtt.Milliseconds())
				}
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

func (l *Logger) open(hostname string) (*os.File, error) {
	file, err := os.OpenFile(fmt.Sprintf("%v%v.log", l.prefix, hostname), os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	return file, nil
}

func (l *Logger) Stop() {
	l.stop <- struct{}{}
}

func (l *Logger) Log(host string, packetLoss float64, rtt time.Duration) {
	l.events <- event{
		Host:       host,
		PacketLoss: packetLoss,
		Rtt:        rtt,
	}
}
