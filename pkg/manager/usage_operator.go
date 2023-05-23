package manager

import (
	"context"
	"time"
	"fmt"

	"github.com/SUSE-Enceladus/csp-rancher-usage-operator/pkg/clients/k8s"
	"github.com/SUSE-Enceladus/csp-rancher-usage-operator/pkg/metrics"
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
		ncount, err := m.getNodeCount(ctx)
                if err != nil {
                       errs <- err
		       logrus.Warningf("Unable to get managed node count. Will not update usage record.")
                } else {
                	logrus.Infof("Node Count %d", ncount)
			err = m.k8s.UpdateProductUsage(uint32(ncount))
                }
	}
	logrus.Infof("[manager] exiting")
}

func (m *UsageOperator) getNodeCount(ctx context.Context) (int ,error) {
        nodeCounts, err := m.scraper.ScrapeAndParse()
        if err != nil {
                return 0, fmt.Errorf("unable to determine number of active nodes: %v", err)
        }
        logrus.Debugf("found %d nodes from rancher metrics", nodeCounts.Total)
        return nodeCounts.Total, nil
}

func ticker(ctx context.Context, duration time.Duration) <-chan time.Time {
	ticker := time.NewTicker(duration)
	go func() {
		<-ctx.Done()
		ticker.Stop()
	}()
	return ticker.C
}
