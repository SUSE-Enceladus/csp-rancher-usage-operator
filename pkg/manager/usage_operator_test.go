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

func (s *testScenario) runGetNodeCountScenario(t *testing.T) {

	t.Setenv("K8S_CSP_BILLING_NO_BILL_THRESHOLD", "45")
	mockK8sClient := mocks.NewMockK8sClient()

	if s.result.errResult {
		// Validate the node count for error usecase
		// Create an instance of the mock object
		mockScraper := mocks.NewMockScraper(s.numRancherNodes)
		mockUsageOperator, e := NewUsageOperator(mockK8sClient, mockScraper)
		assert.NoError(t, e)

		nodeCount, err := mockUsageOperator.getNodeCount(context.TODO())
		assert.Error(t, err, fmt.Sprintf("Scenario: %v", s))
		assert.Contains(t, err.Error(), "unable to determine number of active nodes")
		assert.Equal(t, nodeCount, 0) // count is returned as 0 when err != nil

	} else {
		// Validate the node count for success usecase
		mockScraper := mocks.NewMockScraper(s.numRancherNodes)
		mockUsageOperator, e := NewUsageOperator(mockK8sClient, mockScraper)
		assert.NoError(t, e)

		nodeCount, err := mockUsageOperator.getNodeCount(context.TODO())
		assert.NoError(t, err, fmt.Sprintf("Scenario: %v", s))
		assert.Equal(t, nodeCount, s.numRancherNodes, "expecting managed nodes to be %d, got %d instead", s.numRancherNodes, nodeCount)
	}
}

// TestGetNodeCount tests basic scenarios
func TestGetNodeCount(t *testing.T) {
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
		scenario.runGetNodeCountScenario(t)
	}
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
	mockK8sClient := mocks.NewMockK8sClient()
	hook := test.NewGlobal()
	hook.Reset()

	t.Setenv("K8S_CSP_BILLING_NO_BILL_THRESHOLD", "45")

	if s.result.errResult {
		// Validate channel for specific error "unable to determine number of active nodes"
		// Validate log for specific message "Unable to get managed node count. Will not update usage record."
		mockScraper := mocks.NewMockScraper(s.numRancherNodes)
		mockUsageOperator, e := NewUsageOperator(mockK8sClient, mockScraper)
		assert.NoError(t, e)

		ctxWithCancel, cancelFunction := context.WithCancel(context.TODO())
		testChan := make(chan error, 1)

		mockUsageOperator.Start(ctxWithCancel, testChan)
		time.Sleep(35 * time.Second)
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

		//Derive a context with cancel
		ctxWithCancel, cancelFunction := context.WithCancel(context.TODO())

		// create channels
		testChan := make(chan error, 1)

		mockUsageOperator.Start(ctxWithCancel, testChan)
		time.Sleep(35 * time.Second)

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
		// k8sclient handles count as uint32, scraper handles count as int
		assert.Equal(t, mockK8sClient.CurrentManagedNodeCount, uint32(s.numRancherNodes))

		// Validate that no errors in usernotification
		assert.Equal(t, mockK8sClient.CurrentNotificationMessage, "")

		hook.Reset()
		assert.Nil(t, hook.LastEntry())

		close(testChan)
		cancelFunction()
	}
}

func TestNoEnvFlag(t *testing.T) {
	mockK8sClient := mocks.NewMockK8sClient()
	mockScraper := mocks.NewMockScraper(5)
	mockUsageOperator, e := NewUsageOperator(mockK8sClient, mockScraper)
	assert.Error(t, e)
	assert.Contains(t, e.Error(), "unable to read required env vars")
	assert.Nil(t, mockUsageOperator)
}
