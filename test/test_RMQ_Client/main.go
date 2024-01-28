package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func main() {
	//connectStr := os.Args[0]
	connectStr := ""
	//servExchange := os.Args[1]
	//queueName := os.Args[1]
	queueName := ""

	fmt.Println(connectStr)
	fmt.Println(queueName)

	conn, err := amqp.Dial(connectStr)
	if err != nil {
		log.Println(err.Error())
		log.Fatal("Failed to connect to RabbitMQ")
	}
	defer conn.Close()

	ch, err := conn.Channel()
	//ErrorHanding(err, "Failed to open a channel")
	if err != nil {
		log.Println(err.Error())
		log.Fatal("Failed to open a channel")
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName,
		false, // возможность сохранения состояния при перезагрузке сервера. очередь хранится в RAM (random-access memory), поэтому, чтобы обеспечить её персистентное состояние, мы должны дополнительно сохранять очередь куда-то на диск
		false, // очищается ли очередь автоматически, когда соединение закрывается
		false, // может ли быть доступна очередь из других каналов
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Println(err.Error())
		log.Fatal("Failed to open a channel")
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Println(err.Error())
		log.Fatal("Failed to register a consume")
	}

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
	}
	/*бесконечныый цикл
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	select {}
	*/

}
