package stats_side

import (
	"github.com/ivankorobkov/di"
	"github.com/abourget/teamcity"
	"time"
)

func TeamCityClientModule(m *di.Module) {
	m.AddConstructor(NewTeamCityClient)
	m.MarkPackageDep(Config{})
}

type TeamCityClient interface {
	GetTeamCityView() ([]TaskView, error)
}

type teamCityClient struct {
	client *teamcity.Client
	count int
}

// 10 последних задач и их инфо: ветка, статус, время на сборку, ссылка на страницу с тимсити с подробной инфо
type TaskView struct {
	Branch  string
	BuildAt time.Time
	Status  string // check status type
	Url     string
}

func NewTeamCityClient(config Config) teamCityClient {
	if config.TeamCityUser == "test" {
		return NewTestTeamCityClient()
	}

	return teamCityClient{
		count : config.TeamCityViewTaskCount,
		client: teamcity.New(config.TeamCityUrl, config.TeamCityUser, config.TeamCityPassword),
	}
}

func (tc *teamCityClient) GetTeamCityView() ([]TaskView, error) {

	builds, err := tc.client.SearchBuild("all builds order by date created")
	if err != nil {
		return err
	}
	tcv := []TaskView{}
	for _, b := range builds {
		taskView := TaskView{
			Branch: b.BranchName,
			BuildAt: b.StartDate,
			Status: b.StatusText,
			Url: b.WebURL,
		}
		tcv = append(tcv, taskView)
		if len(tcv) == tc.count {
			break
		}
	}

	return tcv
}
