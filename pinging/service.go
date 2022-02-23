package pinging

import (
	"context"
	"fmt"
	nethttp "net/http"

	"github.com/itsnoproblem/pokt-lint/http"
)

const (
	defaultNumPings = 5
)

type Service interface {
	PingHost(ctx context.Context) (*http.PingStats, error)
}

func NewService(client nethttp.Client, url string) (Service, error) {
	pinger := http.NewPinger(client, url)
	return &service{pinger: pinger}, nil
}

type service struct {
	pinger *http.Pinger
}

func (s *service) PingHost(ctx context.Context) (*http.PingStats, error) {
	stats, err := s.pinger.Run()
	if err != nil {
		return nil, fmt.Errorf("PingHost: %s", err)
	}
	return stats, nil
}

func (s *service) SetNumPings(ctx context.Context, num int64) {
	s.pinger.Count = num
}
