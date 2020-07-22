package types

// PathID struct for path binding
type PathID struct {
	ID string `uri:"id" json:"id" binding:"required"`
} //@name IDResponse

// PathMeal struct for path binding
type PathMeal struct {
	ID     string `uri:"id" json:"id" binding:"required"`
	MealID string `uri:"mealId" json:"mealId" binding:"required"`
} //@name MealPathResponse

// PathCategory struct for path binding
type PathCategory struct {
	ID         string `uri:"id" json:"id" binding:"required"`
	CategoryID string `uri:"categoryID" json:"categoryID" binding:"required"`
}

// PathDish struct for path binding
type PathDish struct {
	CateringID string `uri:"id" json:"id" binding:"required"`
	DishID     string `uri:"dishId" json:"dishId" binding:"required"`
}

// PathDishGet struct for path binding
type PathDishGet struct {
	CateringID string `uri:"id" json:"id" binding:"required"`
	CategoryID string `uri:"categoryID" json:"categoryID" binding:"required"`
}

// PathImageDish struct for path binding
type PathImageDish struct {
	CateringID string `uri:"id" json:"id" binding:"required"`
	ImageID    string `uri:"imageId" json:"imageId" binding:"required"`
	DishID     string `uri:"dishId" json:"dishId" binding:"required"`
}
