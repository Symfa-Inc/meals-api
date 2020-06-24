package delivery

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go_api/src/delivery/middleware"
	"go_api/src/usecase"
	"net/http"
	"os"
)

// RedirectFunc wrapper for a Gin Redirect function
// which takes a route as a string and returns original Gin Redirect func
func RedirectFunc(route string) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, route)
		c.Abort()
	}
}

//SetupRouter setting up gin router and config
func SetupRouter() *gin.Engine {
	r := gin.Default()

	auth := usecase.NewAuth()
	catering := usecase.NewCatering()
	meal := usecase.NewMeal()
	dishCategory := usecase.NewDishCategory()
	dish := usecase.NewDish()

	configCors := cors.DefaultConfig()
	configCors.AllowOrigins = []string{os.Getenv("CLIENT_URL")}
	configCors.AllowCredentials = true
	r.Use(cors.New(configCors))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/api-docs/static/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/refresh-token", middleware.Passport().RefreshHandler)
	r.POST("/login", middleware.Passport().LoginHandler)
	r.GET("/logout", middleware.Passport().LogoutHandler)
	authRequired := r.Group("/")
	authRequired.Use(middleware.Passport().MiddlewareFunc())
	{
		authRequired.GET("/is-authenticated", auth.IsAuthenticated)

		cateringGroup := authRequired.Group("/")
		{
			cateringGroup.POST("/caterings", catering.Add)
			cateringGroup.GET("/caterings", catering.Get)
			cateringGroup.DELETE("/caterings/:id", catering.Delete)
			cateringGroup.PUT("/caterings/:id", catering.Update)

			cateringRoutes := cateringGroup.Group("/caterings")
			{
				cateringRoutes.POST("/:id/meals", meal.Add)
				cateringRoutes.GET("/:id/meals", meal.Get)
				cateringRoutes.PUT("/:id/meals/:mealId", meal.Update)

				cateringRoutes.POST("/:id/dish-categories", dishCategory.Add)
				cateringRoutes.GET("/:id/dish-categories", dishCategory.Get)
				cateringRoutes.DELETE("/:id/dish-categories/:categoryId", dishCategory.Delete)
				cateringRoutes.PUT("/:id/dish-categories/:categoryId", dishCategory.Update)

				cateringRoutes.POST("/:id/dish", dish.Add)
			}
		}
	}
	return r
}
