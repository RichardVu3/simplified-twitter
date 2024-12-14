package main

import (
	"encoding/json"
	"os"
	"simplified-twitter/server"
	"strconv"
)

func main() {
	config := server.Config{
		Encoder: json.NewEncoder(os.Stdout),
		Decoder: json.NewDecoder(os.Stdin),
	}
	if len(os.Args) == 1 {
		config.Mode = "s"
	} else {
		threadCount, _ := strconv.Atoi(os.Args[1])
		config.Mode = "p"
		config.ConsumersCount = threadCount
	}
	server.Run(config)
}
