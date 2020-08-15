package dev

import (
	"fmt"
	"github.com/Aiscom-LLC/meals-api/src/config"
	"github.com/Aiscom-LLC/meals-api/src/domain"
	"github.com/Aiscom-LLC/meals-api/src/repository"
	"strconv"
)

// CreateImageDishes creates seeds for image_dishes table
func CreateImageDishes() {
	cateringRepo := repository.NewCateringRepo()
	categoryRepo := repository.NewCategoryRepo()
	dishRepo := repository.NewDishRepo()
	imageRepo := repository.NewImageRepo()

	seedExists := config.DB.
		Where("name = ?", "init image dishes").
		First(&domain.Seed{}).Error
	if seedExists != nil {
		seed := domain.Seed{
			Name: "init image dishes",
		}

		var dishesArray []domain.Dish

		cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
		cateringID := cateringResult.ID.String()

		categoryResult, _ := categoryRepo.GetByKey("name", "супы", cateringID)

		categoryID := categoryResult.ID.String()

		dishesArray, _, _ = dishRepo.Get(cateringID, categoryID)

		var imageDish domain.ImageDish
		for i := range dishesArray {
			for j := 0; j < 3; j++ {
				imageResult, _ := imageRepo.GetByKey("path", "/salad/"+strconv.Itoa(j+1)+".png")
				imageDish.DishID = dishesArray[i].ID
				imageDish.ImageID = imageResult.ID
				config.DB.Create(&imageDish)
			}
		}

		config.DB.Create(&seed)
		fmt.Println("=== Image Dishes seeds created ===")
	} else {
		fmt.Printf("Seed `init image dishes` already exists \n")
	}
}
