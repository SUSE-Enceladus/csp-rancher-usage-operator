// Package k8s contains kubernetes related functionality for the csp adapter (such as cache storage and results output)
package k8s

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rancher/lasso/pkg/controller"
	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"github.com/rancher/wrangler/pkg/clients"
	v1 "github.com/rancher/wrangler/pkg/generated/controllers/core/v1"
	usagerecordsv1 "github.com/SUSE-Enceladus/csp-rancher-usage-operator/api/usagerecords/v1"
	productusage "github.com/SUSE-Enceladus/csp-rancher-usage-operator/generated/clientset/versioned"
	productusagev1 "github.com/SUSE-Enceladus/csp-rancher-usage-operator/generated/clientset/versioned/typed/usagerecords/v1"
	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

const (
	cspBillingAdapterNamespace  = "K8S_CSP_BILLING_ADAPTER_NAMESPACE"
	cspBillingAdapterConfigMap  = "K8S_CSP_BILLING_ADAPTER_CONFIGMAP"
	cspBillingNoBillThreshold   = "K8S_CSP_BILLING_NO_BILL_THRESHOLD"
	cspNotification             = "K8S_OUTPUT_NOTIFICATION"
	hostnameSettingEnv          = "K8S_HOSTNAME_SETTING"
	versionSettingEnv           = "K8S_RANCHER_VERSION_SETTING"
	cspConfigKey                = "data"
	// TODO(gyee): need to update Rancher UI code to look for csp-usage-operator here
	// https://github.com/rancher/dashboard/blob/d112f1781d8efd96daeacbf23fe0f341e3a0b28a/shell/components/AwsComplianceBanner.vue#L20
	cspComponentName            = "csp-adapter"
	baseProductPrefix           = "cpe:/o:suse:rancher:"
)

var (
	cspBillingNamespace        string
	cspBillingConfigMapName	   string
	noBillThreshold            uint32
	outputNotificationName     string
	hostnameSetting            string
	versionSetting             string
)

type Client interface {
	// UpdateCSPConfigOutput stores config to k8s as a configmap with a static/constant name
	UpdateCSPConfigOutput(marshalledData []byte) error
	// UpdateUserNotification creates/updates a RancherUserNotification based on isInCompliance and the provided message
	UpdateUserNotification(clearError bool, message string) error
	// GetRancherHostname finds the hostname for the core rancher install from the settings.
	GetRancherHostname() (string, error)
	// GetRancherVersion finds the version of rancher from the settings
	GetRancherVersion() (string, error)
	// UpdateProductUsage updates the RancherUsageRecord with the current managed node count
	UpdateProductUsage(managedNodes uint32) error
	// get CSP config data
	GetCSPConfigData() (string, error)
}

type Clients struct {
	ProductUsage  productusagev1.ProductUsageInterface
	ConfigMaps    v1.ConfigMapClient
	Secrets       v1.SecretController
	Notifications controller.SharedController
	Settings      controller.SharedController
}

func New(ctx context.Context, rest *rest.Config) (*Clients, error) {
	err := readConstantsFromEnv()
	if err != nil {
		return nil, err
	}
	clients, err := clients.NewFromConfig(rest, nil)
	if err != nil {
		return nil, err
	}

	productUsageClient, err := productusage.NewForConfig(rest)
	if err != nil {
                return nil, err
        }

	if err := clients.Start(ctx); err != nil {
		return nil, err
	}

	localSchemeBuilder := runtime.SchemeBuilder{
		v3.AddToScheme,
	}
	scheme := runtime.NewScheme()
	err = localSchemeBuilder.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}
	factory, err := controller.NewSharedControllerFactoryFromConfig(rest, scheme)
	if err != nil {
		return nil, err
	}
	settingGVR := schema.GroupVersionResource{Group: "management.cattle.io", Version: "v3", Resource: "settings"}
	settingKind := "Setting"
	settingController := factory.ForResourceKind(settingGVR, settingKind, false)
	err = settingController.Start(context.Background(), 1)
	if err != nil {
		return nil, fmt.Errorf("error when starting setting controller %w", err)
	}

	notificationGVR := schema.GroupVersionResource{Group: "management.cattle.io", Version: "v3", Resource: "rancherusernotifications"}
	notificationKind := "RancherUserNotification"
	notificationController := factory.ForResourceKind(notificationGVR, notificationKind, false)
	err = notificationController.Start(context.Background(), 1)
	if err != nil {
		return nil, fmt.Errorf("error when starting notification controller %w", err)
	}

	return &Clients{
		ProductUsage:  productUsageClient.UsagerecordsV1().ProductUsages(),
		ConfigMaps:    clients.Core.ConfigMap(),
		Notifications: notificationController,
		Settings:      settingController,
	}, nil
}

// readConstantsFromEnv sets the outputConfigMapName, outputNotificationName, cacheName, and hostnameSetting after
// reading values from the env - returns an error if one or more values were not found. Values for these are defined
// in _helpers.tpl
func readConstantsFromEnv() error {
	outputNotificationName = os.Getenv(cspNotification)
	cspBillingNamespace = os.Getenv(cspBillingAdapterNamespace)
	cspBillingConfigMapName = os.Getenv(cspBillingAdapterConfigMap)
	var noBillThresholdStr string = os.Getenv(cspBillingNoBillThreshold)
	hostnameSetting = os.Getenv(hostnameSettingEnv)
	versionSetting = os.Getenv(versionSettingEnv)
	var missingEnvVars []string
	if cspBillingNamespace == "" {
		missingEnvVars = append(missingEnvVars, cspBillingAdapterNamespace)
        } 
	if cspBillingConfigMapName == "" {
		missingEnvVars = append(missingEnvVars, cspBillingAdapterConfigMap)
	}
        if noBillThresholdStr == "" {
		missingEnvVars = append(missingEnvVars, cspBillingNoBillThreshold)
	} else {
		noBillThresholdInt, err := strconv.Atoi(noBillThresholdStr)
		if err != nil {
			missingEnvVars = append(missingEnvVars, cspBillingNoBillThreshold)
		} else {
			noBillThreshold = uint32(noBillThresholdInt)
		}
	}

	if outputNotificationName == "" {
		missingEnvVars = append(missingEnvVars, cspNotification)
	}
	if hostnameSetting == "" {
		missingEnvVars = append(missingEnvVars, hostnameSettingEnv)
	}
	if versionSetting == "" {
		missingEnvVars = append(missingEnvVars, versionSettingEnv)
	}
	if len(missingEnvVars) == 0 {
		return nil
	}
	return fmt.Errorf("unable to read required env vars %v", missingEnvVars)
}

func (c *Clients) UpdateCSPConfigOutput(marshalledData []byte) error {
	// since the data from this output is nested, we have to stick this all under one key in raw format
	data := map[string]string{
		cspConfigKey: string(marshalledData),
	}
	currentConfigMap, err := c.ConfigMaps.Get(cspBillingAdapterNamespace, cspBillingConfigMapName, metav1.GetOptions{})
	if apierror.IsNotFound(err) {
		_, err = c.ConfigMaps.Create(&corev1.ConfigMap{
			Data: data,
			ObjectMeta: metav1.ObjectMeta{
				Name:      cspBillingConfigMapName,
				Namespace: cspBillingAdapterNamespace,
			},
		})
		return err
	}
	currentConfigMap = currentConfigMap.DeepCopy()
	currentConfigMap.Data = data
	_, err = c.ConfigMaps.Update(currentConfigMap)
	return err
}

func (c *Clients) UpdateUserNotification(clearError bool, message string) error {
	if clearError {
		// Clear user notifications
		err := c.Notifications.Client().Delete(context.TODO(), "", outputNotificationName, metav1.DeleteOptions{})
		if err != nil && !apierror.IsNotFound(err) {
			// ignore not found errors - this means we didn't have a notification to delete, so we didn't need to adjust
			return err
		}
	} else {
		current := &v3.RancherUserNotification{}
		err := c.Notifications.Client().Get(context.TODO(), "", outputNotificationName, current, metav1.GetOptions{})
		if err != nil {
			if apierror.IsNotFound(err) {
				// not found means we need to make a new notification
				current = &v3.RancherUserNotification{
					ObjectMeta: metav1.ObjectMeta{
						Name: outputNotificationName,
					},
					ComponentName: cspComponentName,
					Message:       message,
				}
				err = c.Notifications.Client().Create(context.TODO(), "", current, current, metav1.CreateOptions{})
			}
			return err
		}
		// update all relevant fields - also updating component name to future-proof against changes made to this field
		current = current.DeepCopy()
		current.Message = message
		current.ComponentName = cspComponentName
		err = c.Notifications.Client().Update(context.TODO(), "", current, current, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Clients) GetRancherHostname() (string, error) {
	setting := &v3.Setting{}
	err := c.Settings.Client().Get(context.TODO(), "", hostnameSetting, setting, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	// server-url includes the protocol prefix - we need the actual hostname to be returned
	hostname := strings.TrimPrefix(setting.Value, "https://")
	return hostname, nil
}

func (c *Clients) GetRancherVersion() (string, error) {
	setting := &v3.Setting{}
	err := c.Settings.Client().Get(context.TODO(), "", versionSetting, setting, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return setting.Value, nil
}

func (c *Clients) UpdateProductUsage(managedNodes uint32) error {
        currentUsage, err := c.ProductUsage.Get(context.TODO(), "rancher-usage-record", metav1.GetOptions{})
	reportingTime := time.Now().Format(time.RFC3339)
        if apierror.IsNotFound(err) {
		rancherVersion, err := c.GetRancherVersion()
		if err != nil {
                	return err
        	}
                _, err = c.ProductUsage.Create(context.TODO(), &usagerecordsv1.ProductUsage{
                        ObjectMeta: metav1.ObjectMeta{
                                Name: "rancher-usage-record",
                        },
			// TODO(gyee): check with PM about base product. Use "Rancher ${VERISON}" for now
			BaseProduct: baseProductPrefix + rancherVersion,
			ManagedNodeCount: managedNodes,
			ReportingTime: reportingTime,
                }, metav1.CreateOptions{})
                return err
        }
        currentUsage = currentUsage.DeepCopy()
        currentUsage.ManagedNodeCount = managedNodes
	currentUsage.ReportingTime = reportingTime
        _, err = c.ProductUsage.Update(context.TODO(), currentUsage, metav1.UpdateOptions{})
        return err
}

func (c *Clients) GetCSPConfigData() (string, error) {
	currentConfigMap, err := c.ConfigMaps.Get(cspBillingNamespace, cspBillingConfigMapName, metav1.GetOptions{})
	if err != nil {
        	return "", err
        }
	return currentConfigMap.Data[cspConfigKey], err
}
