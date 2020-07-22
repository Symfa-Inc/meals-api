package dev

import (
	"fmt"
	"go_api/src/config"
	"go_api/src/domain"
	"io/ioutil"
	"os"
)

// CreateImages creates seeds for images table
func CreateImages() {
	seedExist := config.DB.Where("name = ?", "init images").First(&domain.Seed{}).Error
	if seedExist != nil {
		seed := domain.Seed{
			Name: "init images",
		}

		dir, _ := os.Getwd()
		salads, _ := ioutil.ReadDir(dir + "/src/static/images/salad")
		for _, salad := range salads {
			categoryName := "salad"
			image := domain.Image{
				Path:     "/salad/" + salad.Name(),
				Category: &categoryName,
			}
			config.DB.Create(&image)
		}

		desserts, _ := ioutil.ReadDir(dir + "/src/static/images/dessert")
		for _, dessert := range desserts {
			categoryName := "dessert"
			image := domain.Image{
				Path:     "/dessert/" + dessert.Name(),
				Category: &categoryName,
			}
			config.DB.Create(&image)
		}

		garnishes, _ := ioutil.ReadDir(dir + "/src/static/images/garnish")
		for _, garnish := range garnishes {
			categoryName := "garnish"
			image := domain.Image{
				Path:     "/garnish/" + garnish.Name(),
				Category: &categoryName,
			}
			config.DB.Create(&image)
		}

		soups, _ := ioutil.ReadDir(dir + "/src/static/images/soup")
		for _, soup := range soups {
			categoryName := "soup"
			image := domain.Image{
				Path:     "/soup/" + soup.Name(),
				Category: &categoryName,
			}
			config.DB.Create(&image)
		}

		firstCourses, _ := ioutil.ReadDir(dir + "/src/static/images/first_course")
		for _, firstCourse := range firstCourses {
			categoryName := "first_course"
			image := domain.Image{
				Path:     "/first_course/" + firstCourse.Name(),
				Category: &categoryName,
			}
			config.DB.Create(&image)
		}

		config.DB.Create(&seed)
		fmt.Println("=== Images seeds created ===")
	} else {
		fmt.Printf("Seed `init images` already exists \n")
	}
}
