package pinging

import (
	"context"
	"fmt"
	"github.com/itsnoproblem/pokt-lint/http"
	"regexp"
)

// Service - pinging service performs ping tests
type Service interface {
	PingHost(ctx context.Context) (*http.PingStats, error)
	SetNumPings(ctx context.Context, num int64)
}

// NewService returns a new pinging service
func NewService(client http.Client, url string) (Service, error) {

	slashPatt := regexp.MustCompile(`/+$`)
	url = slashPatt.ReplaceAllString(url, "")

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
