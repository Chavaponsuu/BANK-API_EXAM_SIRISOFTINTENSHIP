package routes

import (
	"github.com/gin-gonic/gin"
	_ "github.com/krizad/go-gin-api/docs"
	"github.com/krizad/go-gin-api/internal/adapters/http/handlers"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(r *gin.Engine, eh *handlers.ExampleHandler, ah *handlers.AccountHandler, th *handlers.TransactionHandler) {
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
		transactions := api.Group("/transactions")
		{
			transactions.GET("", th.GetAllTransactionHistory)
		}

		accounts := api.Group("/accounts")
		{
			accounts.GET("", ah.GetAccountList)
			accounts.POST("", ah.CreateAccount)
			accounts.GET("/:account_number", ah.GetAccountByNumber)
			accounts.POST("/:account_number/deposit", th.Deposit)
			accounts.POST("/:account_number/withdraw", th.Withdraw)
			accounts.GET("/:account_number/transactions", th.GetTransactionHistory)

			accounts.PATCH("/:account_number/close", ah.CloseAccountHandler)
		}

	}

	r.GET("/health", eh.HealthCheck)
}
