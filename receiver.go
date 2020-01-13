package ramqp

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type Receiver struct {
	ExchangeType string
	ExchangeName string
	// The QueueName may be empty, in which case the server will
	// generate a unique name
	QueueName string
	RouteKey  string
	// OnReceive can not be nil, it will be invoke when receive msg
	OnReceive func(*amqp.Delivery) bool

	channel *amqp.Channel

	// queue declare params
	queueDurable   bool
	queueAutoDel   bool
	queueExclusive bool
	queueNoWait    bool
	queueArgs      amqp.Table

	// exchange declare params
	exchangeDurable  bool
	exchangeAutoDel  bool
	exchangeInternal bool
	exchangeNoWait   bool
	exchangeArgs     amqp.Table

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

func (r *Receiver) init(options []Opt) error {
	if r.OnReceive == nil {
		return fmt.Errorf("OnReceive can not be nil")
	}

	r.dealOptions(options)
	return nil
}

func (r *Receiver) dealOptions(options []Opt) {
	if len(options) < 1 {
		return
	}

	for _, opt := range options {
		opt(r)
	}
}

func (r *Receiver) createMsgQueue(conn *amqp.Connection) (<-chan amqp.Delivery, error) {
	var err error
	r.channel, err = conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("create channel err: %v", err)
	}

	_, err = r.channel.QueueDeclare(
		r.QueueName,
		r.queueDurable,
		r.queueAutoDel,
		r.queueExclusive,
		r.queueNoWait,
		r.queueArgs,
	)
	if err != nil {
		return nil, fmt.Errorf("queue init err: %v", err)
	}

	if len(r.ExchangeName) > 0 && len(r.ExchangeType) > 0 {
		err = r.channel.ExchangeDeclare(
			r.ExchangeName,
			r.ExchangeType,
			r.exchangeDurable,
			r.exchangeAutoDel,
			r.exchangeInternal,
			r.exchangeNoWait,
			r.exchangeArgs,
		)
		if err != nil {
			return nil, fmt.Errorf("exchange declare err: %v", err)
		}

		err = r.channel.QueueBind(
			r.QueueName,
			r.RouteKey,
			r.ExchangeName,
			r.bindNoWait,
			r.bindArgs,
		)
		if err != nil {
			return nil, fmt.Errorf("bind queue err: %v", err)
		}
	}

	err = r.channel.Qos(r.prefetchCount, r.prefetchSize, false)
	if err != nil {
		return nil, fmt.Errorf("set Qos err: %v", err)
	}

	msgQueue, err := r.channel.Consume(
		r.QueueName,
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

func (r *Receiver) listen(conn *amqp.Connection) {
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

func (r *Receiver) handleMsg(msgQueue <-chan amqp.Delivery) {
	closeErr := make(chan *amqp.Error)
	r.channel.NotifyClose(closeErr)

	for {
		select {
		case msg, valid := <-msgQueue:
			if valid {
				go func(msg amqp.Delivery) {
					if !r.OnReceive(&msg) {
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

func (r *Receiver) retry(conn *amqp.Connection) {
	go func() {
		if !conn.IsClosed() {
			time.Sleep(5 * time.Second) // todo exponential backoff
			log.Println("receiver retry to listen ...")
			r.listen(conn)
		}
	}()
}
