package delivery

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go_api/src/delivery/middleware"
	"go_api/src/types"
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
	category := usecase.NewCategory()
	dish := usecase.NewDish()
	image := usecase.NewImage()

	validator := middleware.NewValidator()

	configCors := cors.DefaultConfig()
	configCors.AllowOrigins = []string{os.Getenv("CLIENT_URL"), os.Getenv("CLIENT_MOBILE_URL")}

	configCors.AllowCredentials = true
	r.Use(cors.New(configCors))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	dir, _ := os.Getwd()
	r.Use(static.Serve("/static/", static.LocalFile(dir+"/src/static/images", true)))

	r.GET("/api-docs/static/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/refresh-token", middleware.Passport().RefreshHandler)
	r.POST("/login", middleware.Passport().LoginHandler)
	r.GET("/logout", middleware.Passport().LogoutHandler)
	authRequired := r.Group("/")
	authRequired.Use(middleware.Passport().MiddlewareFunc())
	{
		authRequired.GET("/is-authenticated", auth.IsAuthenticated)

		cateringGroup := authRequired.Group("/")
		cateringGroup.Use(validator.ValidateRoles(types.UserRoleEnum.CateringAdmin, types.UserRoleEnum.SuperAdmin))
		{
			cateringGroup.POST("/caterings", catering.Add)
			cateringGroup.GET("/caterings", catering.Get)
			cateringGroup.DELETE("/caterings/:id", catering.Delete)
			cateringGroup.PUT("/caterings/:id", catering.Update)

			cateringGroup.GET("/images", image.Get)

			cateringRoutes := cateringGroup.Group("/caterings")
			{
				cateringRoutes.POST("/:id/meals", meal.Add)
				cateringRoutes.GET("/:id/meals", meal.Get)
				cateringRoutes.PUT("/:id/meals/:mealId", meal.Update)

				cateringRoutes.POST("/:id/categories", category.Add)
				cateringRoutes.GET("/:id/categories", category.Get)
				cateringRoutes.DELETE("/:id/categories/:categoryId", category.Delete)
				cateringRoutes.PUT("/:id/categories/:categoryId", category.Update)

				cateringRoutes.POST("/:id/dishes", dish.Add)
				cateringRoutes.DELETE("/:id/dishes/:dishId", dish.Delete)
				cateringRoutes.GET("/:id/dishes", dish.Get)
				cateringRoutes.PUT("/:id/dishes/:dishId", dish.Update)

				cateringRoutes.POST("/:id/images", image.Add)
				cateringRoutes.DELETE("/:id/images/:imageId", image.Delete)
			}
		}
	}
	return r
}
