package manager

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/SUSE-Enceladus/csp-rancher-usage-operator/pkg/mocks"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

type testScenario struct {
	numRancherNodes int
	result          testResult
}

type testResult struct {
	errResult bool
}

func TestStart(t *testing.T) {
	scenarios := []testScenario{
		{
			numRancherNodes: 20,
			result: testResult{
				errResult: false,
			},
		},
		{
			numRancherNodes: -1,
			result: testResult{
				errResult: true,
			},
		},
		{
			numRancherNodes: 0,
			result: testResult{
				errResult: false,
			},
		},
	}
	for _, scenario := range scenarios {
		scenario.runStartScenario(t)
	}

}

func (s *testScenario) runStartScenario(t *testing.T) {
	merror := mocks.Error{
		Trigger:   false,
		Condition: "",
	}

	mockK8sClient := mocks.NewMockK8sClient(merror)

	hook := test.NewGlobal()
	hook.Reset()

	t.Setenv("K8S_CSP_BILLING_NO_BILL_THRESHOLD", "45")
	t.Setenv("K8S_CSP_BILLING_MANAGER_INTERVAL", "2s")

	if s.result.errResult {
		// Validate channel for specific error "unable to determine number of active nodes"
		// Validate log for specific message "Unable to get managed node count. Will not update usage record."
		mockScraper := mocks.NewMockScraper(s.numRancherNodes)
		mockUsageOperator, e := NewUsageOperator(mockK8sClient, mockScraper)
		assert.NoError(t, e)

		nodeCount, e := mockUsageOperator.getNodeCount(context.TODO())
		assert.Error(t, e, fmt.Sprintf("Scenario: %v", s))
		assert.Contains(t, e.Error(), "unable to determine number of active nodes")
		assert.Equal(t, nodeCount, 0) // count is returned as 0 when err != nil

		ctxWithCancel, cancelFunction := context.WithCancel(context.TODO())
		testChan := make(chan error, 1)

		mockUsageOperator.Start(ctxWithCancel, testChan)
		time.Sleep(3 * time.Second)
		err := <-testChan
		assert.Contains(t, err.Error(), "unable to determine number of active nodes")

		// Validate specific log messages
		var joinedString, sep string
		for _, entry := range hook.Entries {
			joinedString += sep + entry.Message
			sep = " "
		}
		assert.Contains(t, joinedString, "Calculating managed nodes")
		assert.Contains(t, joinedString, "Unable to get managed node count. Will not update usage record.")

		hook.Reset()
		assert.Nil(t, hook.LastEntry())

		close(testChan)
		cancelFunction()

	} else {
		// Validate channel for no errors
		// Validate log for specific message "Node Count x"
		mockScraper := mocks.NewMockScraper(s.numRancherNodes)
		mockUsageOperator, e := NewUsageOperator(mockK8sClient, mockScraper)
		assert.NoError(t, e)

		nodeCount, e := mockUsageOperator.getNodeCount(context.TODO())
		assert.NoError(t, e, fmt.Sprintf("Scenario: %v", s))
		assert.Equal(t, nodeCount, s.numRancherNodes, "expecting managed nodes to be %d, got %d instead", s.numRancherNodes, nodeCount)

		//Derive a context with cancel
		ctxWithCancel, cancelFunction := context.WithCancel(context.TODO())

		// create channels
		testChan := make(chan error, 1)

		mockUsageOperator.Start(ctxWithCancel, testChan)
		time.Sleep(3 * time.Second)

		assert.Equal(t, len(testChan), 0)

		// Validate specific log messages
		var joinedString, sep string
		for _, entry := range hook.Entries {
			joinedString += sep + entry.Message
			sep = " "
		}
		msgExpected := "Node Count " + strconv.Itoa(s.numRancherNodes)
		assert.Contains(t, joinedString, "Calculating managed nodes")
		assert.Contains(t, joinedString, msgExpected)

		// Validate Product Usage Updated
		assert.Equal(t, mockK8sClient.CurrentManagedNodeCount, s.numRancherNodes)

		// Validate that no errors in usernotification
		assert.Equal(t, mockK8sClient.CurrentNotificationMessage, "")

		hook.Reset()
		assert.Nil(t, hook.LastEntry())

		close(testChan)
		cancelFunction()
	}
}

func TestNoEnvFlag(t *testing.T) {
	// No environment flags
	merror := mocks.Error{
		Trigger:   false,
		Condition: "",
	}

	mockK8sClient := mocks.NewMockK8sClient(merror)
	mockScraper := mocks.NewMockScraper(5)
	mockUsageOperator, e := NewUsageOperator(mockK8sClient, mockScraper)
	assert.Error(t, e)
	assert.Contains(t, e.Error(), "unable to read required env vars")
	assert.Nil(t, mockUsageOperator)
}

type testScenarioInvalidCSPConfigUserNotification struct {
	numRancherNodes         int
	expectedNotificationMsg string
	Error                   mocks.Error
}

func TestInvalidCSPConfigUserNotification(t *testing.T) {
	scenarios := []testScenarioInvalidCSPConfigUserNotification{
		{
			numRancherNodes:         1,
			expectedNotificationMsg: "Unable to access billing API.",
			Error: mocks.Error{
				Trigger:   true,
				Condition: "ErrorBillingNotExists",
			},
		},
		{
			numRancherNodes:         2,
			expectedNotificationMsg: "Unable to access billing API.",
			Error: mocks.Error{
				Trigger:   true,
				Condition: "ErrorBillingFalse",
			},
		},
		{
			numRancherNodes:         3,
			expectedNotificationMsg: "Unable to meter usage.",
			Error: mocks.Error{
				Trigger:   true,
				Condition: "ErrorExpiry",
			},
		},
		//{
		// Behavavior for this use case needs to be determined
		//numRancherNodes:         4,
		//expectedNotificationMsg: "????????",
		//Error: mocks.Error{
		//	Trigger:   true,
		//	Condition: "ErrorLastBilledEmpty",
		//},
		//},
		{
			numRancherNodes:         5,
			expectedNotificationMsg: "days since last billed.",
			Error: mocks.Error{
				Trigger:   true,
				Condition: "ErrorLastBilledThreshold",
			},
		},
	}
	for _, scenario := range scenarios {
		scenario.CheckUserNotifications(t)
	}

}

func (s *testScenarioInvalidCSPConfigUserNotification) CheckUserNotifications(t *testing.T) {
	t.Setenv("K8S_CSP_BILLING_NO_BILL_THRESHOLD", "45")
	t.Setenv("K8S_CSP_BILLING_MANAGER_INTERVAL", "2s")

	mockK8sClient := mocks.NewMockK8sClient(s.Error)
	mockScraper := mocks.NewMockScraper(s.numRancherNodes)
	mockUsageOperator, e := NewUsageOperator(mockK8sClient, mockScraper)

	assert.NoError(t, e)

	//Derive a context with cancel
	ctxWithCancel, cancelFunction := context.WithCancel(context.TODO())

	// create channels
	testChan := make(chan error, 1)

	mockUsageOperator.Start(ctxWithCancel, testChan)
	time.Sleep(3 * time.Second)

	// Validate for specific errors in usernotification
	assert.Contains(t, mockK8sClient.CurrentNotificationMessage, s.expectedNotificationMsg)

	close(testChan)
	cancelFunction()

}

type testScenarioCSPConfigError struct {
	numRancherNodes int
	expectedError   string
	Error           mocks.Error
}

func TestCSPConfigError(t *testing.T) {
	// no user notifications
	scenarios := []testScenarioCSPConfigError{
		{
			numRancherNodes: 1,
			expectedError:   "json: unsupported type:",
			Error: mocks.Error{
				Trigger:   true,
				Condition: "ErrorMarshalUnsupportedType",
			},
		},
		{
			numRancherNodes: 2,
			expectedError:   "cannot parse",
			Error: mocks.Error{
				Trigger:   true,
				Condition: "ErrorLastBilledNonRFC3339",
			},
		},
	}
	for _, scenario := range scenarios {
		scenario.CheckError(t)
	}
}

func (s *testScenarioCSPConfigError) CheckError(t *testing.T) {
	t.Setenv("K8S_CSP_BILLING_NO_BILL_THRESHOLD", "45")
	t.Setenv("K8S_CSP_BILLING_MANAGER_INTERVAL", "2s")
	mockK8sClient := mocks.NewMockK8sClient(s.Error)
	mockScraper := mocks.NewMockScraper(s.numRancherNodes)
	mockUsageOperator, e := NewUsageOperator(mockK8sClient, mockScraper)

	assert.NoError(t, e)

	//Derive a context with cancel
	ctxWithCancel, cancelFunction := context.WithCancel(context.TODO())

	// create channels
	testChan := make(chan error, 1)

	mockUsageOperator.Start(ctxWithCancel, testChan)
	time.Sleep(3 * time.Second)

	// validate error
	err := <-testChan

	expectedSubstring := s.expectedError
	assert.Contains(t, err.Error(), expectedSubstring, "Expected error message to contain substring: "+expectedSubstring)

	close(testChan)
	cancelFunction()

}

func TestInvalidCSPConfigMarshalDataMalformed(t *testing.T) {
	t.Setenv("K8S_CSP_BILLING_NO_BILL_THRESHOLD", "45")
	t.Setenv("K8S_CSP_BILLING_MANAGER_INTERVAL", "2s")
	merror := mocks.Error{
		Trigger:   true,
		Condition: "ErrorMarshaldataMalformed",
	}
	hook := test.NewGlobal()
	hook.Reset()
	mockK8sClient := mocks.NewMockK8sClient(merror)
	mockScraper := mocks.NewMockScraper(1)
	mockUsageOperator, e := NewUsageOperator(mockK8sClient, mockScraper)
	assert.NoError(t, e)

	//Derive a context with cancel
	ctxWithCancel, cancelFunction := context.WithCancel(context.TODO())

	// create channels
	testChan := make(chan error, 1)

	mockUsageOperator.Start(ctxWithCancel, testChan)
	time.Sleep(3 * time.Second)

	// Validate specific log messages
	var joinedString, sep string
	for _, entry := range hook.Entries {
		joinedString += sep + entry.Message
		sep = " "
	}
	assert.Contains(t, joinedString, "Failed to unmarshal csp_config")
	close(testChan)
	cancelFunction()

}
