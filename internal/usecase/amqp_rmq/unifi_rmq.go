package amqp_rmq

import (
	"github.com/streadway/amqp"
	"log"
)

type UnifiRmq struct {
	connectString string
	servExchange  string
}

func NewRmqUnifi(connectStr, servExchange string) *UnifiRmq {
	return &UnifiRmq{
		connectString: connectStr,
		servExchange:  servExchange,
	}

}

// https://russianblogs.com/article/53791654151/
func (ur *UnifiRmq) Publish(message, queueName string) error {
	// Подключаемся к серверу RabbitMQ
	//conn, err := lib.RabbitMQConn()
	conn, err := ur.RabbitMQConn()
	//lib.ErrorHanding(err, "Failed to connect to RabbitMQ")
	if err != nil {
		log.Println("Failed to connect to RabbitMQ")
		return err
	}

	// закрываем соединение
	defer conn.Close()

	// Создаем новый канал
	ch, err := conn.Channel()
	//ErrorHanding(err, "Failed to open a channel")
	if err != nil {
		log.Println("Failed to open a channel")
		return err
	}

	// Закрываем канал
	defer ch.Close()

	// Объявить или создать очередь для хранения сообщений
	q, err := ch.QueueDeclare(
		//"simple:queue", // Имя очереди
		queueName,
		false, // возможность сохранения состояния при перезагрузке сервера. очередь хранится в RAM (random-access memory), поэтому, чтобы обеспечить её персистентное состояние, мы должны дополнительно сохранять очередь куда-то на диск
		false, // очищается ли очередь автоматически, когда соединение закрывается
		false, // может ли быть доступна очередь из других каналов
		false, // no-wait
		nil,   // arguments
	)
	//lib.ErrorHanding(err, "Failed to declare a queue")
	if err != nil {
		log.Println("Failed to declare a queue")
		return err
	}

	//data := simpleDemo{		Name: "Tom",		Addr: "Beijing",	}
	//dataBytes, err := json.Marshal(data)
	//if err != nil {		lib.ErrorHanding(err, "struct to json failed")	}

	err = ch.Publish(
		ur.servExchange, // exchange
		q.Name,          // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			//DeliveryMode: amqp.Persistent,
			ContentType: "text/plain",
			//Body:        dataBytes,
			Body: []byte(message),
		})
	//log.Printf(" [x] Sent %s", dataBytes)
	//lib.ErrorHanding(err, "Failed to publish a message")
	if err != nil {
		log.Println("Failed to publish a message")
		return err
	}

	return nil
}

// Функция подключения RabbitMQ
func (ur *UnifiRmq) RabbitMQConn() (conn *amqp.Connection, err error) {
	// Создаем новое соединение
	conn, err = amqp.Dial(ur.connectString)
	// возвращаем соединение и ошибку
	return
}

// Функция обработки ошибок
func ErrorHanding(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
