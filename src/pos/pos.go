package pos

import (
	"github.com/go-mysql-org/go-mysql/mysql"
	"go-mysql-replication/src/global"
)

type Pos interface {
	Initialize() error
	Save() error
	Get()  *mysql.Position
}

func NewPos() (Pos, error) {
	switch {
	case global.PosIsFile():
		pos := NewFilePos()
		err := pos.Initialize()
		if err != nil {
			return nil, err
		}
		return pos, nil
	case global.PosIsZk():
		pos := NewZkPos()
		err := pos.Initialize()
		if err != nil {
			return nil, err
		}
		return pos, nil
	default:
		pos := NewFilePos()
		err := pos.Initialize()
		if err != nil {
			return nil, err
		}
		return pos, nil
	}
}
