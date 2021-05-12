package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jun2900/online-library/database"
	"github.com/jun2900/online-library/models"
	"github.com/streadway/amqp"
)

const uploadPath = "./uploads"

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	//RabbitMq connection
	conn, err := amqp.Dial(database.RabbitMqUrl)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"insert_paper_content", // name
		false,                  // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			paper := &models.Paper{}
			json.Unmarshal(d.Body, paper)
			fmt.Println(paper.ID)

			f, err := os.Create(filepath.Join(uploadPath, fmt.Sprintf("%d-%s.pdf", paper.ID, paper.Title)))
			if err != nil {
				log.Printf("Error on processing paper content")
			}
			defer f.Close()

			if _, err := f.Write(paper.Content); err != nil {
				log.Printf("Error on writing the paper content")
			}
			if err := f.Sync(); err != nil {
				log.Printf("Error on sync commit the current paper content")
			}
			log.Printf("Paper inserted")
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
