package tools

import "github.com/streadway/amqp"

type Client struct {
	conn *amqp.Connection
}
