package manager

import (
	"context"
	"fmt"
	"testing"

	"github.com/SUSE-Enceladus/csp-rancher-usage-operator/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

type testScenario struct {
	numRancherNodes     int
	result              testResult
}

type testResult struct {
	errResult           bool
}

func (s *testScenario) runScenario(t *testing.T) {
	mockK8sClient := mocks.NewMockK8sClient()
	mockScraper := mocks.NewMockScraper(s.numRancherNodes)
	mockUsageOperator := UsageOperator{
		k8s:     mockK8sClient,
		scraper: mockScraper,
	}

	nodeCount, err := mockUsageOperator.getNodeCount(context.TODO())
	// check that the results were as expected
	if s.result.errResult {
		assert.Error(t, err, fmt.Sprintf("Scenario: %v", s))
	} else {
		assert.NoError(t, err, fmt.Sprintf("Scenario: %v", s))
		assert.Equal(t, nodeCount, s.numRancherNodes, "expecting managed nodes to be %d, got %d instead", s.numRancherNodes, nodeCount)
	}
}

// TestCheckout tests basic scenarios where we don't have anything previously checked out
func TestCheckout(t *testing.T) {
	scenarios := []testScenario{
		{
			numRancherNodes:     20,
			result: testResult{
				errResult:           false,
			},
		},
	}
	for _, scenario := range scenarios {
		scenario.runScenario(t)
	}
}
