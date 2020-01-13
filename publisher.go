package ramqp

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type publishing struct {
	content  amqp.Publishing
	key      string
	exchange string
}

type POpt func(p *Publisher)

type Publisher struct {
	ExchangeType string
	ExchangeName string
	NotifyReturn func(msg amqp.Return)

	channel *amqp.Channel
	msgChan chan publishing
	retErr  chan error

	// exchange declare params
	exchangeDurable  bool
	exchangeAutoDel  bool
	exchangeInternal bool
	exchangeNoWait   bool
	exchangeArgs     amqp.Table

	publishMandatory bool
}

func (p *Publisher) init(options []POpt) {
	p.retErr = make(chan error)
	p.msgChan = make(chan publishing)

	if len(options) < 1 {
		return
	}

	for _, opt := range options {
		opt(p)
	}
}

func (p *Publisher) createChannel(conn *amqp.Connection) error {
	var err error
	p.channel, err = conn.Channel()
	if err != nil {
		return err
	}

	if len(p.ExchangeName) > 0 && len(p.ExchangeType) > 0 {
		err = p.channel.ExchangeDeclare(
			p.ExchangeName,
			p.ExchangeType,
			p.exchangeDurable,
			p.exchangeAutoDel,
			p.exchangeInternal,
			p.exchangeNoWait,
			p.exchangeArgs,
		)
		if err != nil {
			return fmt.Errorf("exchange declare err: %v", err)
		}
	}

	return nil
}

func (p *Publisher) listen(conn *amqp.Connection) {
	err := p.createChannel(conn)
	if err != nil {
		log.Println("create msg queue err: ", err)
		p.retry(conn)
	}

	p.handleMsg()

	p.retry(conn)
}

func (p *Publisher) handleMsg() {
	closeErr := make(chan *amqp.Error)
	p.channel.NotifyClose(closeErr)
	notifyReturn := make(chan amqp.Return)
	p.channel.NotifyReturn(notifyReturn)

	for {
		select {
		case msg, valid := <-p.msgChan:
			if valid {
				err := p.channel.Publish(msg.exchange, msg.key, p.publishMandatory, false, msg.content)
				p.retErr <- err
			}
		case Return, ok := <-notifyReturn:
			if p.NotifyReturn != nil && ok {
				go p.NotifyReturn(Return)
			}
		case <-closeErr:
			return
		}
	}
}

func (p *Publisher) retry(conn *amqp.Connection) {
	go func() {
		if !conn.IsClosed() {
			time.Sleep(5 * time.Second) // todo exponential backoff
			log.Println("publisher retry to listen ...")
			p.listen(conn)
		}
	}()
}

func (p *Publisher) Publish(exchange, key string, msg amqp.Publishing) (err error) {
	pMsg := publishing{
		exchange: exchange,
		key:      key,
		content:  msg,
	}

	p.msgChan <- pMsg
	err = <-p.retErr
	return err
}
