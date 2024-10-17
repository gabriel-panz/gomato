package types

import "time"

type TimerConfig struct {
	Id        int64
	Name      string
	WorkTime  time.Duration
	PauseTime time.Duration
	NotificationLevel
}

type NotificationLevel int

const (
	None NotificationLevel = iota
	Audio
	AudioVisual
)
