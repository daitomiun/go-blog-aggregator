package main

import (
	"fmt"
	"log"

	"github.com/daitonium/go-blog-aggregator/internal/config"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatalf("Cannot read config: %v", err)
	}
	conf.SetUser("mateo")
	conf, err = config.Read()
	if err != nil {
		log.Fatalf("Cannot read config: %v", err)
	}
	fmt.Println(conf)
}
