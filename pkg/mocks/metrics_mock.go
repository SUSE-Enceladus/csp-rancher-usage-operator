package mocks

import (
	"errors"

	"github.com/SUSE-Enceladus/csp-rancher-usage-operator/pkg/metrics"
)

type MockScraper struct {
	Nodes int
}

func NewMockScraper(numNodes int) *MockScraper {
	return &MockScraper{
		Nodes: numNodes,
	}
}

func (m *MockScraper) ScrapeAndParse() (*metrics.NodeCounts, error) {
	if m.Nodes >= 0 {
		return &metrics.NodeCounts{
			Total: m.Nodes,
		}, nil
	}
	// Trigger error
	return &metrics.NodeCounts{
		Total: m.Nodes,
	}, errors.New("")
}
