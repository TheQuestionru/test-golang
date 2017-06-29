package dashboard

import (
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

	assert.Equal(t, 10, len(dashboard.Rows[0].Elements[0].Servers))
	assert.Equal(t, []string{"SUCCESS", "FAILURE", "ERROR"}, dashboard.Rows[0].Elements[1].BuildStatuses)
	assert.Equal(t, "Alison", dashboard.Rows[0].Elements[1].LastChangesAuthor.String)
	assert.Equal(t, "20150714T121353+0000", dashboard.Rows[0].Elements[1].LastChangesDate.String)
}
