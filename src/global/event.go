package global

type Event struct {
	Schema string
	Table string
	Action string
	Rows []map[string]interface{}
}

func NewEvent(schema string,table string,action string,rows []map[string]interface{}) *Event {
	return &Event{Schema: schema,Table: table,Action: action,Rows: rows}
}