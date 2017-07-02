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

func TestSideStats_BuildStats(t *testing.T) {
	tasks, err := test.Stats.BuildStats()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 5, len(tasks))
}
