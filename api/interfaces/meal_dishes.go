package domain

import (
	"github.com/Aiscom-LLC/meals-api/domain"
)

// MealDishRepository is mealDish interface for repository
type MealDishRepository interface {
	Add(mealDish domain.MealDish) error
	Delete(mealID string) error
}
