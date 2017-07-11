package stats_side

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/TheQuestionru/thequestion/server/lib/logger"
	"github.com/TheQuestionru/thequestion/server/schema"
	"github.com/TheQuestionru/thequestion/server/types"
	"github.com/ivankorobkov/di"
	"github.com/yfronto/newrelic"
)

func TestModule(m *di.Module) {
	m.Import(Module)
	m.Import(logger.TestModule)
	m.AddConstructor(NewTestConfig)
	m.AddConstructor(NewTestGaClient)
	m.AddConstructor(NewNrClient)
	m.AddConstructor(NewTcClient)
	m.AddConstructor(NewTest)
}

type TestStats struct {
	SideStats
	gaClient GaClient
	nrClient NrClient
	tcClient TcClient
}

func NewTest(s SideStats, gaClient GaClient, nrClient NrClient, tcClient TcClient) *TestStats {
	ret := &TestStats{s, gaClient, nrClient, tcClient}
	return ret
}

func NewTestConfig() Config {
	return Config{
		GoogleServiceKeyFile: "test",
		Enabled:              true,
		GoogleAnalyticsIds:   map[string]string{"TheQuestion": "ga:91655992"},
		NewRelicApiKey:       "test",
		TcUrl:                "http://localhost:8111",
		TcUsername:           "root",
		TcPassword:           "064695",
	}
}

type testTcClient struct {
}

func NewTestTcClient() TcClient {
	return &testTcClient{}
}

func (t *testTcClient) TcGetBuilds() ([]*schema.TcBuild, error) {
	return []*schema.TcBuild{
		{
			ID:     1,
			Status: "",
		},
		{
			ID:     2,
			Status: "",
		},
	}, nil
}

type testNrClient struct {
	count int
}

func NewTestNrClient() NrClient {
	return &testNrClient{}
}

func (t *testNrClient) GetServersStats() ([]newrelic.Server, error) {
	servers := []newrelic.Server{}

	for i := 0; i < 10; i++ {
		count := t.count
		t.count++

		servers = append(servers, newrelic.Server{
			Name: fmt.Sprintf("test server name %v", count),
		})
	}
	return servers, nil
}

type testGaClient struct {
	summary   *schema.AnalyticsGaRow
	questions map[int64]*schema.AnalyticsGaRow
	timestamp time.Time
	ids       map[string]bool
}

func NewTestGaClient(config Config) GaClient {
	ids := map[string]bool{}
	for _, key := range config.GoogleAnalyticsIds {
		ids[key] = true
	}
	return &testGaClient{
		ids: ids,
	}
}

func (c *testGaClient) GaGetSummaryData(id string, from time.Time, to time.Time) (*schema.AnalyticsGaRow, error) {
	c.timestamp = to.Truncate(Day)
	c.summary.Timestamp = types.NewDate(c.timestamp)
	return c.summary, c.check(id, from, to)
}

func (c *testGaClient) GaGetQuestionsData(id string, from time.Time, to time.Time) (map[int64]*schema.AnalyticsGaRow, error) {
	c.timestamp = to.Truncate(Day)
	for _, qq := range c.questions {
		qq.Timestamp = types.NewDate(c.timestamp)
	}
	return c.questions, c.check(id, from, to)
}

func (c *testGaClient) GaGetRealtime(id string) (int64, error) {
	return rand.Int63n(int64(10000)), nil
}

func (c *testGaClient) check(id string, from time.Time, to time.Time) error {
	if !c.ids[id] {
		return errors.New("Id mismatch")
	}
	return nil
}

func (c *testGaClient) setSummary(row *schema.AnalyticsGaRow) {
	c.summary = row
}

func (c *testGaClient) setQuestions(rows map[int64]*schema.AnalyticsGaRow) {
	c.questions = rows
}

func (s *TestStats) TestSummary(t *testing.T, gaId string) *schema.AnalyticsGaRow {
	summary := &schema.AnalyticsGaRow{
		GaId:                gaId,
		Users:               1,
		Sessions:            2,
		Pageviews:           3,
		PageviewsPerSession: 2,
		AvgSessionDuration:  4,
		BounceRate:          6,
		PercentNewSessions:  8,
	}
	s.gaClient.(*testGaClient).setSummary(summary)
	return summary
}

var ErrTestDfpNetworkMissing = errors.New("Network incorrect")
var ErrTestDfpIncorrectJobId = errors.New("Job id incorrect")
var ErrTestDfpIncorrectJobStatus = errors.New("Job not ready")
var ErrTestDfpIncorrectUrl = errors.New("Incorrect url")
