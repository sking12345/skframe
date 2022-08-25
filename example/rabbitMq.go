package example

import (
	"fmt"
	"skframe/bootstrap"
	"skframe/pkg/rabbitMQ"
	"time"
)

func TestRabbitMq() {
	bootstrap.SetRabbitMq()
	rabbitMQ.Client.ReceiveQueue("queueName1", func(data []byte, string2 string) (ack bool) {
		fmt.Println(string(data))
		return false
	})
	for i := 0; i < 20; i++ {
		rabbitMQ.Client.PushQueue("queueName1", []byte(fmt.Sprintf("hell-%d", i)))
	}

	rabbitMQ.Client.ClearQueueMessage("queuenamexxx1")
	rabbitMQ.Client.QuitReceive("queuenamexxx1")

	time.Sleep(100 * time.Second)
	rabbitMQ.Client.Close()
}
