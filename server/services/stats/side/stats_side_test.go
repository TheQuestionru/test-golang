package stats_side

import (
	"os"
	"testing"

	"github.com/ivankorobkov/di"
	"github.com/stretchr/testify/assert"
)

var test struct {
	Stats *TestStats
}

func TestMain(m *testing.M) {
	di.MustFill(&test, TestModule)

	os.Exit(m.Run())
}

func TestSideStats_Realtime(t *testing.T) {
	realtime, err := test.Stats.Realtime()
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEqual(t, 0, realtime)
}

func TestSideStats_ServersStats(t *testing.T) {
	servers, err := test.Stats.ServersStats()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 10, len(servers))
}

func TestSideStats_TCBuildInfo(t *testing.T) {
	tcBuildInfo, err := test.Stats.TCBuildInfo()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []string{"SUCCESS", "FAILURE", "ERROR"}, tcBuildInfo.BuildStatuses)
	assert.Equal(t, "Alison", tcBuildInfo.LastChangesAuthor)
	assert.Equal(t, "20150714T121353+0000", tcBuildInfo.LastChangesDate)
}
