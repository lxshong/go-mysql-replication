package endpoint

import (
	"fmt"
	"go-mysql-replication/src/global"
)

type stdioEndpoint struct {
}

func (r stdioEndpoint) Consume(event global.Event) error {
	fmt.Println(event)
	return nil
}

func NewStdioEndpoint() Endpoint {
	return &stdioEndpoint{}
}
