package core

import (
	"github.com/streadway/amqp"
)

type Queue interface {
	Publish(content string) error
	Consume(handler func(content string) (ack bool, requeue bool, err error)) error
	Close() error
	Clear() error
}

type RbQueue struct {
	ExchangeName string
	QueueName    string
	RoutingKey   string
	conn         *amqp.Connection
	channel      *amqp.Channel
}

func NewRbQueue(url string, exchangeName string, queueName string, routingKey string, prefetchCount int) (*RbQueue, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	err = ch.Qos(prefetchCount, 0, false)
	if err != nil {
		return nil, err
	}
	err = ch.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	err = ch.QueueBind(queueName, routingKey, exchangeName, false, nil)
	if err != nil {
		return nil, err
	}
	q := &RbQueue{
		exchangeName,
		queueName,
		routingKey,
		conn,
		ch,
	}
	return q, nil
}

func (q RbQueue) Publish(content string) error {
	return q.channel.Publish(q.ExchangeName,
		q.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(content),
			DeliveryMode: 2,                    // 持久化
			Priority:     GetPriority(content), // 优先级
		})
}

func (q RbQueue) Consume(handler func(content string) (ack bool, requeue bool, err error)) error {
	msgs, err := q.channel.Consume(q.QueueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	for d := range msgs {
		msg := d
		ack, requeue, err := handler(string(msg.Body))

		if err != nil {
			msg.Nack(false, requeue)
			return err
		}
		if ack {
			msg.Ack(false)
		} else {
			msg.Nack(false, requeue)
		}
	}
	return nil
}

func (q RbQueue) Close() error {
	return q.conn.Close()
}

func (q RbQueue) Clear() error {
	_, err := q.channel.QueuePurge(q.QueueName, false)
	return err
}
