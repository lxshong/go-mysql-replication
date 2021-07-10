package main

import (
	"flag"
	"fmt"
	"go-mysql-replication/src/global"
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
	fmt.Println(global.Cfg())
}
