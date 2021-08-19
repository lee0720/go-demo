package consumer

import (
	"fmt"

	"github.com/nsqio/go-nsq"
	"gitlab.com/lilh/go-demo/internal/config"
)

func Consumer() {
	nsqConfig := nsq.NewConfig()
	c, err := nsq.NewConsumer(config.GetConfig().NSQConfig.Topic, config.GetConfig().NSQConfig.ConsumerChannel, nsqConfig)
	if err != nil {
		panic(err)
	}
	c2, err := nsq.NewConsumer(config.GetConfig().NSQConfig.Topic, config.GetConfig().NSQConfig.ConsumerChannel, nsqConfig)
	if err != nil {
		panic(err)
	}

	c.AddHandler(&TestHandler{
		Name: "consumer1",
	})

	c2.AddHandler(&TestHandler{
		Name: "consumer2",
	})

	err = c.ConnectToNSQD(config.GetConfig().NSQConfig.Host + ":" + config.GetConfig().NSQConfig.Port)
	if err != nil {
		panic(err)
	}

	err = c2.ConnectToNSQD(config.GetConfig().NSQConfig.Host + ":" + config.GetConfig().NSQConfig.Port)
	if err != nil {
		panic(err)
	}
	select {}
}

type TestHandler struct {
	Name string
}

func (h *TestHandler) HandleMessage(message *nsq.Message) error {
	listID := string(message.Body)
	fmt.Println(h.Name, " has consumed this message, message content is ", listID)
	message.Finish()
	return nil
}
