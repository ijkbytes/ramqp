package ramqp

import (
	"fmt"
	"github.com/streadway/amqp"
	"testing"
	"time"
)

func TestRabbitMQ_Receiver(t *testing.T) {
	mq := New("amqp://test:test@127.0.0.1:5672/test")

	receiver1 := &Receiver{
		ExchangeType: "topic",
		ExchangeName: "exchange_name",
		QueueName:    "test_queue",
		RouteKey:     "test_queue",
		OnReceive: func(msg *amqp.Delivery) bool {
			fmt.Println("Receiver1 receive msg: ", string(msg.Body))
			return true
		},
	}

	receiver2 := &Receiver{
		ExchangeType: "topic",
		ExchangeName: "exchange_name",
		QueueName:    "test_queue",
		RouteKey:     "test_queue",
		OnReceive: func(msg *amqp.Delivery) bool {
			fmt.Println("Receiver2 receive msg: ", string(msg.Body))
			return true
		},
	}

	mq.RegisterReceiver(receiver1)
	mq.RegisterReceiver(receiver2, WithConsumeAutoAck())
	if err := mq.Start(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(3 * time.Minute)

	mq.Stop()

	t.Log("Stop RabbitMQ")

	time.Sleep(2 * time.Minute)

}
