package mocks

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
