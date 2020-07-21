package types

type PathId struct {
	ID string `uri:"id" json:"id" binding:"required"`
} //@name IDResponse

type PathMeal struct {
	ID     string `uri:"id" json:"id" binding:"required"`
	MealId string `uri:"mealId" json:"mealId" binding:"required"`
} //@name MealPathResponse

type PathCategory struct {
	ID         string `uri:"id" json:"id" binding:"required"`
	CategoryID string `uri:"categoryId" json:"categoryId" binding:"required"`
}

type PathDish struct {
	CateringID string `uri:"id" json:"id" binding:"required"`
	DishID     string `uri:"dishId" json:"dishId" binding:"required"`
}

type PathDishGet struct {
	CateringID string `uri:"id" json:"id" binding:"required"`
	CategoryID string `uri:"categoryId" json:"categoryId" binding:"required"`
}

type PathImageDish struct {
	CateringID string `uri:"id" json:"id" binding:"required"`
	ImageId    string `uri:"imageId" json:"imageId" binding:"required"`
	DishId     string `uri:"dishId" json:"dishId" binding:"required"`
}
