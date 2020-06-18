package types

type PathId struct {
	ID string `uri:"id" json:"id" binding:"required"`
} //@name IDResponse

type PathMeal struct {
	ID     string `uri:"id" json:"id" binding:"required"`
	MealID string `uri:"mealId" json:"mealId" binding:"required"`
} //@name MealPathResponse

type PathDishCategory struct {
	ID         string `uri:"id" json:"id" binding:"required"`
	CategoryID string `uri:"categoryId" json:"categoryId" binding:"required"`
}
