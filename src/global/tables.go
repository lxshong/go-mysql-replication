package global

import (
	"fmt"
	"github.com/go-mysql-org/go-mysql/client"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/schema"
)

type Tables map[string]map[string]*schema.Table

// 规则匹配
func (t *Tables) Match(database string, table string) bool {
	_, ok := (*t)[database][table]
	if !ok {
		return false
	}
	return true
}

// 获取表信息
func (t *Tables) GetTable(database string, table string) *schema.Table {
	tb, ok := (*t)[database][table]
	if !ok {
		return nil
	}
	return tb
}

func (t *Tables) Change(rules2 rules) error {
	cnf := Cfg()

	conn,err := client.Connect(fmt.Sprintf("%s:%d",cnf.Host,cnf.Port),cnf.User,cnf.Password,"information_schema")
	if err != nil {
		return err
	}
	defer conn.Close()
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
				(*t)[tableSchema] = map[string]*schema.Table{}
			}
			conn1,err := client.Connect(fmt.Sprintf("%s:%d",cnf.Host,cnf.Port),cnf.User,cnf.Password,tableSchema)
			if err != nil {
				return err
			}
			tb,err := schema.NewTable(conn1,tableSchema,tableName)
			if err != nil {
				return err
			}
			conn1.Close()
			(*t)[tableSchema][tableName] = tb
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
