package services

import (
	"github.com/TheQuestionru/thequestion/server/services/dashboard"
	"github.com/ivankorobkov/di"
)

func Module(m *di.Module) {
	m.Import(dashboard.Module)
	m.Import(stats_side.Module)
	m.Import(stats_side.GaClientModule)
	m.Import(stats_side.NrClientModule)
	m.Import(stats_side.TcClientModule)
	m.AddConstructor(New)
}

type Services struct {
	SideStats stats_side.SideStats
	Dashboard dashboard.Dashboard
}

func New(
	sideStats stats_side.SideStats,
	dashboard dashboard.Dashboard,

) Services {

	return Services{
		SideStats: sideStats,
		Dashboard: dashboard,
	}
}
