package ramqp

import "github.com/streadway/amqp"

type Opt func(*Receiver)

func WithQueueDurable() Opt {
	return func(wrap *Receiver) {
		wrap.queueDurable = true
	}
}

func WithQueueAutoDel() Opt {
	return func(wrap *Receiver) {
		wrap.queueAutoDel = true
	}
}

func WithQueueExclusive() Opt {
	return func(wrap *Receiver) {
		wrap.queueExclusive = true
	}
}

func WithQueueNoWait() Opt {
	return func(wrap *Receiver) {
		wrap.queueNoWait = true
	}
}

func WithQueueArgs(args amqp.Table) Opt {
	return func(wrap *Receiver) {
		wrap.queueArgs = args
	}
}

func WithBindNoWait() Opt {
	return func(wrap *Receiver) {
		wrap.bindNoWait = true
	}
}

func WithBindArgs(args amqp.Table) Opt {
	return func(wrap *Receiver) {
		wrap.bindArgs = args
	}
}

func WithPrefetchCount(c int) Opt {
	return func(wrap *Receiver) {
		wrap.prefetchCount = c
	}
}

func WithPrefetchSize(s int) Opt {
	return func(wrap *Receiver) {
		wrap.prefetchSize = s
	}
}

func WithConsumeAutoAck() Opt {
	return func(wrap *Receiver) {
		wrap.cAutoAck = true
	}
}

func WithConsumeExclusive() Opt {
	return func(wrap *Receiver) {
		wrap.cExclusive = true
	}
}

func WithConsumeNoLocal() Opt {
	return func(wrap *Receiver) {
		wrap.cNoLocal = true
	}
}

func WithConsumeNoWait() Opt {
	return func(wrap *Receiver) {
		wrap.cNoWait = true
	}
}

func WithConsumeArgs(args amqp.Table) Opt {
	return func(wrap *Receiver) {
		wrap.cArgs = args
	}
}

func WithExchangeDurable() Opt {
	return func(wrap *Receiver) {
		wrap.exchangeDurable = true
	}
}

func WithExchangeAutoDel() Opt {
	return func(wrap *Receiver) {
		wrap.exchangeAutoDel = true
	}
}

func WithExchangeInternal() Opt {
	return func(wrap *Receiver) {
		wrap.exchangeInternal = true
	}
}

func WithExchangeNoWait() Opt {
	return func(wrap *Receiver) {
		wrap.exchangeNoWait = true
	}
}

func WithExchangeArgs(args amqp.Table) Opt {
	return func(wrap *Receiver) {
		wrap.exchangeArgs = args
	}
}

func WithPExchangeDurable() POpt {
	return func(wrap *Publisher) {
		wrap.exchangeDurable = true
	}
}

func WithPExchangeAutoDel() POpt {
	return func(wrap *Publisher) {
		wrap.exchangeAutoDel = true
	}
}

func WithPExchangeInternal() POpt {
	return func(wrap *Publisher) {
		wrap.exchangeInternal = true
	}
}

func WithPExchangeNoWait() POpt {
	return func(wrap *Publisher) {
		wrap.exchangeNoWait = true
	}
}

func WithPExchangeArgs(args amqp.Table) POpt {
	return func(wrap *Publisher) {
		wrap.exchangeArgs = args
	}
}

func WithPMandatory() POpt {
	return func(wrap *Publisher) {
		wrap.publishMandatory = true
	}
}
