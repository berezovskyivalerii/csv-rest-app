package domain

type Product struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

type ProductUpdate struct{
	Name  string `json:"name"`
	Price int64  `json:"price"`
}