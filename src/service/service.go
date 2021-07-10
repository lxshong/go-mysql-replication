package service

import "github.com/go-mysql-org/go-mysql/replication"

var (
	_replicationService *replicationService
	_handle *handler
)

type EventCallBack func(event *replication.BinlogEvent) error

func init() {
	_replicationService = NewReplicationService()
	_handle = NewHandle(DefaultCap())
}

func Run() error {
	err := _replicationService.initialize()
	if err != nil {
		return err
	}

	err = _handle.initialize()
	if err != nil {
		return err
	}

	err = _replicationService.run(_handle.Callback)
	if err != nil {
		return err
	}
	return nil
}
