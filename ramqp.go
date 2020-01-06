package ramqp

import (
	"github.com/streadway/amqp"
	"go.uber.org/atomic"
	"log"
	"time"
)

type Ramqp struct {
	url   string
	conn  *amqp.Connection
	stop  atomic.Bool
	stopC chan struct{}

	receivers []*receiverWrap

	// todo publisher
	//publisher []*publisherWrap
}

func New(url string) *Ramqp {
	return &Ramqp{
		url:   url,
		stopC: make(chan struct{}),
	}
}

func (mq *Ramqp) RegisterReceiver(recv Receiver, options ...Opt) {
	r := &receiverWrap{
		receiver: recv,
	}
	r.dealOptions(options)
	mq.receivers = append(mq.receivers, r)
}

func (mq *Ramqp) refresh() error {
	var err error
	mq.conn, err = amqp.Dial(mq.url)
	if err != nil {
		log.Println("connection err: ", err)
		return err
	}

	for _, recv := range mq.receivers {
		go recv.listen(mq.conn)
	}

	closeErr := make(chan *amqp.Error)
	mq.conn.NotifyClose(closeErr)

	// keep retry
	go func() {
		select {
		case <-closeErr:
		RETRY:
			if !mq.stop.Load() {
				if err := mq.refresh(); err != nil {
					log.Println("connection has been closed, retry connect now ...")
					time.Sleep(5 * time.Second) // todo exponential backoff
					goto RETRY
				}
			}
		case <-mq.stopC:
			_ = mq.conn.Close()
		}
	}()

	return nil
}

func (mq *Ramqp) Start() error {
	mq.stop.Store(false)
	return mq.refresh()
}

func (mq *Ramqp) Stop() {
	mq.stop.Store(true)
	mq.stopC <- struct{}{}
}