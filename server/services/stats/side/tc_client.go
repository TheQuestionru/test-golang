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
	GetProjectList() ([]teamcity.Project, error)
}

type tcClient struct {
	apiClient teamcity.Client
}

func NewTcClient(config Config) TcClient {
	if config.TeamCityAddr == "test" {
		return NewTestTcClient()
	}

	return &tcClient{
		apiClient: teamcity.NewClient(config.TeamCityAddr, teamcity.GuestAuth()),
	}
}

func (t tcClient) GetProjectList() ([]teamcity.Project, error) {
	return t.apiClient.GetProjects()
}
