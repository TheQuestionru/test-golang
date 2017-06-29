package dashboard

import (
	"github.com/TheQuestionru/test-golang/server/lib/logger"
	"github.com/TheQuestionru/test-golang/server/services/stats/side"
	"github.com/ivankorobkov/di"
)

func TestModule(m *di.Module) {
	m.Import(Module)
	m.Import(logger.TestModule)
	m.Import(stats_side.TestModule)
	m.AddConstructor(NewTestDashboard)
}

type TestDashboard struct {
	Dashboard

	counter int
}

func NewTestDashboard(dashboard Dashboard) *TestDashboard {
	return &TestDashboard{dashboard, 0}
}
