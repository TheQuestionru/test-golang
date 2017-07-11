package stats_side

import (
	"github.com/TheQuestionru/thequestion/server/schema"
	"github.com/ivankorobkov/di"
	"github.com/kapitanov/go-teamcity"
)

const (
	buildsSize int = 40
)

func TcClientModule(m *di.Module) {
	m.AddConstructor(NewTcClient)
	m.MarkPackageDep(Config{})
}

type TcClient interface {
	TcGetBuilds() ([]*schema.TcBuild, error)
}

type tcClient struct {
	teamcitySvc teamcity.Client
}

func NewTcClient(config Config) TcClient {
	if config.TcUrl == "http://localhost:8111" {
		return NewTestTcClient()
	}
	return &tcClient{
		teamcitySvc: teamcity.NewClient(config.TcUrl, teamcity.BasicAuth(config.TcUsername, config.TcPassword)),
	}
}

func (t *tcClient) TcGetBuilds() ([]*schema.TcBuild, error) {
	builds, err := t.teamcitySvc.GetBuilds(buildsSize)
	if err != nil {
		return nil, err
	}
	return t.parseBuilds(builds), nil
}

func (t *tcClient) parseBuilds(builds []teamcity.Build) []*schema.TcBuild {
	var result []*schema.TcBuild
	for _, build := range builds {
		result = append(result, &schema.TcBuild{
			ID:     build.ID,
			Status: build.StatusText,
		})
	}
	return result
}
