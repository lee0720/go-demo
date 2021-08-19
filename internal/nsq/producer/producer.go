package producer

import (
	"fmt"
	"time"

	"github.com/nsqio/go-nsq"
	"gitlab.com/lilh/go-demo/internal/config"
)

func Producer() {
	nsqConfig := nsq.NewConfig()
	p, err := nsq.NewProducer(config.GetConfig().NSQConfig.Host+":"+config.GetConfig().NSQConfig.Port, nsqConfig)

	if err != nil {
		panic(err)
	}
	for i := 0; i < 1000; i++ {
		content := fmt.Sprintf("lilh%d", i)
		err = p.Publish(config.GetConfig().NSQConfig.Topic, []byte(content))
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}
}
