package domain

type Order struct {
	ID     int       `bson:"_id,omitempty" json:"id"`
	UserID int       `bson:"user_id" json:"user_id"`
	Status string       `bson:"status" json:"status"`
	Items  []OrderItem  `bson:"items" json:"items"`
}

type OrderItem struct {
	ProductID int	`bson:"product_id" json:"product_id"`
	Quantity  int    `bson:"quantity" json:"quantity"`
}