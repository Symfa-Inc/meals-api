package api

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Aiscom-LLC/meals-api/api/middleware"
	"github.com/Aiscom-LLC/meals-api/repository/enums"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//SetupRouter setting up gin router and config
func SetupRouter() *gin.Engine {
	r := gin.Default()

	if err := os.Mkdir("logs", 0777); err != nil {
		fmt.Println(err)
	}
	file, _ := os.Create("logs/gin: " + time.Now().UTC().String() + ".log")
	fileErr, _ := os.Create("logs/err: " + time.Now().UTC().String() + ".log")
	gin.DefaultWriter = io.MultiWriter(file)
	gin.DefaultErrorWriter = io.MultiWriter(fileErr)

	auth := NewAuth()
	cateringUser := NewCateringUser()
	clientUser := NewClientUser()
	catering := NewCatering()
	cateringSchedule := NewCateringSchedule()
	clientSchedule := NewClientSchedule()
	client := NewClient()
	meal := NewMeal()
	category := NewCategory()
	dish := NewDish()
	image := NewImage()
	order := NewOrder()
	address := NewAddress()

	validator := middleware.NewValidator()

	configCors := cors.DefaultConfig()
	configCors.AllowOrigins = []string{os.Getenv("CLIENT_URL"), os.Getenv("CLIENT_MOBILE_URL"), os.Getenv("TASTY_URL"), os.Getenv("TASTY_MOBILE_URL")}

	configCors.AllowCredentials = true
	r.Use(cors.New(configCors))

	dir, _ := os.Getwd()
	r.Use(static.Serve("/static/", static.LocalFile(dir+"/static/images", true)))

	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		return fmt.Sprintf("IP=%s - [Date=%s] Method=\"%s\" Path=%s Request Prototype=%s Status code=%d Latency=%s User agent=\"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	r.GET("/api-docs/static/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/is-authenticated", auth.IsAuthenticated)
	r.POST("/login", middleware.Passport().LoginHandler)
	r.GET("/logout", middleware.Passport().LogoutHandler)
	r.POST("/recovery-password", auth.RecoveryPassword)

	authRequired := r.Group("/")
	authRequired.Use(middleware.Passport().MiddlewareFunc())
	{

		suAdmin := authRequired.Group("/")
		suAdmin.Use(validator.ValidateRoles(
			enums.UserRoleEnum.SuperAdmin,
		))
		{
			// caterings
			suAdmin.POST("/caterings", catering.Add)
			suAdmin.DELETE("/caterings/:id", catering.Delete)
		}

		caAdminSuAdmin := authRequired.Group("/")
		caAdminSuAdmin.Use(validator.ValidateRoles(
			enums.UserRoleEnum.SuperAdmin,
			enums.UserRoleEnum.CateringAdmin,
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
			caAdminSuAdmin.GET("/caterings/:id/clients/:clientId/orders", order.GetCateringClientOrders)

			// catering dishes
			caAdminSuAdmin.POST("/caterings/:id/dishes", dish.Add)
			caAdminSuAdmin.DELETE("/caterings/:id/dishes/:dishId", dish.Delete)
			caAdminSuAdmin.PUT("/caterings/:id/dishes/:dishId", dish.Update)
			caAdminSuAdmin.GET("/caterings/:id/catering-dishes", dish.GetCateringDish)
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
			enums.UserRoleEnum.ClientAdmin,
			enums.UserRoleEnum.SuperAdmin,
		))
		{
			// client schedules
			clAdminSuAdmin.PUT("/clients/:id/schedules/:scheduleId", clientSchedule.Update)

			// client addresses
			clAdminSuAdmin.POST("/clients/:id/addresses", address.Add)
			clAdminSuAdmin.DELETE("/clients/:id/addresses/:addressId", address.Delete)
			clAdminSuAdmin.PUT("/clients/:id/addresses/:addressId", address.Update)

			// client orders
			clAdminSuAdmin.GET("/clients/:id/orders", order.GetClientOrders)
			clAdminSuAdmin.PUT("/clients/:id/orders", order.ApproveOrders)

			clAdminSuAdmin.PUT("/clients/:id/auto-approve", client.UpdateAutoApprove)
		}

		clAdminUser := authRequired.Group("/")
		clAdminUser.Use(validator.ValidateRoles(
			enums.UserRoleEnum.SuperAdmin,
			enums.UserRoleEnum.ClientAdmin,
			enums.UserRoleEnum.User,
		))
		{
			// client orders
			clAdminUser.POST("/users/:id/orders", order.Add)
			clAdminUser.DELETE("/users/:id/orders/:orderId", order.CancelOrder)
			clAdminUser.GET("/users/:id/orders", order.GetUserOrder)

			clAdminUser.GET("/clients/:id/order-status", order.GetOrderStatus)
		}

		allUsers := authRequired.Group("/")
		allUsers.Use(validator.ValidateRoles(
			enums.UserRoleEnum.SuperAdmin,
			enums.UserRoleEnum.CateringAdmin,
			enums.UserRoleEnum.ClientAdmin,
			enums.UserRoleEnum.User,
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

			// auth
			allUsers.PUT("/auth/change-password", auth.ChangePassword)
		}

		allAdmins := authRequired.Group("/")
		allAdmins.Use(validator.ValidateRoles(
			enums.UserRoleEnum.SuperAdmin,
			enums.UserRoleEnum.CateringAdmin,
			enums.UserRoleEnum.ClientAdmin,
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

			// orders
			allAdmins.GET("/clients/:id/orders-file", order.GetClientOrdersExcel)
		}
	}
	return r
}
