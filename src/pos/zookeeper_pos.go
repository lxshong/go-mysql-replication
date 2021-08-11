package pos

import (
	"encoding/json"
	"errors"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/samuel/go-zookeeper/zk"
	"go-mysql-replication/src/global"
	"strings"
	"time"
)

type zkPos struct {
	pos     *mysql.Position
	conn    *zk.Conn
	servers []string
	zkpath  string
}

func NewZkPos() Pos {
	return &zkPos{conn: nil, pos: &mysql.Position{}, servers: nil}
}

func (receiver *zkPos) createZkConn() (*zk.Conn, error) {
	conn, _, err := zk.Connect(receiver.servers, 5*time.Second)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (receiver *zkPos) Initialize() error {
	receiver.servers = global.Cfg().PosZk
	if len(receiver.servers) == 0 {
		return errors.New("zk servers empty")
	}
	conn, err := receiver.createZkConn()
	if err != nil {
		return err
	}
	receiver.conn = conn
	receiver.zkpath = "/go-mysql-replication"
	return nil
}

func (receiver *zkPos) Save() error {
	posStr, err := json.Marshal(receiver.pos)
	if err != nil {
		return err
	}
	if b, err := existsNode(receiver.conn, receiver.zkpath); err != nil {
		return err
	} else {
		if b {
			if err := updateNode(receiver.conn, receiver.zkpath, []byte(posStr)); err != nil {
				return err
			}
		} else {
			if err := createNode(receiver.conn, receiver.zkpath, []byte(posStr)); err != nil {
				return err
			}
		}
		return nil
	}
}

func (receiver *zkPos) read() error {
	posStr, err := getNode(receiver.conn, receiver.zkpath)
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

func (receiver *zkPos) Get() *mysql.Position {
	return receiver.pos
}

// 创建节点
func createNode(conn *zk.Conn, path string, data []byte) error {
	var flags int32 = 0
	acls := zk.WorldACL(zk.PermAll)
	_, err := conn.Create(path, data, flags, acls)
	return err
}

// 更新节点
func updateNode(conn *zk.Conn, path string, data []byte) error {
	_, sate, _ := conn.Get(path)
	_, err := conn.Set(path, data, sate.Version)
	return err
}

// 节点是否存在
func existsNode(conn *zk.Conn, path string) (bool, error) {
	b, _, err := conn.Exists(path)
	return b, err
}

// 删除节点
func deleteNode(conn *zk.Conn, path string) error {
	_, sate, _ := conn.Get(path)
	return conn.Delete(path, sate.Version)
}

// 查询
func getNode(conn *zk.Conn, path string) (string, error) {
	bs, _, err := conn.Get(path)
	return string(bs), err
}
