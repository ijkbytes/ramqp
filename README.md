# Ramqp
Ramqp是对Go版本[amqp](github.com/streadway/amqp)的简单封装，提供了简洁明了的调用方式，并且支持断线重连。理论上支持所有实现了amqp协议的消息中间件，比如著名的RabbitMQ。

# How to install

```
go get -u -v github.com/ijkbytes/ramqp
```

# How to use

下面的例子展示了如何使用Ramqp：

```go
...

msg := make(chan string)

mq := ramqp.New("amqp://test:test@127.0.0.1:5672/test")
receiver := &ramqp.Receiver{
        ExchangeName: "exchange_name",
        ExchangeType: "topic",
        QueueName:    "test_queue",
        RouteKey:     "test_queue",
        OnReceive:    func(msg *amqp.Delivery) bool {
            msg <- string(msg.Body)
            return true
        },
}
publisher = &rabbitmq.Publisher{}

// register receiver and enable auto ack
mq.RegisterReceiver(&ReceiverA{}, ramqp.WithConsumeAutoAck())
// register publisher
mq.RegisterPubliser(publisher)
if err := mq.Start(); err != nil {
    panic(err)
}

publisher.Publish("exchange_name", "test_queue", amqp.Publishing{
	Body: []byte("this is a test."),
})

fmt.Println("receive msg: ", <-msg)

...

```

# Features
- [x] TCP断线自动重连
- [x] 多通道共用一条TCP连接
- [x] 通道异常重连
- [x] 队列被删除自动重建
- [x] 实现消费者
- [x] 实现生产者
