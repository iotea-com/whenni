package main

import (
	"fmt"

	"github.com/iotea-com/iotea/services/controller/api"
)

func main() {
	server := api.New()
	fmt.Println("Starting controller server on port 9002...")
	server.Listen(9002)
}
