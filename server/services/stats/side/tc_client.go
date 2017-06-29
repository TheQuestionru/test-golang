package stats_side

import (
	"github.com/kapitanov/go-teamcity"
	"github.com/ivankorobkov/di"
)

func TCModule(m *di.Module) {
	m.AddConstructor(NewTcClient)
	m.MarkPackageDep(Config{})
}

type TcClient interface {
	GetProjectsStatus() ([]teamcity.Project, error)
}

type tcClient struct {
	apiClient teamcity.Client
}

func NewTcClient(config Config) TcClient {
	if config.NewRelicApiKey == "test" {
		return NewTestTcClient()
	}

	return &tcClient{
		apiClient: teamcity.NewClient(config.TeamCityAddress, teamcity.GuestAuth()),
	}
}


func (t *tcClient) GetProjectsStatus() ([]teamcity.Project, error) {
	return t.apiClient.GetProjects()
}