package tools

import (
	"fmt"

	"github.com/streadway/amqp"
)

// Client RabbitMQ 客戶端
func NewClient(uri string) (*Client, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("mq dial failed: %w", err)
	}

	return &Client{conn: conn}, nil
}

// Channel 取得一個新的 channel
func (c *Client) Channel() (*amqp.Channel, error) {
	return c.conn.Channel()
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
