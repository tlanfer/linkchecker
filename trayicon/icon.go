package trayicon

import (
	_ "embed"
	"github.com/getlantern/systray"
	"time"
)

//go:embed fine.ico
var fine []byte

//go:embed yell.ico
var yell []byte

func New() Ui {
	u := Ui{
		quit: make(chan any),
	}
	go systray.Run(u.onReady, u.onExit)
	return u
}

type Ui struct {
	quit chan any
}

func (u *Ui) Broken() {
	systray.SetIcon(yell)
	go func() {
		time.Sleep(5 * time.Second)
		systray.SetIcon(fine)
	}()

}

func (u *Ui) onReady() {
	systray.SetTitle("LinkChecker")
	systray.SetIcon(fine)

	menuItem := systray.AddMenuItem("Quit", "Quit LinkChecker")

	<-menuItem.ClickedCh
	u.quit <- "quit"
	systray.Quit()
	time.Sleep(100 * time.Millisecond)
}

func (u *Ui) onExit() {}

func (u *Ui) Wait() {
	<-u.quit
}
