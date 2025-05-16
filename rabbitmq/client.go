package tools

import (
	"fmt"

	"github.com/streadway/amqp"
)

// Client RabbitMQ 客戶端
func NewRabbitMQClient(uri string) (*Client, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("mq dial failed: %w", err)
	}

	return &Client{conn: conn}, nil
}

// Publish 發送訊息
func (c *Client) Publish(exchange, routingKey string, body []byte) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return ch.Publish(
		exchange,   // exchange
		routingKey, // routingKey
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

// Close 關閉連線
func (c *Client) Close() error {
	return c.conn.Close()
}
