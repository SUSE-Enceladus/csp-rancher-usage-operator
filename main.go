package main

import (
	"fmt"
	"os"

	"github.com/rancher/csp-rancher-usage-operator/pkg/clients/k8s"
	"github.com/rancher/csp-rancher-usage-operator/pkg/manager"
	"github.com/rancher/csp-rancher-usage-operator/pkg/metrics"
	"github.com/rancher/wrangler/pkg/k8scheck"
	"github.com/rancher/wrangler/pkg/ratelimit"
	"github.com/rancher/wrangler/pkg/signals"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
)

var (
	Version   = "dev"
	GitCommit = "HEAD"
)

func main() {
	if err := run(); err != nil {
		logrus.Fatalf("csp-rancher-usage-operator failed to run with error: %v", err)
	}
}

const (
	debugEnv = "CATTLE_DEBUG"
)

func run() error {
	if os.Getenv(debugEnv) == "true" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.Infof("csp-rancher-usage-operator version %s is starting", fmt.Sprintf("%s (%s)", Version, GitCommit))

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	cfg.RateLimiter = ratelimit.None
	ctx := signals.SetupSignalContext()
	err = k8scheck.Wait(ctx, *cfg)
	if err != nil {
		return err
	}

	k8sClients, err := k8s.New(ctx, cfg)
	if err != nil {
		return err
	}

	hostname, err := k8sClients.GetRancherHostname()
	if err != nil {
		return fmt.Errorf("failed to start, unable to get hostname: %v", err)
	}

	m := manager.NewUsageOperator(k8sClients, metrics.NewScraper(hostname, cfg))

	errs := make(chan error, 1)
	m.Start(ctx, errs)
	go func() {
		for err := range errs {
			logrus.Errorf("usage operator manager error: %v", err)
		}
	}()

	<-ctx.Done()

	return nil
}
