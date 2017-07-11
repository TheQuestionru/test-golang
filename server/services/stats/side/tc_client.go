package stats_side

import (
	"encoding/json"
	"fmt"
	"github.com/TheQuestionru/thequestion/server/schema"
	"github.com/ivankorobkov/di"
	"net/http"
)

func TcClientModule(m *di.Module) {
	m.AddConstructor(NewTcClient)
	m.MarkPackageDep(Config{})
}

type (
	TcClient interface {
		GetAgents() ([]schema.Agent, error)
	}

	tcClient struct {
		config Config
		client http.Client
	}
)

func NewTcClient(config Config) TcClient {
	if config.TeamcityEndpoint == "test" && config.TeamcityAuthHeader == "test" {
		return NewTestTcClient()
	}

	return &tcClient{
		config: config,
		client: http.Client{},
	}
}

func (t *tcClient) GetAgents() ([]schema.Agent, error) {
	agents := &schema.AgentsResponse{}
	err := t.load(
		"GET",
		"agents",
		agents,
	)
	if err != nil {
		return nil, err
	}
	return agents.Agents, nil
}

func (t *tcClient) load(method, url string, response interface{}) error {
	req, err := http.NewRequest(
		method,
		fmt.Sprintf(
			"%s/httpAuth/app/rest/%s",
			t.config.TeamcityEndpoint,
			url,
		),
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Add(
		"Authorization",
		t.config.TeamcityAuthHeader,
	)

	req.Header.Add(
		"Accept",
		"application/json",
	)

	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(response)
}
