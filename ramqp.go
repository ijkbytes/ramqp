package ramqp

import (
	"github.com/streadway/amqp"
	"go.uber.org/atomic"
	"log"
	"math/rand"
	"strings"
	"time"
)

type Ramqp struct {
	urls   []string
	conn  *amqp.Connection
	stop  atomic.Bool
	stopC chan struct{}

	receivers []*Receiver

	publishers []*Publisher
}

// supports multiple urls, separated by commas. For example:
// "amqp://guest:guest@10.0.1.21:5672, amqp://guest:guest@10.0.1.22:5672"
func New(urls string) *Ramqp {
	return &Ramqp{
		urls:   strings.Split(urls, ","),
		stopC: make(chan struct{}),
	}
}

func (mq *Ramqp) RegisterReceiver(recv *Receiver, options ...Opt) error {
	if err := recv.init(options); err != nil {
		return err
	}
	mq.receivers = append(mq.receivers, recv)
	return nil
}

func (mq *Ramqp) RegisterPubliser(pub *Publisher, options ...POpt) error {
	if err := pub.init(options); err != nil {
		return err
	}
	mq.publishers = append(mq.publishers, pub)
	return nil
}

func (mq *Ramqp) randomUrl() string {
	rand.Seed(time.Now().UnixNano())
	i := rand.Int63n(int64(len(mq.urls)))
	return strings.TrimSpace(mq.urls[i])
}

func (mq *Ramqp) refresh() error {
	var err error
	if mq.conn == nil || mq.conn.IsClosed() {
		mq.conn, err = amqp.Dial(mq.randomUrl())
	}
	if err != nil {
		log.Println("connection err: ", err)
		return err
	}

	for _, recv := range mq.receivers {
		go recv.listen(mq.conn)
	}

	for _, pub := range mq.publishers {
		go pub.listen(mq.conn)
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
	close(mq.stopC)
}
