package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/amtrindade/go-intensive/internal/infra/database"
	"github.com/amtrindade/go-intensive/internal/usecase"
	"github.com/amtrindade/go-intensive/pkg/rabbitmq"
	_ "github.com/mattn/go-sqlite3"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Car struct {
	Model string
	Color string
}

func (c Car) Start() {
	println(c.Model + " has been started!")
}

func main() {

	db, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		panic(err)
	}
	defer db.Close() // espera rodar e executa o close da conexão
	orderRepository := database.NewOrderRepository(db)
	uc := usecase.NewCalculateFinalPrice(orderRepository)

	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()
	msgRabbitmqChannel := make(chan amqp.Delivery)
	go rabbitmq.Consume(ch, msgRabbitmqChannel) // escutando a fila // trava e por isso T2
	rabbitmqWorker(msgRabbitmqChannel, uc)      //T1

	// input := usecase.OrderInput{
	// 	ID:    "125",
	// 	Price: 10.0,
	// 	Tax:   1.0,
	// }
	// output, err := uc.Execute(input)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(output)
}

func rabbitmqWorker(msgChan chan amqp.Delivery, uc *usecase.CalculateFinalPrice) {
	fmt.Println("Starting rabbitmq")
	for msg := range msgChan {
		var input usecase.OrderInput
		err := json.Unmarshal(msg.Body, &input)
		if err != nil {
			panic(err)
		}
		output, err := uc.Execute(input)
		if err != nil {
			panic(err)
		}
		msg.Ack(false)
		fmt.Println("Mensagem processada e salva no banco:", output)
	}
}
