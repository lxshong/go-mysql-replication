package pos

import (
	"encoding/json"
	"errors"
	"github.com/go-mysql-org/go-mysql/mysql"
	"go-mysql-replication/src/global"
	"go-mysql-replication/src/tools"
	"strings"
)

type filePos struct {
	pos  *mysql.Position
	file string
}

func NewFilePos() Pos {
	return &filePos{file: global.Cfg().PosFile, pos: &mysql.Position{}}
}

func (receiver *filePos) Initialize() error {
	if len(strings.TrimSpace(receiver.file)) == 0 {
		return errors.New("pos file empty")
	}

	err := receiver.read()
	if err != nil {
		return err
	}
	//receiver.pos.Name = "master-bin.000496"
	//receiver.pos.Pos = 552911692
	return nil
}

func (receiver *filePos) Save() error {
	posStr, err := json.Marshal(receiver.pos)
	if err != nil {
		return err
	}
	err = tools.FilePutContent(receiver.file, string(posStr))
	if err != nil {
		return err
	}
	return nil
}

func (receiver *filePos) read() error {
	posStr, err := tools.FileGetContent(receiver.file)
	if err != nil {
		return err
	}

	if len(strings.TrimSpace(posStr)) == 0 {
		return nil
	}
	err = json.Unmarshal([]byte(posStr), &receiver.pos)
	if err != nil {
		return err
	}

	return nil
}

func (receiver *filePos) Get() *mysql.Position {
	return receiver.pos
}
