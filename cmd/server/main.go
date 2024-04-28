package main

import (
	"github.com/marcinkonwiak/batch-requests-server/server"
)

func main() {
	s := server.NewServer()

	s.Logger.Fatal(s.Start(":1323"))
}
