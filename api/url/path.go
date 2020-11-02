package url

// PathID struct for path binding
type PathID struct {
	ID string `uri:"id" json:"id" binding:"required"`
} //@name IDResponse

// PathMeal struct for path binding
type PathMeal struct {
	ID     string `uri:"id" json:"id" binding:"required"`
	MealID string `uri:"mealId" json:"mealId" binding:"required"`
} //@name MealPathResponse

// PathUser struct for path binding
type PathUser struct {
	ID     string `uri:"id" json:"id" binding:"required"`
	UserID string `uri:"userId" json:"userId" binding:"required"`
}

// PathCategory struct for path binding
type PathCategory struct {
	ID         string `uri:"id" json:"id" binding:"required"`
	ClientID   string `uri:"clientId" json:"clientId" binding:"required"`
	CategoryID string `uri:"categoryID" json:"categoryID" binding:"required"`
}

// PathClient struct for path binding
type PathClient struct {
	ID       string `uri:"id" json:"id" binding:"required"`
	ClientID string `uri:"clientId" json:"clientId" binding:"required"`
}

// PathUpdateMeal struct for path binding
type PathUpdateMeal struct {
	ID       string `uri:"id" json:"id" binding:"required"`
	ClientID string `uri:"clientId" json:"clientId" binding:"required"`
	MealID   string `uri:"mealId" json:"mealId" binding:"required"`
}

// PathDish struct for path binding
type PathDish struct {
	CateringID string `uri:"id" json:"id" binding:"required"`
	DishID     string `uri:"dishId" json:"dishId" binding:"required"`
}

// PathDishID struct for path binding
type PathDishID struct {
	ID     string `uri:"id" json:"id" binding:"required"`
	DishID string `uri:"dishId" json:"dishId" binding:"required"`
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

// PathSchedule struct for path binding
type PathSchedule struct {
	ID         string `uri:"id" json:"id" binding:"required"`
	ScheduleID string `uri:"scheduleId" json:"scheduleId" binding:"required"`
}

// PathAddress struct for path binding
type PathAddress struct {
	ID        string `uri:"id" json:"id" binding:"required"`
	AddressID string `uri:"addressId" json:"addressId" binding:"required"`
}

// PathOrder struct for path binding
type PathOrder struct {
	ID      string `uri:"id" json:"id" binding:"required"`
	OrderID string `uri:"orderId" json:"orderId" binding:"required"`
}
