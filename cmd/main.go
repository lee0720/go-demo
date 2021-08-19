package main

import (
	"flag"
	"strings"

	"gitlab.com/lilh/go-demo/internal/config"
	"gitlab.com/lilh/go-demo/internal/nsq/consumer"
	"gitlab.com/lilh/go-demo/internal/nsq/producer"
)

var configFileName = flag.String("cfn", "config", "name of configs file")
var configFilePath = flag.String("cfp", "./configs", "path of configs file")

func main() {
	flag.Parse()

	err := config.InitConfig(*configFileName, strings.Split(*configFilePath, ","))

	if err != nil {
		panic(err)
	}

	go producer.Producer()

	consumer.Consumer()

}
