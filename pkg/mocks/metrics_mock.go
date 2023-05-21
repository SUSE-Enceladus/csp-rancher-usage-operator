package mocks

import (
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
	// TODO: Error case
	return &metrics.NodeCounts{
		Total: m.Nodes,
	}, nil
}
