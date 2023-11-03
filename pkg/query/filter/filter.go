package filter

import (
	"fmt"
	"strings"
)

type Filter struct {
	Field     string
	Condition string
	Value     string
	Value2    string
	ValueInt  bool
}

// GetDBFilterMapCondition is a function for getting DB filter map condition.
func GetDBFilterMapCondition(filters []Filter) string {
	if filters == nil {
		return ""
	}

	var conditionBuilder strings.Builder
	count := len(filters)

	for i, filter := range filters {
		if i == 0 || i == count-1 {
			conditionBuilder.WriteString("(")
		} else {
			conditionBuilder.WriteString(" AND (")
		}

		if filter.Value2 != "" {
			conditionBuilder.WriteString(fmt.Sprintf("(%s %s %s AND %s)", filter.Field, filter.Condition, filter.Value, filter.Value2))
		} else {
			value := filter.Value
			if !filter.ValueInt {
				value = fmt.Sprintf("'%s'", filter.Value)
			}
			conditionBuilder.WriteString(fmt.Sprintf("%s %s %s", filter.Field, filter.Condition, value))
		}

		conditionBuilder.WriteString(")")
	}

	return conditionBuilder.String()
}
