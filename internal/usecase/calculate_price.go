package usecase

import "github.com/amtrindade/go-instesive/internal/entity"

type OrderInput struct {
	ID    string
	Price float64
	Tax   float64
}

type OrderOutPut struct {
	ID         string
	Price      float64
	Tax        float64
	FinalPrice float64
}

// SOLID - "D" - Dependency Inversion Principle
type CalculateFinalPrice struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewCalculateFinalPrice(orderRepository entity.OrderRepositoryInterface) *CalculateFinalPrice {
	return &CalculateFinalPrice{
		OrderRepository: orderRepository,
	}
}

func (c *CalculateFinalPrice) Execute(input OrderInput) (*OrderOutPut, error) {
	order, err := entity.NewOrder(input.ID, input.Price, input.Tax)
	if err != nil {
		return nil, err
	}
	err = order.CalculateFinalPrice()
	if err != nil {
		return nil, err
	}
	err = c.OrderRepository.Save(order)
	if err != nil {
		return nil, err
	}
	return &OrderOutPut{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}, nil
}
