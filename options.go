package ramqp

import "github.com/streadway/amqp"

type Opt func(*receiverWrap)

func WithQueueDurable() Opt {
	return func(wrap *receiverWrap) {
		wrap.queueDurable = true
	}
}

func WithQueueAutoDel() Opt {
	return func(wrap *receiverWrap) {
		wrap.queueAutoDel = true
	}
}

func WithQueueExclusive() Opt {
	return func(wrap *receiverWrap) {
		wrap.queueExclusive = true
	}
}

func WithQueueNoWait() Opt {
	return func(wrap *receiverWrap) {
		wrap.queueNoWait = true
	}
}

func WithQueueArgs(args amqp.Table) Opt {
	return func(wrap *receiverWrap) {
		wrap.queueArgs = args
	}
}

func WithBindNoWait() Opt {
	return func(wrap *receiverWrap) {
		wrap.bindNoWait = true
	}
}

func WithBindArgs(args amqp.Table) Opt {
	return func(wrap *receiverWrap) {
		wrap.bindArgs = args
	}
}

func WithPrefetchCount(c int) Opt {
	return func(wrap *receiverWrap) {
		wrap.prefetchCount = c
	}
}

func WithPrefetchSize(s int) Opt {
	return func(wrap *receiverWrap) {
		wrap.prefetchSize = s
	}
}

func WithConsumeAutoAck() Opt {
	return func(wrap *receiverWrap) {
		wrap.cAutoAck = true
	}
}

func WithConsumeExclusive() Opt {
	return func(wrap *receiverWrap) {
		wrap.cExclusive = true
	}
}

func WithConsumeNoLocal() Opt {
	return func(wrap *receiverWrap) {
		wrap.cNoLocal = true
	}
}

func WithConsumeNoWait() Opt {
	return func(wrap *receiverWrap) {
		wrap.cNoWait = true
	}
}

func WithConsumeArgs(args amqp.Table) Opt {
	return func(wrap *receiverWrap) {
		wrap.cArgs = args
	}
}
