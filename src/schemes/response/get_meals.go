package response

type GetMealsModel struct {
	DishStruct
	CategoryID   string `gorm:"column:category_id" json:"categoryId"`
	CategoryName string `gorm:"column:category_name" json:"categoryName"`
}

type DishStruct struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Weight       string `json:"weight,omitempty"`
	Price        string `json:"price,omitempty"`
	Images       string `json:"images,omitempty"`
	Desc         string `json:"desc,omitempty"`
}

type GetMealsResponse struct {
	CategoryName []DishStruct `json:"categoryName"`
}
