package main

import (
	"log"

	"github.com/Aiscom-LLC/meals-api/src/domain"

	"github.com/Aiscom-LLC/meals-api/src/config"
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func main() {
	m := gormigrate.New(config.DB, gormigrate.DefaultOptions, []*gormigrate.Migration{})
	//{
	//		ID: "16",
	//		Migrate: func(tx *gorm.DB) error {
	//
	//			// it's a good pratice to copy the struct inside the function,
	//			// so side effects are prevented if the original struct changes during the time
	//			type UserOrders struct {
	//				domain.UserOrders
	//			}
	//			return tx.AutoMigrate(&UserOrders{}).Error
	//		},
	//		Rollback: func(tx *gorm.DB) error {
	//
	//			return tx.DropTable("user_orders").Error
	//		},
	//	},
	//})

	m.InitSchema(func(tx *gorm.DB) error {
		err := tx.AutoMigrate(
			&domain.Seed{},
			&domain.User{},
			&domain.Catering{},
			&domain.CateringUser{},
			&domain.CateringSchedule{},
			&domain.Client{},
			&domain.ClientUser{},
			&domain.Address{},
			&domain.ClientSchedule{},
			&domain.Meal{},
			&domain.Category{},
			&domain.Dish{},
			&domain.ImageDish{},
			&domain.Image{},
			&domain.MealDish{},
			&domain.Order{},
			&domain.OrderDishes{},
			&domain.UserOrders{},
		)
		if err != nil {
			return err.Error
		}
		if err := tx.Exec("ALTER TABLE catering_users ADD CONSTRAINT fk_caterings_users FOREIGN KEY (catering_id) REFERENCES caterings (id)").Error; err != nil {
			return err
		}
		return nil
	})

	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v\n", err)
	}
	log.Println("=== ADD MIGRATIONS ===")
}
