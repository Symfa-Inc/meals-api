package dev

import (
	"fmt"
	"go_api/src/config"
	"go_api/src/domain"
	"go_api/src/repository"
	"strconv"
)

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
		cateringId := cateringResult.ID.String()

		categoryResult, _ := categoryRepo.GetByKey("name", "супы", cateringId)

		categoryId := categoryResult.ID.String()

		dishesArray, _, _ = dishRepo.Get(cateringId, categoryId)

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
