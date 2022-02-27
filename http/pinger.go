package http

import (
	"fmt"
	"github.com/itsnoproblem/pokt-lint/timer"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultPingCount = 10
	pingDelayMS      = 500
)

type Pinger struct {
	Count  int64
	URL    string
	stats  *PingStats
	client http.Client
}

type PingStats struct {
	NumSent   int64   `json:"num_sent"`
	NumOk     int64   `json:"num_ok"`
	MinTimeMS float64 `json:"min_time_ms"`
	MaxTimeMS float64 `json:"max_time_ms"`
	AvgTimeMS float64 `json:"avg_time_ms"`
}

func NewPinger(client http.Client, url string) *Pinger {
	return &Pinger{
		Count:  defaultPingCount,
		URL:    url,
		client: client,
	}
}

func (p *Pinger) Run() (*PingStats, error) {
	var total, min, max time.Duration
	success := int64(0)

	parsed, err := url.Parse(p.URL)
	if err != nil {
		return nil, fmt.Errorf("Pinger.Run: %s", err)
	}

	url := fmt.Sprintf("%s://%s%s", parsed.Scheme, parsed.Host, "/v1")

	for i := int64(0); i < p.Count; i++ {
		t := timer.Start()
		resp, err := p.client.Get(url)
		if err != nil {
			log.Default().Printf("Ping %s: %s", url, err)
			continue
		}
		duration := t.Elapsed()
		total += duration

		if min == 0 || duration < min {
			min = duration
		}

		if duration > max {
			max = duration
		}

		if resp.StatusCode == 200 {
			success++
		}

		log.Default().Printf("Ping %s: %d (%s ms)", url, resp.StatusCode, duration.String())
		time.Sleep(pingDelayMS * time.Millisecond)
	}

	p.stats = &PingStats{
		NumSent:   p.Count,
		NumOk:     success,
		MinTimeMS: float64(min.Microseconds()) / 1000,
		MaxTimeMS: float64(max.Microseconds()) / 1000,
		AvgTimeMS: float64(total.Microseconds()/p.Count) / 1000,
	}

	return p.stats, nil
}

func (p *Pinger) Statistics() *PingStats {
	return p.stats
}
