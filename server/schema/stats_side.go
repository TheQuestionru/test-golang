package schema

import (
	"github.com/TheQuestionru/thequestion/server/types"
	"time"
)

// Dfp

type AnalyticsDfpRow struct {
	Id          int64      `json:"-" id:"true" generated:"true"`
	DfpNetwork  string     `json:"dfpNetwork" db:"dfp_network"`
	DfpId       string     `json:"dfpId" db:"dfp_id"`
	DfpName     string     `json:"dfpName" db:"dfp_name"`
	Timestamp   types.Date `json:"timestamp" db:"timestamp"`
	Impressions int32      `json:"impressions" db:"impressions"`
	Clicks      int32      `json:"clicks" db:"clicks"`
}

type AnalyticsDfpStatsQuery struct {
	DfpNetwork string              `json:"dfpNetwork" db:"dfp_network"`
	DfpId      types.NullString    `json:"dfpId" db:"dfp_id"`
	Start      time.Time           `json:"start" db:"start"`
	End        time.Time           `json:"end" db:"end"`
	GroupBy    types.GroupByPeriod `json:"groupBy" db:"-"`
}

// Ga

type AnalyticsGaRow struct {
	Id                  int64           `json:"-" id:"true" generated:"true"`
	GaId                string          `json:"gaId" db:"ga_id"`
	QuestionId          types.NullInt64 `json:"questionId" db:"question_id"`
	Timestamp           types.Date      `json:"timestamp" db:"timestamp"`
	Users               int32           `json:"users" db:"users"`
	Sessions            int32           `json:"sessions" db:"sessions"`
	Pageviews           int32           `json:"pageviews" db:"pageviews"`
	PageviewsPerSession float64         `json:"pageviewsPerSession" db:"pageviews_per_session"`
	AvgSessionDuration  int32           `json:"avgSessionDuration" db:"avg_session_duration"`
	BounceRate          float64         `json:"bounceRate" db:"bounce_rate"`
	PercentNewSessions  float64         `json:"percentNewSessions" db:"percent_new_sessions"`
}

func (v *AnalyticsGaRow) Add(a *AnalyticsGaRow) *AnalyticsGaRow {
	v.Users += a.Users
	v.Sessions += a.Sessions
	v.Pageviews += a.Pageviews
	v.PageviewsPerSession = (v.PageviewsPerSession + a.PageviewsPerSession) / 2
	v.AvgSessionDuration = (v.AvgSessionDuration + a.AvgSessionDuration) / 2
	v.BounceRate = (v.BounceRate + a.BounceRate) / 2
	v.PercentNewSessions = (v.PercentNewSessions + a.PercentNewSessions) / 2
	return v
}

type AnalyticsGaStatsQuery struct {
	GaId       string              `json:"gaId" db:"ga_id"`
	QuestionId types.NullInt64     `json:"questionId" db:"question_id"`
	Start      types.Date          `json:"start" db:"start"`
	End        types.Date          `json:"end" db:"end"`
	GroupBy    types.GroupByPeriod `json:"groupBy" db:"-"`
}
