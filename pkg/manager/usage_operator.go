package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/SUSE-Enceladus/csp-rancher-usage-operator/pkg/clients/k8s"
	"github.com/SUSE-Enceladus/csp-rancher-usage-operator/pkg/metrics"
	"github.com/sirupsen/logrus"
)

type UsageOperator struct {
	k8s     k8s.Client
	scraper metrics.Scraper
}

const (
        cspBillingNoBillThreshold   = "K8S_CSP_BILLING_NO_BILL_THRESHOLD"
        cspBillingManagerInterval   = "K8S_CSP_BILLING_MANAGER_INTERVAL"
)

var (
        noBillThreshold  int
	managerInterval  time.Duration
)


func NewUsageOperator(k k8s.Client, s metrics.Scraper) (*UsageOperator, error) {
	err := readConstantsFromEnv()
        if err != nil {
                return nil, err
        }

	return &UsageOperator{
		k8s:     k,
		scraper: s,
	}, nil
}

func readConstantsFromEnv() error {
        var noBillThresholdStr string = os.Getenv(cspBillingNoBillThreshold)
        var managerIntervalStr string = os.Getenv(cspBillingManagerInterval)
        var missingEnvVars []string
	var err error
        if noBillThresholdStr == "" {
                missingEnvVars = append(missingEnvVars, cspBillingNoBillThreshold)
        } else {
		noBillThreshold, err = strconv.Atoi(noBillThresholdStr)
                if err != nil {
                        missingEnvVars = append(missingEnvVars, cspBillingNoBillThreshold)
                }
        }
        if managerIntervalStr == "" {
                missingEnvVars = append(missingEnvVars, cspBillingManagerInterval)
        } else {
		managerInterval, err = time.ParseDuration(managerIntervalStr)
                if err != nil {
                        missingEnvVars = append(missingEnvVars, cspBillingManagerInterval)
                }
        }

        if len(missingEnvVars) == 0 {
                return nil
        }
        return fmt.Errorf("unable to read required env vars %v", missingEnvVars)
}

func (m *UsageOperator) Start(ctx context.Context, errs chan<- error) {
	go m.start(ctx, errs)
}

func (m *UsageOperator) start(ctx context.Context, errs chan<- error) {
	for range ticker(ctx, managerInterval) {
		logrus.Info("Calculating managed nodes")
		ncount, err := m.getNodeCount(ctx)
                if err != nil {
                       errs <- err
		       logrus.Warningf("Unable to get managed node count. Will not update usage record.")

                } else {
                	logrus.Infof("Node Count %d", ncount)
			err = m.k8s.UpdateProductUsage(ncount)
			if err != nil {
				errs <- err
				logrus.Infof("Failed to update product usage.")
			}
                }
		// Check to see if we need to warn the user
                // about unsupported Rancher cluster if we
                // are unable to gather usage data pass x number
                // of days.
		err = m.checkAndUpdateUserNotifications()
                if err != nil {
                	errs <- err
                }
	}
	logrus.Infof("[manager] exiting")
}

func (m *UsageOperator) checkAndUpdateUserNotifications() error {
	mashaledData, err := m.k8s.GetCSPConfigData()
	if err != nil {
		return err
	}
	var config CSPConfig
        err = json.Unmarshal([]byte(mashaledData), &config)
	if err != nil {
		logrus.Infof("Failed to unmarshal csp_config: %v", err)
                return err
        }
	logrus.Debugf("csp-config: %#v", config)
	// TODO(gyee): should we just display the last error or all the errors
	// in the errors field?
	var message_prefix string = "One or more errors detected by CSP billing. This Rancher instance may not be supported till the problem is resolved. Please resolve the following errors or contact SUSE for further assistance. "
	var cspErrors string = ""
	// Notify user if unable to access billing API
        if !config.BillingAPIAccessOK {
		cspErrors = cspErrors + "Unable to access billing API. "
        }
	// Notify user if expire date is in the past
	expireDate, err := time.Parse(time.RFC3339, config.Expire)
	if err != nil {
		return err
	}
	if expireDate.Before(time.Now()) {
		cspErrors = cspErrors + "Unable to meter usage. "
        }
        // Notify user if they haven't been billed
	// FIXME: if user hasn't been billed, what does the lastBilled date is
	// expected to be?
	if config.LastBilled != "" {
		lastBilledDate, err := time.Parse(time.RFC3339, config.LastBilled)
		if err != nil {
                	return err
        	}
		if lastBilledDate.AddDate(0, 0, noBillThreshold).Before(time.Now()) {
			cspErrors = "It's been more than " + strconv.Itoa(noBillThreshold) + " days since last billed. "
		}
	}
	if len(cspErrors) > 0 {
                if len(config.Errors) > 0 {
                        cspErrors = cspErrors + "Errors: " + strings.Join(config.Errors, " ")
                }
                cspErrors = message_prefix + cspErrors
		return m.k8s.UpdateUserNotification(false, cspErrors)
        } else {
		// clear the user notification since there are no errors
		return m.k8s.UpdateUserNotification(true, "")
	}
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
