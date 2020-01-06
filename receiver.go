package ramqp

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type Receiver interface {
	ExchangeType() string
	ExchangeName() string
	QueueName() string
	RouteKey() string
	OnReceive(*amqp.Delivery) bool
}

type receiverWrap struct {
	receiver Receiver
	channel  *amqp.Channel

	// queue declare params
	queueDurable   bool
	queueAutoDel   bool
	queueExclusive bool
	queueNoWait    bool
	queueArgs      amqp.Table

	// queue bind params
	bindNoWait bool
	bindArgs   amqp.Table

	prefetchCount int
	prefetchSize  int

	// consume params
	cAutoAck   bool
	cExclusive bool
	cNoLocal   bool
	cNoWait    bool
	cArgs      amqp.Table
}

func (r *receiverWrap) dealOptions(options []Opt) {
	if len(options) < 1 {
		return
	}

	for _, opt := range options {
		opt(r)
	}
}

func (r *receiverWrap) createMsgQueue(conn *amqp.Connection) (<-chan amqp.Delivery, error) {
	var err error
	r.channel, err = conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("create channel err: %v", err)
	}

	_, err = r.channel.QueueDeclare(
		r.receiver.QueueName(),
		r.queueDurable,
		r.queueAutoDel,
		r.queueExclusive,
		r.queueNoWait,
		r.queueArgs,
	)
	if err != nil {
		return nil, fmt.Errorf("queue init err: %v", err)
	}

	// todo create exchange and bind

	err = r.channel.QueueBind(
		r.receiver.QueueName(),
		r.receiver.RouteKey(),
		r.receiver.ExchangeName(),
		r.bindNoWait,
		r.bindArgs,
	)
	if err != nil {
		return nil, fmt.Errorf("bind queue err: %v", err)
	}

	err = r.channel.Qos(r.prefetchCount, r.prefetchSize, false)
	if err != nil {
		return nil, fmt.Errorf("set Qos err: %v", err)
	}

	msgQueue, err := r.channel.Consume(
		r.receiver.QueueName(),
		"",
		r.cAutoAck,
		r.cExclusive,
		r.cNoLocal,
		r.cNoWait,
		r.cArgs,
	)
	if err != nil {
		return nil, fmt.Errorf("consume err: %v", err)
	}

	return msgQueue, nil
}

func (r *receiverWrap) listen(conn *amqp.Connection) {
	msgQueue, err := r.createMsgQueue(conn)
	if err != nil {
		log.Println("create msg queue err: ", err)
		r.retry(conn)
		return
	}

	// block if success
	r.handleMsg(msgQueue)

	r.retry(conn)
}

func (r *receiverWrap) handleMsg(msgQueue <-chan amqp.Delivery) {
	closeErr := make(chan *amqp.Error)
	r.channel.NotifyClose(closeErr)

	for {
		select {
		case msg, valid := <-msgQueue:
			if valid {
				go func(msg amqp.Delivery) {
					if !r.receiver.OnReceive(&msg) {
						if !r.cAutoAck {
							log.Println("msg process err, requeue now !!!!!!!!!!")
							_ = msg.Nack(false, true)
						}
					}
					if !r.cAutoAck {
						_ = msg.Ack(false)
					}
				}(msg)
			} else {
				return
			}
		case <-closeErr:
			// end handle
			return
		}
	}
}

func (r *receiverWrap) retry(conn *amqp.Connection) {
	go func() {
		if !conn.IsClosed() {
			time.Sleep(5 * time.Second) // todo exponential backoff
			log.Println("receiver retry to listen ...")
			r.listen(conn)
		}
	}()
}
