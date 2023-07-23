package main

import (
	"fmt"
	"os"

	envserver "github.com/codescalersinternships/envserver-Hanya/internal"
)

func main() {
	server, err := envserver.NewServer(3000)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = server.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
