package stats_side

import (
	"github.com/ivankorobkov/di"
	"github.com/kapitanov/go-teamcity"
)

func TcClientModule(m *di.Module) {
	m.AddConstructor(NewTcClient)
	m.MarkPackageDep(Config{})
}

type TcClient interface {
	GetBuilds() ([]teamcity.Build, error)
}

type tcClient struct {
	apiClient teamcity.Client
}

func NewTcClient(config Config) TcClient {
	if config.TeamCityHost == "test" {
		return NewTestTcClient()
	}

	return &tcClient{
		apiClient: teamcity.NewClient(config.TeamCityHost, teamcity.GuestAuth()),
	}
}

func (t *tcClient) GetBuilds() ([]teamcity.Build, error) {
	return t.apiClient.GetBuilds(5)
}
