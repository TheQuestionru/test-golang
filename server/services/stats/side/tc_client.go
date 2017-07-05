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
	CountGetBuilds int
	apiClient      teamcity.Client
}

func NewTcClient(config Config) TcClient {

	if config.TeamCityHost == "test" {
		return NewTestTcClient(config)
	}
	auth := teamcity.BasicAuth(config.TeamCityUser, config.TeamCityPass)
	client := teamcity.NewClient(config.TeamCityHost, auth)

	return &tcClient{
		CountGetBuilds: config.TeamCityCountGetBuilds,
		apiClient:      client,
	}
}

func (tc *tcClient) GetBuildStats() ([]teamcity.Build, error) {

	builds, err := tc.apiClient.GetBuilds(tc.CountGetBuilds)
	if err != nil {
		panic(err)
	}
	lenList := tc.CountGetBuilds
	contGetBuilds := len(builds)
	if lenList > contGetBuilds {
		lenList = contGetBuilds
	}

	BuildsList := make([]teamcity.Build, lenList)

	i := 0
	for _, build := range builds {
		BuildsList[i] = build
		i++
	}

	return BuildsList, nil
}
