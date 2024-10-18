package service

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gabriel-panz/gomato/types"
	"github.com/gen2brain/beeep"
	"github.com/schollz/progressbar/v3"
)

type PomodoroStatus uint8

const (
	Working PomodoroStatus = iota
	Resting
)

type Timer struct {
	working bool
	status  chan PomodoroStatus
	tm      *time.Timer
	ti      *time.Ticker
}

func CreateTimer() *Timer {
	return &Timer{
		working: false,
		status:  make(chan PomodoroStatus),
	}
}

func (p *Timer) StartPomodoro(c *types.TimerConfig) {
	p.tm = time.NewTimer(c.WorkTime)
	p.ti = time.NewTicker(time.Second)
	p.working = true

	go renderProgress(p, c)

	status := Working
	for {
		<-p.tm.C
		if status == Working {
			err := beeep.Alert("Go-mato!", "Pause, you can rest now.", "")
			if err != nil {
				panic(err)
			}

			p.tm.Reset(c.PauseTime)

			status = Resting
		} else {
			err := beeep.Alert("Go-mato!", "Time to work!", "")
			if err != nil {
				panic(err)
			}

			p.tm.Reset(c.WorkTime)
			status = Working
		}
		p.status <- status
	}
}

func (p *Timer) StopPomodoro() {
	log.Println("pomodoro stopping")
}

func renderProgress(p *Timer, c *types.TimerConfig) {
	pbar := progressbar.Default(int64(c.WorkTime.Seconds()))
	for {
		<-p.ti.C
		pbar.Add(1)
		if pbar.IsFinished() {
			s := <-p.status
			if s == Working {
				pbar.ChangeMax64(int64(c.WorkTime.Seconds()))
			} else {
				pbar.ChangeMax64(int64(c.PauseTime.Seconds()))
			}
			pbar.Reset()
			pbar.RenderBlank()
		}
	}
}

func (p *Timer) StartFlow() {
	p.ti = time.NewTicker(time.Second)
	done := make(chan bool)
	startT := time.Now()
	var cT time.Duration
	resting := false

	go waitInput(done)

	for {
		select {
		case <-done:
			if resting {
				resting = false
				startT = time.Now()
				fmt.Printf("\r%s", time.Since(startT))
			} else {
				resting = true
				cT = time.Since(startT) / 5
				fmt.Printf("\r%s", cT.Round(time.Second))
			}
			go waitInput(done)
		case <-p.ti.C:
			if resting {
				cT -= time.Second
				cT = cT.Round(time.Second)
			} else {
				cT = time.Since(startT)
				cT = cT.Round(time.Second)
			}
			fmt.Printf("\r%s", cT)
		}
	}
}

func waitInput(done chan bool) {
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	done <- true
}
