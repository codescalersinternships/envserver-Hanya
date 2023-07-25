package main

import (
	"flag"
	"log"

	envserver "github.com/codescalersinternships/envserver-Hanya/internal"
)

func main() {
	var port int
	flag.IntVar(&port,"p",8080,"port that will be used to run the app")
	server, err := envserver.NewServer(port)
	if err != nil {
		log.Fatal(err)
	}
	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}

}
