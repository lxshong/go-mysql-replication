package service

import (
	"fmt"
	"github.com/go-mysql-org/go-mysql/replication"
	"go-mysql-replication/src/global"
	"log"
	"os"
)

const (
	DEFAULT_CAP = 4096
)

type handler struct {
	cap   int
	queue chan *replication.BinlogEvent
	stop  chan struct{}
	tables global.Tables
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
	err := h.startListener()
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
	log.Fatalln("start listener")
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

	if !h.tables.Match(tableSchema, tableName) {
		return nil
	}

	fmt.Println(tableSchema,tableName)
	switch e.Header.EventType {
	case replication.WRITE_ROWS_EVENTv1, replication.WRITE_ROWS_EVENTv2:
		// 插入
	case replication.DELETE_ROWS_EVENTv1, replication.DELETE_ROWS_EVENTv2:
		// 删除
	case replication.UPDATE_ROWS_EVENTv1, replication.UPDATE_ROWS_EVENTv2:
		// 更新
	}
	os.Exit(1)
	return nil
}
