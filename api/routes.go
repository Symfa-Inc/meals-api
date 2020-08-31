package api

import (
	"os"

	"github.com/Aiscom-LLC/meals-api/api/middleware"
	"github.com/Aiscom-LLC/meals-api/types"
	"github.com/Aiscom-LLC/meals-api/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//SetupRouter setting up gin router and config
func SetupRouter() *gin.Engine {
	r := gin.Default()

	auth := NewAuth()
	cateringUser := usecase.NewCateringUser()
	clientUser := usecase.NewClientUser()
	catering := usecase.NewCatering()
	cateringSchedule := usecase.NewCateringSchedule()
	clientSchedule := usecase.NewClientSchedule()
	client := usecase.NewClient()
	meal := usecase.NewMeal()
	category := usecase.NewCategory()
	dish := usecase.NewDish()
	image := usecase.NewImage()
	order :=  NewOrder()
	//order := usecase.NewOrder()
	address := usecase.NewAddress()

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
	r.GET("/is-authenticated", auth.IsAuthenticated)
	r.POST("/login", middleware.Passport().LoginHandler)
	r.GET("/logout", middleware.Passport().LogoutHandler)

	authRequired := r.Group("/")
	authRequired.Use(middleware.Passport().MiddlewareFunc())
	{

		suAdmin := authRequired.Group("/")
		suAdmin.Use(validator.ValidateRoles(
			types.UserRoleEnum.SuperAdmin,
		))
		{
			// caterings
			suAdmin.POST("/caterings", catering.Add)
			suAdmin.DELETE("/caterings/:id", catering.Delete)
		}

		caAdminSuAdmin := authRequired.Group("/")
		caAdminSuAdmin.Use(validator.ValidateRoles(
			types.UserRoleEnum.SuperAdmin,
			types.UserRoleEnum.CateringAdmin,
		))
		{
			// catering info
			caAdminSuAdmin.PUT("/caterings/:id", catering.Update)
			caAdminSuAdmin.GET("/caterings/:id", catering.GetByID)

			// catering users
			caAdminSuAdmin.PUT("/caterings/:id/users/:userId", cateringUser.Update)
			caAdminSuAdmin.POST("/caterings/:id/users", cateringUser.Add)
			caAdminSuAdmin.GET("/caterings/:id/users", cateringUser.Get)
			caAdminSuAdmin.DELETE("/caterings/:id/users/:userId", cateringUser.Delete)

			// catering categories
			caAdminSuAdmin.POST("/caterings/:id/clients/:clientId/categories", category.Add)
			caAdminSuAdmin.DELETE("/caterings/:id/clients/:clientId/categories/:categoryID", category.Delete)
			caAdminSuAdmin.PUT("/caterings/:id/clients/:clientId/categories/:categoryID", category.Update)

			// catering clients
			caAdminSuAdmin.POST("/caterings/:id/clients", client.Add)
			caAdminSuAdmin.GET("/caterings/:id/clients-orders", client.GetCateringClientsOrders)

			// clients
			caAdminSuAdmin.PUT("/clients/:id", client.Update)
			caAdminSuAdmin.DELETE("/clients/:id", client.Delete)

			// catering client orders
			//caAdminSuAdmin.GET("/caterings/:id/clients/:clientId/orders", order.GetCateringClientOrders)

			// catering dishes
			caAdminSuAdmin.POST("/caterings/:id/dishes", dish.Add)
			caAdminSuAdmin.DELETE("/caterings/:id/dishes/:dishId", dish.Delete)
			caAdminSuAdmin.PUT("/caterings/:id/dishes/:dishId", dish.Update)

			// catering images
			caAdminSuAdmin.GET("/images", image.Get)
			caAdminSuAdmin.POST("/caterings/:id/dishes/:dishId/images", image.Add)
			caAdminSuAdmin.DELETE("/caterings/:id/dishes/:dishId/images/:imageId", image.Delete)
			caAdminSuAdmin.PUT("/caterings/:id/dishes/:dishId/images/:imageId", image.Update)

			// catering meals
			caAdminSuAdmin.POST("/caterings/:id/clients/:clientId/meals", meal.Add)

			// catering schedules
			caAdminSuAdmin.PUT("/caterings/:id/schedules/:scheduleId", cateringSchedule.Update)

		}

		clAdminSuAdmin := authRequired.Group("/")
		clAdminSuAdmin.Use(validator.ValidateRoles(
			types.UserRoleEnum.ClientAdmin,
			types.UserRoleEnum.SuperAdmin,
		))
		{
			// client schedules
			clAdminSuAdmin.PUT("/clients/:id/schedules/:scheduleId", clientSchedule.Update)

			// client addresses
			clAdminSuAdmin.POST("/clients/:id/addresses", address.Add)
			clAdminSuAdmin.DELETE("/clients/:id/addresses/:addressId", address.Delete)
			clAdminSuAdmin.PUT("/clients/:id/addresses/:addressId", address.Update)

			// client orders
			//clAdminSuAdmin.GET("/clients/:id/orders", order.GetClientOrders)
			//clAdminSuAdmin.PUT("/clients/:id/orders", order.ApproveOrders)

			clAdminSuAdmin.PUT("/clients/:id/auto-approve", client.UpdateAutoApprove)
		}

		clAdminUser := authRequired.Group("/")
		clAdminUser.Use(validator.ValidateRoles(
			types.UserRoleEnum.SuperAdmin,
			types.UserRoleEnum.ClientAdmin,
			types.UserRoleEnum.User,
		))
		{
			// client orders
			clAdminUser.POST("/users/:id/orders", order.Add)
			//clAdminUser.DELETE("/users/:id/orders/:orderId", order.CancelOrder)
			//clAdminUser.GET("/users/:id/orders", order.GetUserOrder)

			//clAdminUser.GET("/clients/:id/order-status", order.GetOrderStatus)
		}

		allUsers := authRequired.Group("/")
		allUsers.Use(validator.ValidateRoles(
			types.UserRoleEnum.SuperAdmin,
			types.UserRoleEnum.CateringAdmin,
			types.UserRoleEnum.ClientAdmin,
			types.UserRoleEnum.User,
		))
		{
			// categories
			allUsers.GET("/caterings/:id/clients/:clientId/categories", category.Get)

			// catering meals
			allUsers.GET("/caterings/:id/clients/:clientId/meals", meal.Get)

			// schedules
			allUsers.GET("/caterings/:id/schedules", cateringSchedule.Get)
			allUsers.GET("/clients/:id/schedules", clientSchedule.Get)

			// dishes
			allUsers.GET("/caterings/:id/dishes", dish.Get)
			allUsers.GET("/caterings/:id/dishes/:dishId", dish.GetByID)
		}

		allAdmins := authRequired.Group("/")
		allAdmins.Use(validator.ValidateRoles(
			types.UserRoleEnum.SuperAdmin,
			types.UserRoleEnum.CateringAdmin,
			types.UserRoleEnum.ClientAdmin,
		))
		{
			// clients
			allAdmins.GET("/clients", client.Get)
			allAdmins.GET("/caterings/:id/clients", client.GetByCateringID)
			allAdmins.GET("/clients/:id", client.GetByID)

			// caterings
			allAdmins.GET("/caterings", catering.Get)

			// clients users
			allAdmins.DELETE("/clients/:id/users/:userId", clientUser.Delete)
			allAdmins.PUT("/clients/:id/users/:userId", clientUser.Update)
			allAdmins.POST("/clients/:id/users", clientUser.Add)
			allAdmins.GET("/clients/:id/users", clientUser.Get)

			// client addresses
			allAdmins.GET("/clients/:id/addresses", address.Get)
		}
	}
	return r
}
