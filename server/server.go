package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	signals := make(chan os.Signal, 1)
	done := make(chan struct{})
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals

		log.Println("Exit program")
		close(done)
	}()

	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			err := channel.PublishWithContext(ctx,
				"",         // exchange
				queue.Name, // routing key
				false,      // mandatory
				false,      // immediate
				amqp091.Publishing{
					ContentType: "text/plain",
					Body:        []byte("Hello world"),
				})
			if err != nil {
				log.Println(err.Error())
			}

			log.Println("Message is sended")
		}
	}
}
