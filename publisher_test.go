package ramqp

import (
	"fmt"
	"github.com/streadway/amqp"
	"testing"
	"time"
)

func TestRabbitMQ_Publiser(t *testing.T) {
	success := make(chan struct{})
	failed := time.Tick(5 * time.Second)

	mq := New("amqp://test:test@127.0.0.1:5672/test")
	publisher := &Publisher{
		ExchangeName: "exchange_name",
		ExchangeType: "topic",
		NotifyReturn: func(msg amqp.Return) {
			fmt.Println("Message route err, msg: ", string(msg.Body))
			success <- struct{}{}
		},
	}
	mq.RegisterPublisher(publisher, WithPMandatory())
	if err := mq.Start(); err != nil {
		t.Fatal(err)
	}

	err := publisher.Publish("exchange_name", "no_queue", amqp.Publishing{
		ContentType: "application/octet-stream",
		Body:        []byte("this is a test !"),
	})
	if err != nil {
		t.Fatal(err)
	}

	select {
	case <-success:
		break
	case <-failed:
		t.Fatal("NotifyReturn hasn't been invoke")
	}

	mq.Stop()
	time.Sleep(2 * time.Second)
}
