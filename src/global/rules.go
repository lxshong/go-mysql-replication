package global

import "regexp"

type rules map[string][]string

// 规则匹配
func (r rules)Match(database string, table string) bool {
	tablesPattern,ok := r[database]
	if !ok {
		return false
	}
	for _, tablePattern := range tablesPattern {
		if match,_ := regexp.Match(tablePattern,[]byte(table));match{
			return true
		}
	}
	return false
}
