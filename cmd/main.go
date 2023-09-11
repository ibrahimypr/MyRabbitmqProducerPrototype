package main

import "smtp/rabbitmq"

func main() {
	rmq := rabbitmq.InitalizeRabbitMQ()
	defer rmq.Close()
	rabbitmq.SendMessage(rmq, "test", "hello world!")
}