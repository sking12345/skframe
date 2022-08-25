package bootstrap

import (
	"skframe/pkg/config"
	"skframe/pkg/rabbitMQ"
)

func SetRabbitMq() {
	rabbitMQ.ConnectRabbit(config.GetString("rabbit.url"))
}
