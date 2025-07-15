package order

import (
	"net/url"
	"time"

	it "github.com/renxzen/go-getters/examples/order/item"
)

//go:generate go-getters -structs=Order -output=getters.go

type Order struct {
	ID         int
	Name       string
	Items      []it.Item
	TotalPrice float64
	Image      url.URL
	Status     string
	CreatedAt  *time.Time
}
