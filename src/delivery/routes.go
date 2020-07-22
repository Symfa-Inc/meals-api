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
	cateringSchedule := usecase.NewCateringSchedule()
	clientSchedule := usecase.NewClientSchedule()
	client := usecase.NewClient()
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

		clientGroup := authRequired.Group("/")
		clientGroup.Use(validator.ValidateRoles(
			types.UserRoleEnum.ClientAdmin,
			types.UserRoleEnum.SuperAdmin,
		))
		{
			clientGroup.GET("/clients/:id/schedules", clientSchedule.Get)
			clientGroup.PUT("/clients/:id/schedules/:scheduleId", clientSchedule.Update)
		}

		cateringGroup := authRequired.Group("/")
		cateringGroup.Use(validator.ValidateRoles(
			types.UserRoleEnum.CateringAdmin,
			types.UserRoleEnum.SuperAdmin,
		))
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

				cateringRoutes.POST("/:id/clients", client.Add)
				cateringRoutes.GET("/:id/clients", client.Get)
				cateringRoutes.DELETE("/:id/clients/:clientId", client.Delete)
				cateringRoutes.PUT("/:id/clients/:clientId", client.Update)

				cateringRoutes.GET("/:id/schedules", cateringSchedule.Get)
				cateringRoutes.PUT("/:id/schedules/:scheduleId", cateringSchedule.Update)

				cateringRoutes.POST("/:id/categories", category.Add)
				cateringRoutes.GET("/:id/categories", category.Get)
				cateringRoutes.DELETE("/:id/categories/:categoryID", category.Delete)
				cateringRoutes.PUT("/:id/categories/:categoryID", category.Update)

				cateringRoutes.POST("/:id/dishes", dish.Add)
				cateringRoutes.DELETE("/:id/dishes/:dishId", dish.Delete)
				cateringRoutes.GET("/:id/dishes", dish.Get)
				cateringRoutes.GET("/:id/dishes/:dishId", dish.GetByID)
				cateringRoutes.PUT("/:id/dishes/:dishId", dish.Update)

				cateringRoutes.POST("/:id/images", image.Add)
				cateringRoutes.DELETE("/:id/images/:imageId/dish/:dishId", image.Delete)
				cateringRoutes.PUT("/:id/images/:imageId/dish/:dishId", image.Update)
			}
		}
	}
	return r
}
