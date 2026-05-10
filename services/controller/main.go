package main

import (
	"fmt"

	"github.com/iotea-com/iotea/services/controller/api"
)

func main() {
	server := api.New()
	fmt.Println("Starting controller server on port 8080...")
	server.Listen(8080)
}
