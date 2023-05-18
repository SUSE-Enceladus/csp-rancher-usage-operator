package manager

import (
	"context"
	"time"

	"github.com/rancher/csp-rancher-usage-operator/pkg/clients/k8s"
	"github.com/rancher/csp-rancher-usage-operator/pkg/metrics"
	"github.com/sirupsen/logrus"
)

type UsageOperator struct {
	cancel  context.CancelFunc
	k8s     k8s.Client
	scraper metrics.Scraper
}

func NewUsageOperator(k k8s.Client, s metrics.Scraper) *UsageOperator {
	return &UsageOperator{
		k8s:     k,
		scraper: s,
	}
}

func (m *UsageOperator) Start(ctx context.Context, errs chan<- error) {
	go m.start(ctx, errs)
}

const (
	managerInterval = 30 * time.Second
)

func (m *UsageOperator) start(ctx context.Context, errs chan<- error) {
	for range ticker(ctx, managerInterval) {
		// TODO: calculate nodes here
		logrus.Info("Calculating managed nodes")
	}
	logrus.Infof("[manager] exiting")
}

func ticker(ctx context.Context, duration time.Duration) <-chan time.Time {
	ticker := time.NewTicker(duration)
	go func() {
		<-ctx.Done()
		ticker.Stop()
	}()
	return ticker.C
}
