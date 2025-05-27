package domain

type Product struct {
	ID         int  `bson:"_id,omitempty" json:"id"`
	Name       string  `bson:"name" json:"name"`
	Price      float32 `bson:"price" json:"price"`
	Stock      int32    `bson:"stock" json:"stock"`
	CategoryID int32  `bson:"category_id" json:"category_id"`
}

