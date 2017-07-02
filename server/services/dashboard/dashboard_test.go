package dashboard

import (
	"github.com/TheQuestionru/thequestion/server/schema"
	"github.com/TheQuestionru/thequestion/server/types"
	"github.com/ivankorobkov/di"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var test struct {
	Dashboard *TestDashboard
}

func TestMain(m *testing.M) {
	di.MustFill(&test, TestModule)

	os.Exit(m.Run())
}

func TestDashboard_UpdateElements__should_update_elements(t *testing.T) {
	dashboard, err := test.Dashboard.GetDashboard(types.NewEmptyReq())
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(dashboard.Rows))
	assert.Equal(t, 2, len(dashboard.Rows[0].Elements))
	assert.Equal(t, schema.DashboardElementTypeNRServers, dashboard.Rows[0].Elements[0].Type)
	assert.Equal(t, schema.DashboardElementTypeTeamCity, dashboard.Rows[0].Elements[1].Type)
	assert.Equal(t, 10, len(dashboard.Rows[0].Elements[0].Servers))
}
