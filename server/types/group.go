package types

import (
	"strings"
)

type GroupByPeriod string

const (
	GroupByPeriodInvalid GroupByPeriod = ""
	GroupByPeriodDay     GroupByPeriod = "day"
	GroupByPeriodWeek    GroupByPeriod = "week"
	GroupByPeriodMonth   GroupByPeriod = "month"
)

var groupByMap = map[GroupByPeriod]bool{
	GroupByPeriodDay:   true,
	GroupByPeriodWeek:  true,
	GroupByPeriodMonth: true,
}

func (b GroupByPeriod) Clean() GroupByPeriod {
	b = GroupByPeriod(strings.ToLower(CleanString(string(b))))
	if len(b) == 0 {
		return GroupByPeriodInvalid
	}
	if _, ok := groupByMap[b]; !ok {
		return GroupByPeriodInvalid
	}
	return b
}
