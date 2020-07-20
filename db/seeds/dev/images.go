package dev

import (
	"fmt"
	"go_api/src/config"
	"go_api/src/domain"
	"io/ioutil"
	"os"
)

func CreateImages() {
	seedExist := config.DB.Where("name = ?", "init images").First(&domain.Seed{}).Error
	if seedExist != nil {
		seed := domain.Seed{
			Name: "init images",
		}

		dir, _ := os.Getwd()
		salads, _ := ioutil.ReadDir(dir + "/src/static/images/salad")
		for _, salad := range salads {
			image := domain.Image{
				Path:     "/salad/" + salad.Name(),
				Category: "salad",
			}
			config.DB.Create(&image)
		}

		desserts, _ := ioutil.ReadDir(dir + "/src/static/images/dessert")
		for _, dessert := range desserts {
			image := domain.Image{
				Path:     "/dessert/" + dessert.Name(),
				Category: "dessert",
			}
			config.DB.Create(&image)
		}

		garnishes, _ := ioutil.ReadDir(dir + "/src/static/images/garnish")
		for _, garnish := range garnishes {
			image := domain.Image{
				Path:     "/garnish/" + garnish.Name(),
				Category: "garnish",
			}
			config.DB.Create(&image)
		}

		soups, _ := ioutil.ReadDir(dir + "/src/static/images/soup")
		for _, soup := range soups {
			image := domain.Image{
				Path:     "/soup/" + soup.Name(),
				Category: "soup",
			}
			config.DB.Create(&image)
		}

		firstCourses, _ := ioutil.ReadDir(dir + "/src/static/images/first_course")
		for _, firstCourse := range firstCourses {
			image := domain.Image{
				Path:     "/first_course/" + firstCourse.Name(),
				Category: "first_course",
			}
			config.DB.Create(&image)
		}

		config.DB.Create(&seed)
		fmt.Println("=== Images seeds created ===")
	} else {
		fmt.Printf("Seed `init images` already exists \n")
	}
}
