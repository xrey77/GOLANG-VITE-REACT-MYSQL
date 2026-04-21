package dto

type Category struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	Name     string    `json:"category"`
	Products []Product `json:"products" gorm:"foreignKey:CategoryID"`
}

type Product struct {
	ID           uint    `json:"id" gorm:"primaryKey"`
	Descriptions string  `json:"descriptions"`
	Qty          int     `json:"qty"`
	Unit         string  `json:"unit"`
	Costprice    float64 `json:"costprice"`
	Sellprice    float64 `json:"sellprice"`
	CategoryID   uint    `json:"category_id"`
}
