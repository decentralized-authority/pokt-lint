package http

import (
	"github.com/itsnoproblem/pokt-lint/timer"
	"log"
	"time"
)

const (
	defaultPingCount = 10
	pingDelayMS      = 500
)

// Pinger is a ping client capable of returning PingStats
type Pinger struct {
	Count  int64
	URL    string
	stats  *PingStats
	client Client
}

// PingStats are the measured timing and outcome of a ping test
type PingStats struct {
	NumSent   int64          `json:"num_sent"`
	NumOk     int64          `json:"num_ok"`
	MinTimeMS float64        `json:"min_time_ms"`
	MaxTimeMS float64        `json:"max_time_ms"`
	AvgTimeMS float64        `json:"avg_time_ms"`
	Results   []PingResponse `json:"results"`
}

type PingResponse struct {
	DurationMS float64 `json:"duration_ms"`
	StatusCode int     `json:"status_code"`
}

// NewPinger returns a new Pinger for the specified url
func NewPinger(client Client, url string) *Pinger {
	return &Pinger{
		Count:  defaultPingCount,
		URL:    url,
		client: client,
	}
}

// Run executes the ping tests
func (p *Pinger) Run() (*PingStats, error) {
	var total, min, max time.Duration
	var results = make([]PingResponse, p.Count)
	success := int64(0)

	for i := int64(0); i < p.Count; i++ {
		t := timer.Start()
		resp, err := p.client.Get(p.URL)
		if err != nil {
			log.Default().Printf("Ping %s: %s", p.URL, err)
			continue
		}
		duration := t.Elapsed()
		total += duration

		results[i] = PingResponse{
			DurationMS: float64(duration.Microseconds()) / 1000,
			StatusCode: resp.StatusCode,
		}

		if min == 0 || duration < min {
			min = duration
		}

		if duration > max {
			max = duration
		}

		if resp.StatusCode == 200 {
			success++
		}

		log.Default().Printf("Ping %s: %d (%s ms)", p.URL, resp.StatusCode, duration.String())
		time.Sleep(pingDelayMS * time.Millisecond)
	}

	p.stats = &PingStats{
		NumSent:   p.Count,
		NumOk:     success,
		MinTimeMS: float64(min.Microseconds()) / 1000,
		MaxTimeMS: float64(max.Microseconds()) / 1000,
		AvgTimeMS: float64(total.Microseconds()/p.Count) / 1000,
		Results:   results,
	}

	return p.stats, nil
}

// Statistics returns a reference to the PingStats from the pinger's most recent run.  If it hasn't run yet, the return value will be nil
func (p *Pinger) Statistics() *PingStats {
	return p.stats
}
