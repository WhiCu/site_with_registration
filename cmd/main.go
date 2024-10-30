package main

import (
	"reg/internal/server"
)

func main() {
	serv := server.New()
	serv.Go()
}
