package mux

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go_api/src/mux/auth"
	"go_api/src/mux/catering"
	"go_api/src/mux/middleware"
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

	configCors := cors.DefaultConfig()
	configCors.AllowOrigins = []string{os.Getenv("CLIENT_URL")}
	configCors.AllowCredentials = true;
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
			cateringGroup.POST("/caterings", catering.AddCatering)
			cateringGroup.GET("/caterings", catering.GetCaterings)
			cateringGroup.GET("/caterings/:id", catering.GetCatering)
			cateringGroup.DELETE("/caterings/:id", catering.DeleteCatering)
			cateringGroup.PUT("/caterings/:id", catering.UpdateCatering)
		}
	}
	return r
}
