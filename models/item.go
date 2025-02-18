package models

type Item struct {
	Id    int     `json:"id"`
	Name  string  `json:"name"`
	Count int     `json:"count"`
	Price float64 `json:"price"`
}

type ItemEdit struct {
	Name  string  `json:"name"`
	Count int     `json:"count"`
	Price float64 `json:"price"`
}
