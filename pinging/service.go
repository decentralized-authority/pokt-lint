package pinging

import (
	"context"
	"fmt"
	"github.com/go-ping/ping"
)

const (
	defaultNumPings = 5
)

type Service interface {
	PingHost(ctx context.Context) (*ping.Statistics, error)
	SetNumPings(ctx context.Context, num int)
}

func NewService(hostname string) (Service, error) {
	pinger, err := ping.NewPinger(hostname)
	if err != nil {
		return &service{}, nil
	}

	pinger.Count = defaultNumPings
	return &service{pinger: *pinger}, nil
}

type service struct {
	pinger ping.Pinger
}

func (s *service) PingHost(ctx context.Context) (*ping.Statistics, error) {
	err := s.pinger.Run()
	if err != nil {
		return nil, fmt.Errorf("PingHost: %s", err)
	}

	return s.pinger.Statistics(), nil
}

func (s *service) SetNumPings(ctx context.Context, num int) {
	s.pinger.Count = num
}
