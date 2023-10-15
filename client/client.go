package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"testqueue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := channel.Consume(
		queue.Name,  // queue
		"BasbustPC", // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		log.Fatal(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for request := range msgs {
			log.Printf("Message: %#v\n", request.Body)
		}
	}()

	<-signals
	log.Println("Exit program")
}
