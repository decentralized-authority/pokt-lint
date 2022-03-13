package timer

import "time"

// Timer is a utility used to track elapsed time
type Timer struct {
	startedAt time.Time
}

// Elapsed returns the time that has passed since the timer was started
func (t *Timer) Elapsed() time.Duration {
	return time.Since(t.startedAt)
}

// Start returns an activated Timer
func Start() Timer {
	return Timer{
		startedAt: time.Now(),
	}
}
