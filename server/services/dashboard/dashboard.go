package dashboard

import (
	"github.com/TheQuestionru/thequestion/server/lib/logger"
	"github.com/TheQuestionru/thequestion/server/schema"
	"github.com/TheQuestionru/thequestion/server/services/stats/side"
	"github.com/TheQuestionru/thequestion/server/types"
	"github.com/ivankorobkov/di"
	"sort"
)

func Module(m *di.Module) {
	m.Import(logger.Module)
	m.Import(stats_side.Module)
	m.AddConstructor(New)
}

type Dashboard interface {
	GetDashboard(req types.Req) (*schema.DashboardView, error)
}

type dashboard struct {
	logger logger.Logger

	sideStats stats_side.SideStats
}

func New(logger logger.Logger, sideStats stats_side.SideStats) Dashboard {
	return &dashboard{
		logger:    logger.Prefix("dashboard"),
		sideStats: sideStats,
	}
}

func (t *dashboard) GetDashboard(req types.Req) (*schema.DashboardView, error) {
	grid, err := t.getDashboardGrid()
	if err != nil {
		return nil, err
	}

	view := &schema.DashboardView{
		Rows: []*schema.DashboardRowView{},
	}
	for _, row := range grid {
		rowView := &schema.DashboardRowView{
			Elements: []*schema.DashboardElementView{},
		}
		for _, element := range row {
			rowView.Elements = append(rowView.Elements, element)
		}

		view.Rows = append(view.Rows, rowView)
	}

	return view, nil
}

func (t *dashboard) getDashboardGrid() ([][]*schema.DashboardElementView, error) {
	elements := []*schema.DashboardElement{} // simulate db query

	elements = append(elements, &schema.DashboardElement{
		RowNumber: 0,
		ColNumber: 0,
		DashboardElementKey: schema.DashboardElementKey{
			Type: schema.DashboardElementTypeNRServers,
		},
	})

	grid := t.makeGrid(elements)
	view := [][]*schema.DashboardElementView{}
	for _, row := range grid {
		rowView := []*schema.DashboardElementView{}
		for _, element := range row {
			elementView := &schema.DashboardElementView{}
			elementView.DashboardElement = element

			switch elementView.Type {
			case schema.DashboardElementTypeGARealtime:
				realtime, err := t.sideStats.Realtime()
				if err != nil {
					return nil, err
				}

				elementView.Realtime = types.NewNullInt64(realtime)
			case schema.DashboardElementTypeNRServers:
				var err error
				elementView.Servers, err = t.sideStats.ServersStats()
				if err != nil {
					return nil, err
				}
			}

			rowView = append(rowView, elementView)
		}

		view = append(view, rowView)
	}

	return view, nil
}

func (t *dashboard) makeGrid(elements []*schema.DashboardElement) [][]*schema.DashboardElement {
	elementsByRow := map[int32][]*schema.DashboardElement{}
	for _, element := range elements {
		_, ok := elementsByRow[element.RowNumber]
		if ok {
			elementsByRow[element.RowNumber] = append(elementsByRow[element.RowNumber], element)
			continue
		}

		elementsByRow[element.RowNumber] = []*schema.DashboardElement{element}
	}

	grid := [][]*schema.DashboardElement{}
	for _, row := range elementsByRow {
		sort.Slice(row, func(i, j int) bool {
			return row[i].ColNumber < row[j].ColNumber
		})

		grid = append(grid, row)
	}

	sort.Slice(grid, func(i, j int) bool {
		return grid[i][0].RowNumber < grid[j][0].RowNumber
	})

	return grid
}
