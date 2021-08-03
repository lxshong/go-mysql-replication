package service

import (
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/go-mysql-org/go-mysql/schema"
	"go-mysql-replication/src/endpoint"
	"go-mysql-replication/src/global"
	"log"
)

const (
	DEFAULT_CAP  = 4096
	UpdateAction = "update"
	InsertAction = "insert"
	DeleteAction = "delete"
)

type handler struct {
	cap      int
	queue    chan *replication.BinlogEvent
	stop     chan struct{}
	tables   global.Tables
	endpoint endpoint.Endpoint
}

func DefaultCap() int {
	return DEFAULT_CAP
}

func NewHandle(cap int) *handler {
	return &handler{
		cap: cap,
	}
}

func (h *handler) initialize() error {
	h.queue = make(chan *replication.BinlogEvent, h.cap)
	h.stop = make(chan struct{}, 1)
	h.tables = global.Cfg().Tables
	ep, err := endpoint.GetEndpoint()
	if err != nil {
		return err
	} else {
		h.endpoint = ep
	}
	err = h.startListener()
	if err != nil {
		return err
	}
	return nil
}

func (h *handler) Callback(event *replication.BinlogEvent) error {
	h.queue <- event
	return nil
}

func (h *handler) startListener() error {
	log.Println("start listener .. .. ..")
	go func() {
		var ev *replication.BinlogEvent = nil
		for {
			select {
			case ev = <-h.queue:
				err := h.deal(ev)
				if err != nil {
					return
				}
			case <-h.stop:
				return
			}
		}
	}()
	return nil
}

func (h *handler) stopListener() {
	log.Println("stop listener")
	h.stop <- struct{}{}
}

func (h *handler) deal(e *replication.BinlogEvent) error {

	switch e.Event.(type) {
	case *replication.RowsEvent:
		return h.RowsEventDeal(e)
	}

	return nil
}

func (h *handler) RowsEventDeal(e *replication.BinlogEvent) error {
	ev := e.Event.(*replication.RowsEvent)

	tableName := string(ev.Table.Table)
	tableSchema := string(ev.Table.Schema)
	tableInfo := h.tables.GetTable(tableSchema, tableName)
	if tableInfo == nil {
		return nil
	}

	action := ""
	var rows []map[string]interface{}
	switch e.Header.EventType {
	case replication.WRITE_ROWS_EVENTv1, replication.WRITE_ROWS_EVENTv2:
		// 插入
		action = InsertAction
		rows = h.RowDataToMap(ev.Rows, tableInfo)
	case replication.DELETE_ROWS_EVENTv1, replication.DELETE_ROWS_EVENTv2:
		// 删除
		action = DeleteAction
		rows = h.RowDataToMap(ev.Rows, tableInfo)
	case replication.UPDATE_ROWS_EVENTv1, replication.UPDATE_ROWS_EVENTv2:
		// 更新
		action = UpdateAction
		rows = h.RowDataToMap(ev.Rows, tableInfo)
	default:
		return nil
	}
	event := global.NewEvent(tableSchema, tableName, action, rows)
	return h.endpoint.Consume(event)
}

func (h handler) RowDataToMap(rows [][]interface{}, tableInfo *schema.Table) []map[string]interface{} {
	data := make([]map[string]interface{}, 0)
	for _, row := range rows {
		item := map[string]interface{}{}
		for i, c := range row {
			fieldName := tableInfo.Columns[i].Name
			switch c.(type) {
			case []uint8:
				item[fieldName] = string(c.([]byte))
			default:
				item[fieldName] = c
			}
		}
		data = append(data, item)
	}
	return data
}
