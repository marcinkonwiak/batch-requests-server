package main

import (
	"github.com/marcinkonwiak/batch-requests-server/config"
	"github.com/marcinkonwiak/batch-requests-server/server"
	"github.com/spf13/viper"
)

func main() {
	config.LoadConfig()
	s := server.NewServer()

	s.Logger.Fatal(s.Start(":" + viper.GetString("port")))
}
