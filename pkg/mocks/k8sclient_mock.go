package mocks

import (
	"encoding/json"
	"errors"
	"time"
)

const timeFormat = time.RFC3339

type MockK8sClient struct {
	CurrentManagedNodeCount    int
	CurrentNotificationMessage string
	RancherMetricsAPIEndpoint  string
	RancherVersion             string

	Error Error
}

type Error struct {
	Trigger   bool
	Condition string
}

func NewMockK8sClient(Error Error) *MockK8sClient {
	if Error.Trigger {
		return &MockK8sClient{Error: Error}
	}
	return &MockK8sClient{}
}

func (m *MockK8sClient) UpdateUserNotification(clearError bool, message string) error {
	if !clearError {
		m.CurrentNotificationMessage = message
	}
	return nil
}

func (m *MockK8sClient) GetRancherMetricsAPIEndpoint() (string, error) {
	if m.Error.Trigger && m.Error.Condition == "ErrorRancherMetricsAPIEndpoint" {
		return "", errors.New("trigger mock error for GetRancherMetricsAPIEndpoint")
	}
	return m.RancherMetricsAPIEndpoint, nil
}

func (m *MockK8sClient) GetRancherVersion() (string, error) {
	if m.Error.Trigger && m.Error.Condition == "ErrorRancherVersion" {
		return "", errors.New("trigger mock error for GetRancherVersion")
	}
	return m.RancherVersion, nil
}

func (m *MockK8sClient) UpdateProductUsage(managedNodes int) error {
	//TODO: mock errorneous conditions
	m.CurrentManagedNodeCount = managedNodes
	return nil
}

func (m *MockK8sClient) GetCSPConfigData() (string, error) {
	var data map[string]interface{}
	timestamp := time.Now().Format(timeFormat)
	expire := time.Now().AddDate(0, 1, 0).Format(timeFormat)
	lastBilled := time.Now().AddDate(0, 0, -25).Format(timeFormat)

	// valid data
	data = map[string]interface{}{
		"timestamp":             timestamp,
		"billing_api_access_ok": true,
		"expire":                expire,
		"errors":                []string{""},
		"last_billed":           lastBilled,
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

	// mocking various error conditions
	if m.Error.Trigger && m.Error.Condition == "ErrorBillingNotExists" {
		// Errorneous Condition - removing billing_api_access
		delete(data, "billing_api_access_ok")
	}
	if m.Error.Trigger && m.Error.Condition == "ErrorBillingFalse" {
		// Errorneous Condition - setting billing_api_access to false
		data["billing_api_access_ok"] = false
	}
	if m.Error.Trigger && m.Error.Condition == "ErrorExpiry" {
		// Errorneous Condition - expire date is in the past
		data["expire"] = time.Now().AddDate(0, 0, -2).Format(timeFormat)
	}
	if m.Error.Trigger && m.Error.Condition == "ErrorLastBilledEmpty" {
		// Errorneous Condition - lastBilled empty
		data["last_billed"] = ""
	}
	if m.Error.Trigger && m.Error.Condition == "ErrorLastBilledThreshold" {
		// Errorneous Condition - lastBilled beyond threshold
		data["last_billed"] = time.Now().AddDate(0, 0, -47).Format(timeFormat)
	}
	if m.Error.Trigger && m.Error.Condition == "ErrorMarshalUnsupportedType" {
		// Errorneous Condition - marshal error by adding invalid field type
		data["email"] = make(chan string)
	}
	if m.Error.Trigger && m.Error.Condition == "ErrorMarshaldataMalformed" {
		// Errorneous Condition - marshal error for malformed
		data["billing_api_access_ok"] = "STRINGISTEADOFBOOL"
	}
	if m.Error.Trigger && m.Error.Condition == "ErrorLastBilledNonRFC3339" {
		// Errorneous Condition - lastBilled in non-RFC3339
		data["last_billed"] = time.Now().AddDate(0, 0, -25).Format(time.Stamp)
	}
	// Convert the JSON object to a string
	jEncData, err := json.Marshal(data)
	return string(jEncData), err
}
