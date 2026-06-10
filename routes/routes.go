package routes

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/krizad/go-gin-api/handlers"
)

func Setup(r *gin.Engine, eh *handlers.ExampleHandler, ah *handlers.AccountHandler) {
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		examples := api.Group("/examples")
		{
			examples.GET("", eh.ListExamples)
			examples.POST("", eh.CreateExample)

			examples.GET("/search", eh.GetExampleByEmail)

			examples.GET("/:id", eh.GetExample)
			examples.PUT("/:id", eh.UpdateExample)
			examples.PATCH("/:id", eh.PatchExample)
			examples.DELETE("/:id", eh.DeleteExample)
		}

		accounts := api.Group("/accounts")
		{
			accounts.GET("", ah.GetAccountList)
			accounts.POST("", ah.CreateAccount)
			accounts.GET("/:account_number", ah.GetAccountByNumber)
		}
	}

	r.GET("/health", eh.HealthCheck)
}
