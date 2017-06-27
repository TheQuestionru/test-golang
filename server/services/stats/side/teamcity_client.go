package stats_side

import (
	"github.com/ivankorobkov/di"
	"github.com/kapitanov/go-teamcity"
)

type (
	TcClient interface {
		GetProjects() ([]teamcity.Project, error)
	}

	tcClient struct {
		client teamcity.Client
	}
)

func TcClientModule(m *di.Module) {
	m.AddConstructor(NewTcClient)
	m.MarkPackageDep(Config{})
}

func NewTcClient(config Config) TcClient {
	if config.TeamCityAddress == "test" {
		return NewTestTcClient()
	}

	return &tcClient{
		client: teamcity.NewClient(
			config.TeamCityAddress,
			teamcity.GuestAuth(),
		),
	}
}

func (t *tcClient) GetProjects() ([]teamcity.Project, error) {
	return t.client.GetProjects()
}
