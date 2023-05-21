package manager

import (
	"context"
	"time"
	"fmt"

	rancherusagerecordv1 "github.com/SUSE-Enceladus/csp-rancher-usage-operator/api/record/v1"
	rancherusage "github.com/SUSE-Enceladus/csp-rancher-usage-operator/generated/clientset/versioned"
	"github.com/SUSE-Enceladus/csp-rancher-usage-operator/pkg/clients/k8s"
	"github.com/SUSE-Enceladus/csp-rancher-usage-operator/pkg/metrics"
//	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/sirupsen/logrus"
)

type UsageOperator struct {
	cancel  context.CancelFunc
	k8s     k8s.Client
	ruo     *rancherusage.Clientset
	scraper metrics.Scraper
}

func NewUsageOperator(k k8s.Client, ruo *rancherusage.Clientset, s metrics.Scraper) *UsageOperator {
	return &UsageOperator{
		k8s:     k,
		ruo:     ruo,
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
                }
                logrus.Infof("Node Count %d", ncount)
		//updateRancherUsageRecord(ctx, m.ruo, 3)
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

func updateRancherUsageRecord(ctx context.Context, ruo *rancherusage.Clientset, managedNodeCount uint32) {
	rur := &rancherusagerecordv1.RancherUsageRecord{
		ObjectMeta: metav1.ObjectMeta{
			Name: "gyee-test",
		},
		BaseProduct: "blah",
		ManagedNodeCount: 800,
		ReportingTime: "sfsfsdfd",
	}
	_, err := ruo.RecordV1().RancherUsageRecords().Create(ctx, rur, metav1.CreateOptions{})
	if err != nil {
                logrus.Errorf("unable to create usage record %s", err)
        }
	/*ur, err := ruo.RecordV1().RancherUsageRecords().Get(ctx, "rancher-usage-record", metav1.GetOptions{})
	if err != nil {
		logrus.Errorf("unable to lookup usage record %s", err)
	}
	if ur != nil {
		logrus.Infof("%+v", ur)
	} else {
		logrus.Infof("Rancher usage record not found, creating a new one")
	}*/
}

func ticker(ctx context.Context, duration time.Duration) <-chan time.Time {
	ticker := time.NewTicker(duration)
	go func() {
		<-ctx.Done()
		ticker.Stop()
	}()
	return ticker.C
}
