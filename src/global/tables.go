package global

import (
	"fmt"
	"github.com/go-mysql-org/go-mysql/client"
	"github.com/go-mysql-org/go-mysql/mysql"
)

type Tables map[string]map[string]bool

// 规则匹配
func (t *Tables) Match(database string, table string) bool {
	_, ok := (*t)[database][table]
	if !ok {
		return false
	}
	return true
}

func (t *Tables) Add(database string, table string) error {
	(*t)[database][table] = true
	return nil
}

func (t *Tables) Change(rules2 rules) error {
	cnf := Cfg()
	conn,err := client.Connect(fmt.Sprintf("%s:%d",cnf.Host,cnf.Port),cnf.User,cnf.Password,"information_schema")
	if err != nil {
		return err
	}
	if err := conn.Ping();err != nil {
		return err
	}
	var tableSchema string
	var tableName string
	var result mysql.Result
	err = conn.ExecuteSelectStreaming(`SELECT table_schema,table_name FROM information_schema.Tables`, &result, func(row []mysql.FieldValue) error {
		for idx, val := range row {
			field := string(result.Fields[idx].Name)
			val := string(val.Value().([]byte))
			if field == "table_schema" {
				tableSchema = val
			}else {
				tableName = val
			}
		}
		if rules2.Match(tableSchema, tableName) {
			if _, ok := (*t)[tableSchema];!ok {
				(*t)[tableSchema] = map[string]bool{}
			}
			(*t)[tableSchema][tableName] = true
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}


