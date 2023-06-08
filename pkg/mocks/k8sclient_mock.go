package mocks

import (
	"encoding/json"
	"time"
)

/*import (
	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)*/

type MockK8sClient struct {
	CurrentSupportConfig       []byte
	CurrentManagedNodeCount    uint32
	CurrentNotificationMessage string
	RancherHostName            string
	RancherVersion             string
}

func NewMockK8sClient() *MockK8sClient {
	return &MockK8sClient{
		CurrentSupportConfig: nil,
	}
}

func (m *MockK8sClient) UpdateCSPConfigOutput(marshalledData []byte) error {
	//todo: mock error
	m.CurrentSupportConfig = marshalledData
	return nil
}

func (m *MockK8sClient) UpdateUserNotification(isInCompliance bool, message string) error {
	if !isInCompliance {
		m.CurrentNotificationMessage = message
	}
	return nil
}

func (m *MockK8sClient) GetRancherHostname() (string, error) {
	return m.RancherHostName, nil
}

func (m *MockK8sClient) GetRancherVersion() (string, error) {
	return m.RancherVersion, nil
}

func (m *MockK8sClient) UpdateProductUsage(managedNodes uint32) error {
	m.CurrentManagedNodeCount = managedNodes
	return nil
}

func (m *MockK8sClient) GetCSPConfigData() (string, error) {
	data := map[string]interface{}{
		"timestamp":             time.Now().Format(time.RFC3339),
		"billing_api_access_ok": true,
		"expire":                time.Now().AddDate(0, 1, 0).Format(time.RFC3339),
		"errors":                []string{""},
		"last_billed":           time.Now().AddDate(0, 0, -25).Format(time.RFC3339),
		"usage": map[string]interface{}{
			"foo": 10,
		},
		"customer_csp_data": map[string]interface{}{
			"document": map[string]interface{}{
				"foo": "bar",
			},
			"signature":      "signature",
			"pkcs7":          "pkcs7",
			"cloud_provider": "aws",
		},
		"base_product": "cpe:/o:suse:rancher:v2.7.3",
	}

	// Convert the JSON object to a string
	jEncData, err := json.Marshal(data)
	return string(jEncData), err

}
