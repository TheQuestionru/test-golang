package stats_side

import (
	"github.com/ivankorobkov/di"
	"github.com/yfronto/newrelic"
)

func NrClientModule(m *di.Module) {
	m.AddConstructor(NewNrClient)
	m.MarkPackageDep(Config{})
}

type NrClient interface {
	GetServersStats() ([]newrelic.Server, error)
}

type nrClient struct {
	apiClient *newrelic.Client
}

func NewNrClient(config Config) NrClient {
	if config.NewRelicApiKey == TestNewRelicApiKey {
		return NewTestNrClient()
	}

	return &nrClient{
		apiClient: newrelic.NewClient(config.NewRelicApiKey),
	}
}

func (t *nrClient) GetServersStats() ([]newrelic.Server, error) {
	opt := new(newrelic.ServersOptions)

	return t.apiClient.GetServers(opt)
}
