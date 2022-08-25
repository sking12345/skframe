package rabbitMQ

import (
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"skframe/pkg/logger"
	"sync"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	//队列名称
	QueueName string
	//交换机名称
	Exchange string
	//bind Key 名称
	Key   string
	Event map[string]chan int
}

var Client *RabbitMQ

var once sync.Once

const (
	stopChannel = 0
)

func ConnectRabbit(url string) {
	once.Do(func() {
		Client, _ = NewClient(url)
	})
}

func Close() {
	Client.conn.Close()
}

func (ptr *RabbitMQ) Close() {
	ptr.channel.Close()
	ptr.conn.Close()

}

func NewClient(url string) (*RabbitMQ, error) {
	con, err := amqp.Dial(url)
	if err != nil {
		logger.Error("rabbit", zap.Error(err))
		return nil, err
	}
	channel, err := con.Channel()
	if err != nil {
		logger.Error("rabbit", zap.Error(err))
		return nil, err
	}
	return &RabbitMQ{
		conn:    con,
		channel: channel,
		Event:   map[string]chan int{},
	}, nil
}

func (ptr *RabbitMQ) queueDeclare(queueName string) (amqp.Queue, error) {
	return ptr.channel.QueueDeclare(queueName, //是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞处理
		false,
		//额外的属性
		nil)
}

func (ptr *RabbitMQ) PushQueue(queueName string, message []byte) bool {
	_, err := ptr.queueDeclare(queueName)
	if err != nil {
		logger.Error("rabbit", zap.Error(err))
		return false
	}
	err = ptr.channel.Publish(
		ptr.Exchange,
		queueName,
		//如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
			//MessageId:
			//Timestamp: time.Now().,
		},
	)
	if err != nil {
		logger.Error("rabbit", zap.Error(err))
		return false
	}
	return true
}

func (ptr *RabbitMQ) ReceiveQueue(queueName string, handler func(data []byte, MessageId string) bool) bool {
	queue, err := ptr.queueDeclare(queueName)
	if err != nil {
		logger.Error("rabbit", zap.Error(err))
		return false
	}
	msgs, err := ptr.channel.Consume(
		queue.Name,
		"", // consumer
		//是否自动应答
		false, // auto-ack
		//是否独有
		false, // exclusive
		//设置为true，表示 不能将同一个Conenction中生产者发送的消息传递给这个Connection中 的消费者
		false, // no-local
		//列是否阻塞
		true, // no-wait
		nil,  // args
	)
	if err != nil {
		logger.Error("rabbit", zap.Error(err))
		return false
	}
	ptr.Event[queueName] = make(chan int, 1)
	go func() {
		select {
		case status := <-ptr.Event[queueName]:
			if status == stopChannel {
				logger.Info("rabbit", zap.String(queueName, "end"))
				return
			}
		case <-msgs:
			for data := range msgs {
				data.Reject(!handler(data.Body, data.MessageId))
			}
		}
	}()

	return true
}

func (ptr *RabbitMQ) setExchange(exChangeName string, kind string) bool {
	err := ptr.channel.ExchangeDeclare(
		exChangeName,
		kind,
		true,
		false,
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error("rabbit", zap.Error(err))
		return false
	}
	return true
}

func (ptr *RabbitMQ) exChangeModelServer(exChangeName, kind, routeKey string, data []byte) bool {
	if ptr.setExchange(exChangeName, kind) == false {
		return false
	}
	err := ptr.channel.Publish(
		exChangeName,
		routeKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain", //application/json
			Body:        data,
		},
	)
	if err != nil {
		logger.Error("rabbit", zap.Error(err))
		return false
	}
	return true
}

func (ptr *RabbitMQ) exChangeClient(exChangeName, kind, queueName, routeKey string, handler func(data []byte) (ack bool)) bool {
	if ptr.setExchange(exChangeName, kind) == false {
		return false
	}
	queue, err := ptr.queueDeclare(queueName) //随机生成队名
	if err != nil {
		logger.Error("rabbit", zap.Error(err))
		return false
	}

	//绑定队列到 exchange 中
	err = ptr.channel.QueueBind(
		queue.Name,
		routeKey,
		exChangeName,
		false,
		nil,
	)
	if err != nil {
		logger.Error("rabbit", zap.Error(err))
		return false
	}
	//消费消息
	messages, err := ptr.channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		true,
		nil,
	)
	ptr.Event[queueName] = make(chan int, 1)
	go func() {
		select {
		case status := <-ptr.Event[queueName]:
			if status == stopChannel {
				logger.Info("rabbit", zap.String(queueName, "end"))
				return
			}
		case <-messages:
			for data := range messages {
				data.Reject(!handler(data.Body))
			}
		}
	}()
	return true
}

func (ptr *RabbitMQ) PushPublish(exChangeName string, msgType string, data []byte) bool { //没有队列
	return ptr.exChangeModelServer(exChangeName, "fanout", "", data)
}

func (ptr *RabbitMQ) ReceivePublish(exChangeName string, queueName string, handler func(data []byte) (ack bool)) bool {
	return ptr.exChangeClient(exChangeName, "fanout", queueName, "", handler)
}

func (ptr *RabbitMQ) PushRoute(exChangeName string, routeKey string, data []byte) bool { //有队列
	return ptr.exChangeModelServer(exChangeName, "direct", routeKey, data)
}

func (ptr *RabbitMQ) ReceiveRoute(exChangeName string, routeKey string, queueName string, handler func(data []byte) (ack bool)) bool {

	return ptr.exChangeClient(exChangeName, "direct", queueName, routeKey, handler)
}

func (ptr *RabbitMQ) PushTopic(exChangeName string, routeKey string, data []byte) bool { //队列匹配模式
	return ptr.exChangeModelServer(exChangeName, "direct", routeKey, data)
}

func (ptr *RabbitMQ) ReceiveTopic(exChangeName string, routeKey string, queueName string, handler func(data []byte) (ack bool)) bool {
	return ptr.exChangeClient(exChangeName, "direct", queueName, routeKey, handler)
}

func (ptr *RabbitMQ) ClearQueueMessage(queueName string) bool { //清除队列message
	_, err := ptr.channel.QueueDelete(queueName, true, true, true)
	if err != nil {
		logger.Error("rabbit", zap.Error(err))
		return false
	}
	return true
}

func (ptr *RabbitMQ) QueueUnbind(queueName string, RoutingKey string, exchange string) bool { //取消交互机与队列的绑定
	if err := ptr.channel.QueueUnbind(queueName, RoutingKey, exchange, nil); err != nil {
		logger.Error("rabbit", zap.Error(err))
		return false
	}
	return true
}

func (ptr *RabbitMQ) QuitReceive(queueName string) {
	if ptr.Event[queueName] != nil {
		ptr.Event[queueName] <- stopChannel
	}

}
