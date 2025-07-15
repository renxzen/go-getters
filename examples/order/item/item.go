package item

//go:generate go-getters -structs=Item -output=getters.go

type Item struct {
	Name   *string
	Price  *float64
	Count  *int64
	Active *bool
}
