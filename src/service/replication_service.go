package service

import "github.com/go-mysql-org/go-mysql/replication"

type ReplicationService struct {
	syconfig replication.BinlogSyncerConfig
	syncer *replication.BinlogSyncer
}

func (s *ReplicationService) initialize() error {
	s.syconfig = replication.BinlogSyncerConfig{}

	s.syncer = replication.NewBinlogSyncer(s.syconfig)
	return nil
}
