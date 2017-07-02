package stats_side

import (
	"github.com/abourget/teamcity"
	"github.com/ivankorobkov/di"
)

func TcClientModule(m *di.Module) {
	m.AddConstructor(NewTcClient)
	m.MarkPackageDep(Config{})
}

type TcClient interface {
	GetBuildStats() ([]*teamcity.Build, error)
}

type tcClient struct {
	CountGetTask int
	apiClient    *teamcity.Client
}

func NewTcClient(config Config) TcClient {

	return &tcClient{
		CountGetTask: config.TeamCityCountGetTask,
		apiClient:    teamcity.New(config.TeamCityHost, config.TeamCityUser, config.TeamCityPass),
	}
}

func (tc *tcClient) GetBuildStats() ([]*teamcity.Build, error) {

	builds, err := tc.apiClient.SearchBuild("all")
	if err != nil {
		panic(err)
	}
	lenList := tc.CountGetTask
	contGetTask := len(builds)
	if lenList > contGetTask {
		lenList = contGetTask
	}

	taskList := make([]*teamcity.Build, lenList, lenList)

	i := 0
	for _, task := range builds {
		taskList[i] = task
		i++
	}

	return taskList, nil
}
