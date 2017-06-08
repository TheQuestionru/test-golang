package stats_side

import (
	"github.com/TheQuestionru/thequestion/server/lib/logger"
	"github.com/TheQuestionru/thequestion/server/types"
	"github.com/ivankorobkov/di"
	"github.com/yfronto/newrelic"
	"regexp"
	"time"
)

var questionUri = regexp.MustCompile(`^/questions/(\d+)`)

func Module(m *di.Module) {
	m.Import(logger.Module)
	m.AddConstructor(New)
	m.MarkPackageDeps(struct {
		Config
		GaClient
		NrClient
	}{})
}

type Config struct {
	GoogleServiceKeyFile string            `yaml:"GoogleServiceKeyFile"`
	GoogleAnalyticsIds   map[string]string `yaml:"GoogleAnalyticsIds"`
	GoogleDfpNetworkIds  map[string]string `yaml:"GoogleDfpNetworkIds"`
	Enabled              bool              `yaml:"Enabled"`
	Schedule             string            `yaml:"Schedule"`
	NewRelicApiKey       string            `yaml:"NewRelicApiKey"`
	TeamcityAuthHeader   string            `yaml:"TeamcityAuthHeader"`
	TeamcityEndpoint     string            `yaml:"TeamcityEndpoint"`
}

type SideStats interface {
	RunUpdateGa()

	Realtime() (int64, error)
	ServersStats() ([]newrelic.Server, error)
}

type sideStats struct {
	logger   logger.Logger
	gaClient GaClient
	nrClient NrClient
	config   Config
}

const (
	MaxRetry    = 5
	Day         = time.Hour * 24
	DfpPoolTime = time.Second * 2
	DfpPoolMax  = 10
	rtGAKeyName = "TheQuestion"
)

func New(logger logger.Logger, config Config, gaClient GaClient,
	nrClient NrClient) SideStats {
	return &sideStats{
		logger:   logger.Prefix("side-stats"),
		gaClient: gaClient,
		nrClient: nrClient,
		config:   config,
	}
}

func (s *sideStats) RunUpdateGa() {
	date := time.Now().Truncate(Day)
	retry := 0
	for {
		err := s.tryRunUpdateGa(date)
		if err == nil {
			break
		}
		retry++
		s.logger.Error("Ga stats failed", logger.Payload{"error": err, "retry": retry})
		if retry > MaxRetry {
			break
		}
	}
}

func (t *sideStats) Realtime() (int64, error) {
	id := t.config.GoogleAnalyticsIds[rtGAKeyName]
	return t.gaClient.GaGetRealtime(id)
}

func (t *sideStats) ServersStats() ([]newrelic.Server, error) {
	return t.nrClient.GetServersStats()
}

func (s *sideStats) tryRunUpdateGa(timestamp time.Time) error {
	for _, gaId := range s.config.GoogleAnalyticsIds {
		s.logger.Info("Updating", logger.Payload{"gaId": gaId})

		found := false

		if found {
			s.logger.Info("Update early, skipping", logger.Payload{"gaId": gaId})
			continue
		}

		summary, err := s.gaClient.GaGetSummaryData(gaId, timestamp.Add(-Day), timestamp)
		if err != nil {
			return err
		}
		summary.Timestamp = types.NewDate(timestamp)
		s.logger.Info("Summary received", logger.Payload{"gaId": gaId})

		questions, err := s.gaClient.GaGetQuestionsData(gaId, timestamp.Add(-Day), timestamp)
		if err != nil {
			return err
		}
		s.logger.Info("Questions received", logger.Payload{"gaId": gaId, "count": len(questions)})

		qIds := make([]int64, len(questions))
		j := 0
		for id, qq := range questions {
			qIds[j] = id
			qq.Timestamp = types.NewDate(timestamp)
			j++
		}

		// saving ga stats to db

		s.logger.Info("Ga stats saved", logger.Payload{"gaId": gaId})
	}

	return nil
}
