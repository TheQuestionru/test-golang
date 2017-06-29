package stats_side

import (
	"strconv"
	"time"

	"github.com/TheQuestionru/thequestion/server/lib/logger"
	"github.com/TheQuestionru/thequestion/server/schema"
	"github.com/TheQuestionru/thequestion/server/types"
	"github.com/ivankorobkov/di"
	"google.golang.org/api/analytics/v3"
)

func GaClientModule(m *di.Module) {
	m.Import(logger.Module)
	m.AddConstructor(NewGaClient)
	m.MarkPackageDep(Config{})
}

type GaClient interface {
	GaGetSummaryData(id string, from time.Time, to time.Time) (*schema.AnalyticsGaRow, error)
	GaGetQuestionsData(id string, from time.Time, to time.Time) (map[int64]*schema.AnalyticsGaRow, error)
	GaGetRealtime(id string) (int64, error)
}

var ErrBadApiAnswer = types.NewError("erorrs.BadApiAnswer") // "Bad api answer"

const (
	gaMetricsList = "ga:users,ga:sessions,ga:pageviews,ga:pageviewsPerSession," +
		"ga:avgSessionDuration,ga:bounceRate,ga:percentNewSessions"
	gGaMetricsSize = 7
	gaDateTpl      = "2006-01-02"
	gaMaxInt       = int64(0x7FFFFFFF)

	rtRealtimeMetric = "rt:activeUsers"
)

type gaClient struct {
	logger       logger.Logger
	analyticsSvc *analytics.Service
	realtimeSvc  *analytics.DataRealtimeService
	config       Config
}

func NewGaClient(logger logger.Logger, config Config) GaClient {
	if config.GoogleServiceKeyFile == TestGoogleServiceKeyFile {
		return NewTestGaClient(config)
	}

	client, err := newJwtClient(config.GoogleServiceKeyFile, analytics.AnalyticsReadonlyScope)
	if err != nil {
		panic(err)
	}

	analyticsSvc, err := analytics.New(client)
	if err != nil {
		panic(err)
	}

	return &gaClient{
		logger:       logger,
		analyticsSvc: analyticsSvc,
		realtimeSvc:  analytics.NewDataRealtimeService(analyticsSvc),
		config:       config,
	}
}

func (t *gaClient) GaGetRealtime(id string) (int64, error) {
	call := t.realtimeSvc.Get(id, rtRealtimeMetric)

	data, err := call.Do()
	if err != nil {
		return int64(0), err
	}

	return strconv.ParseInt(data.Rows[0][0], 10, 64)
}

func (s *gaClient) GaGetSummaryData(id string, from time.Time, to time.Time) (*schema.AnalyticsGaRow, error) {
	res, err := s.analyticsSvc.Data.Ga.
		Get(id, from.Format(gaDateTpl), to.Format(gaDateTpl), gaMetricsList).Do()
	if err != nil {
		return nil, err
	}

	if len(res.Rows) != 1 {
		return nil, ErrBadApiAnswer
	}

	row, err := parseAnalyticsRow(res.Rows[0])
	if err != nil {
		return nil, err
	}

	row.GaId = id
	return row, nil
}

func (s *gaClient) GaGetQuestionsData(id string, from time.Time, to time.Time) (map[int64]*schema.AnalyticsGaRow, error) {
	offset := int64(1)
	var total int64
	var rows = [][]string{}
	for {
		res, err := s.analyticsSvc.Data.Ga.
			Get(id, from.Format(gaDateTpl), to.Format(gaDateTpl), gaMetricsList).
			Dimensions("ga:pagePath").
			MaxResults(gaMaxInt).StartIndex(offset).Filters("ga:pagePath=~^/questions/").Do()
		if err != nil {
			return nil, err
		}

		total = res.TotalResults

		if len(res.Rows) == 0 {
			break
		}
		offset += int64(len(res.Rows))
		rows = append(rows, res.Rows...)
	}

	if int(total) != len(rows) {
		return nil, ErrBadApiAnswer
	}

	ret := map[int64]*schema.AnalyticsGaRow{}
	for _, rowX := range rows {
		if len(rowX) != gGaMetricsSize+1 {
			return nil, ErrBadApiAnswer
		}

		match := questionUri.FindStringSubmatch(rowX[0])
		if match == nil {
			continue
		}

		questionId, _ := strconv.ParseInt(match[1], 10, 64)
		row, err := parseAnalyticsRow(rowX[1:])
		if err != nil {
			return nil, err
		}

		row.GaId = id
		row.QuestionId = types.NewNullInt64(questionId)
		if v, ok := ret[questionId]; ok {
			v.Add(row)
		} else {
			ret[questionId] = row
		}
	}
	return ret, nil
}

func parseAnalyticsRow(row []string) (*schema.AnalyticsGaRow, error) {
	if len(row) != gGaMetricsSize {
		return nil, ErrBadApiAnswer
	}

	users, err := strconv.ParseInt(row[0], 10, 32)
	if err != nil || users > gaMaxInt {
		return nil, ErrBadApiAnswer
	}

	sessions, err := strconv.ParseInt(row[1], 10, 32)
	if err != nil || sessions > gaMaxInt {
		return nil, ErrBadApiAnswer
	}

	pageviews, err := strconv.ParseInt(row[2], 10, 32)
	if err != nil || pageviews > gaMaxInt {
		return nil, ErrBadApiAnswer
	}

	pageviewsPerSession, err := strconv.ParseFloat(row[3], 64)
	if err != nil {
		return nil, ErrBadApiAnswer
	}

	avgSessionDuration, err := strconv.ParseFloat(row[4], 64)
	if err != nil {
		return nil, ErrBadApiAnswer
	}

	bounceRate, err := strconv.ParseFloat(row[5], 64)
	if err != nil {
		return nil, ErrBadApiAnswer
	}

	percentNewSessions, err := strconv.ParseFloat(row[6], 64)
	if err != nil {
		return nil, ErrBadApiAnswer
	}

	return &schema.AnalyticsGaRow{
		Users:               int32(users),
		Sessions:            int32(sessions),
		Pageviews:           int32(pageviews),
		PageviewsPerSession: pageviewsPerSession,
		AvgSessionDuration:  int32(avgSessionDuration),
		BounceRate:          bounceRate,
		PercentNewSessions:  percentNewSessions,
	}, nil
}
