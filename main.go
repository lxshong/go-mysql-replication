package main

import (
	"flag"
	"fmt"
	"go-mysql-replication/src/global"
	"go-mysql-replication/src/service"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "config.yml", "application config file")
}
func main() {

	err := global.InitConfig(configPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = service.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
