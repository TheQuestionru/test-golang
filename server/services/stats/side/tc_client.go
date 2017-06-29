package stats_side

import (
	"errors"
	"github.com/GSokol/go-teamcity"
	"github.com/ivankorobkov/di"
)

func TcClientModule(m *di.Module) {
	m.AddConstructor(NewTcClient)
	m.MarkPackageDep(Config{})
}

type TCBuildInfo struct {
	BuildStatuses     []string `json:"lastBuilds"`
	LastChangesAuthor string   `json:"lastChangesAuthor"`
	LastChangesDate   string   `json:"lastChangesDate"`
}

type TcClient interface {
	GetTCBuildInfo() (*TCBuildInfo, error)
}

type tcClient struct {
	client teamcity.Client
	config Config
}

const (
	tcBuildNumber = 5
)

var ErrTcNoChangesYet = errors.New("No changes yet")

func setTcConfigDefaults(config *Config) {
	if config.TeamCityBuildNumber == 0 {
		config.TeamCityBuildNumber = tcBuildNumber
	}
}

func NewTcClient(config Config) TcClient {
	var authorizer teamcity.Authorizer
	if config.TeamCityAddress == TestTeamCityAddress {
		return NewTestTcClient()
	}

	if config.TeamCityLogin == "" && config.TeamCityPassword == "" {
		authorizer = teamcity.GuestAuth()
	} else {
		if config.TeamCityLogin == "" || config.TeamCityPassword == "" {
			panic("TC Login or Password is not specified")
		}

		authorizer = teamcity.BasicAuth(config.TeamCityLogin, config.TeamCityPassword)
	}

	setTcConfigDefaults(&config)

	return &tcClient{
		client: teamcity.NewClient(
			config.TeamCityAddress,
			authorizer,
		),
		config: config,
	}
}

func (t *tcClient) setupBuilds(result *TCBuildInfo) error {
	builds, err := t.client.GetBuilds(t.config.TeamCityBuildNumber)
	if err != nil {
		return err
	}

	result.BuildStatuses = make([]string, len(builds), len(builds))
	for i, build := range builds {
		result.BuildStatuses[i] = build.StatusText
	}

	return nil
}

func (t *tcClient) setupChanges(result *TCBuildInfo) error {
	changes, err := t.client.GetChanges(1)
	if err != nil {
		return err
	}

	if len(changes) == 0 {
		return ErrTcNoChangesYet
	}

	result.LastChangesAuthor = changes[0].Username
	result.LastChangesDate = changes[0].Date

	return nil
}

func (t *tcClient) GetTCBuildInfo() (*TCBuildInfo, error) {
	tcBuildInfo := TCBuildInfo{}
	err := t.setupBuilds(&tcBuildInfo)
	if err != nil {
		return nil, err
	}

	err = t.setupChanges(&tcBuildInfo)
	if err != nil {
		return nil, err
	}

	return &tcBuildInfo, nil
}
