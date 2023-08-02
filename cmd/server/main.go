package main

import (
	"fmt"

	"github.com/Nimbo1999/go-apis-go-expert/configs"
)

func main() {
	config, err := configs.LoadConfig("./")
	if err != nil {
		panic(err)
	}
	fmt.Println(config.DBDriver)
}
