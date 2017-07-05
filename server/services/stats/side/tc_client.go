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
	GetBuildStats() ([]teamcity.Build, error)
}

type tcClient struct {
	CountGetTask int
	apiClient    teamcity.Client
}

func NewTcClient(config Config) TcClient {

	if config.TeamCityHost == "test" {
		return NewTestTcClient(config)
	}
	auth := teamcity.BasicAuth(config.TeamCityUser, config.TeamCityPass)
	client := teamcity.NewClient(config.TeamCityHost, auth)

	return &tcClient{
		CountGetTask: config.TeamCityCountGetBuilds,
		apiClient:    client,
	}
}

func (tc *tcClient) GetBuildStats() ([]teamcity.Build, error) {

	builds, err := tc.apiClient.GetBuilds(tc.CountGetTask)
	if err != nil {
		panic(err)
	}
	lenList := tc.CountGetTask
	contGetTask := len(builds)
	if lenList > contGetTask {
		lenList = contGetTask
	}

	taskList := make([]teamcity.Build, lenList)

	i := 0
	for _, task := range builds {
		taskList[i] = task
		i++
	}

	return taskList, nil
}
