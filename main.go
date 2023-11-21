package main

import (
	"database/sql"
	"fmt"

	"github.com/amtrindade/go-instesive/internal/infra/database"
	"github.com/amtrindade/go-instesive/internal/usecase"
	_ "github.com/mattn/go-sqlite3"
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
	orderRepository := database.NewOrderRepository(db)
	uc := usecase.NewCalculateFinalPrice(orderRepository)

	input := usecase.OrderInput{
		ID:    "124",
		Price: 10.0,
		Tax:   1.0,
	}
	output, err := uc.Execute(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(output)
}
