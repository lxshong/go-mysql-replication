package service

import (
	"context"
	"github.com/go-mysql-org/go-mysql/replication"
	"go-mysql-replication/src/global"
	"go-mysql-replication/src/pos"
)

type replicationService struct {
	syconfig replication.BinlogSyncerConfig
	syncer   *replication.BinlogSyncer
	pos      pos.Pos
}

func NewReplicationService() *replicationService {
	return &replicationService{}
}

func (s *replicationService) initialize() error {
	config := global.Cfg()
	s.syconfig = replication.BinlogSyncerConfig{}
	s.syconfig.Host = config.Host
	s.syconfig.Port = config.Port
	s.syconfig.User = config.User
	s.syconfig.Password = config.Password
	s.syconfig.ServerID = config.SlaveID
	s.syconfig.Charset = config.Charset
	s.syconfig.Flavor = config.Flavor
	s.syncer = replication.NewBinlogSyncer(s.syconfig)
	p, err := pos.NewPos()
	if err != nil {
		return err
	}
	s.pos = p
	return nil
}

func (s *replicationService) run(callback EventCallBack) error {
	pst := s.pos.Get()
	streamer, _ := s.syncer.StartSync(*pst)

	for {
		ev, err := streamer.GetEvent(context.Background())
		if err != nil {
			return err
		}
		switch ev.Event.(type) {
		case *replication.RotateEvent:
			rotateEvent := ev.Event.(*replication.RotateEvent)
			pst.Name = string(rotateEvent.NextLogName)
			pst.Pos = uint32(rotateEvent.Position)
		case *replication.RowsEvent:
			pst.Pos = ev.Header.LogPos
			err = callback(ev)
			if err != nil {
				return err
			}
		default:
			pst.Pos = ev.Header.LogPos
		}
		//if err := s.pos.Save(); err != nil {
		//	return err
		//}

	}
	return nil
}
