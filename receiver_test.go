package ramqp

import (
	"fmt"
	"github.com/streadway/amqp"
	"testing"
	"time"
)

type ReceiverA struct {
}

func (r *ReceiverA) ExchangeType() string {
	return "topic"
}

func (r *ReceiverA) ExchangeName() string {
	return "exchange_name"
}

func (r *ReceiverA) QueueName() string {
	return "test_queue"
}

func (r *ReceiverA) RouteKey() string {
	return "test_queue"
}

func (r *ReceiverA) OnReceive(msg *amqp.Delivery) bool {
	fmt.Println("ReceiverA receive msg: ", string(msg.Body))
	return true
}

type ReceiverB struct {
}

func (r *ReceiverB) ExchangeType() string {
	return "topic"
}

func (r *ReceiverB) ExchangeName() string {
	return "exchange_name"
}

func (r *ReceiverB) QueueName() string {
	return "test_queue"
}

func (r *ReceiverB) RouteKey() string {
	return "test_queue"
}

func (r *ReceiverB) OnReceive(msg *amqp.Delivery) bool {
	fmt.Println("ReceiverB receive msg: ", string(msg.Body))
	return false
}

func TestRabbitMQ_Receiver(t *testing.T) {
	mq := New("amqp://test:test@127.0.0.1:5672/test")

	mq.RegisterReceiver(&ReceiverA{})
	mq.RegisterReceiver(&ReceiverA{}, WithConsumeAutoAck())
	mq.RegisterReceiver(&ReceiverB{})
	mq.RegisterReceiver(&ReceiverB{}, WithConsumeAutoAck())
	if err := mq.Start(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Minute)

	mq.Stop()

	t.Log("Stop RabbitMQ")

	time.Sleep(1 * time.Minute)

}
