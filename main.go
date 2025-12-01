package main

import (
	"fmt"

	"github.com/daitonium/go-blog-aggregator/internal/config"
)

func main() {
	conf := config.Read()
	conf.SetUser("mateo")
	conf = config.Read()
	fmt.Println(conf)
}
