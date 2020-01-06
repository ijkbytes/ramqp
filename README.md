# Ramqp
Ramqp是对Go版本[amqp](github.com/streadway/amqp)的简单封装，提供了简洁明了的调用方式，并且支持断线重连。理论上支持所有实现了amqp协议的消息中间件，比如著名的RabbitMQ。

# How to install

```
go get -u -v github.com/ijkbytes/ramqp
```

# How to use

下面的例子展示了如何使用Ramqp：

```go
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

...

mq := ramqp.New("amqp://test:test@127.0.0.1:5672/test")
// enable auto ack
mq.RegisterReceiver(&ReceiverA{}, ramqp.WithConsumeAutoAck())
if err := mq.Start(); err != nil {
    panic(err)
}

...

```

# Features
- [x] TCP断线自动重连
- [x] 多通道共用一条TCP连接
- [x] 通道异常重连
- [x] 队列被删除自动重建
- [x] 实现消费者
- [ ] 实现实现生产者
- [ ] ...
