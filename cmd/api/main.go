// @title           Billing Service API
// @version         1.0
// @description     Subscription billing service for managing plans, subscriptions, and payment methods.
// @host            localhost:8080
// @BasePath        /

// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and your JWT token
package main

import (
	"billingService/backend/internal/accounts"
	"billingService/backend/internal/invoice"
	"billingService/backend/internal/middleware"
	"billingService/backend/internal/plans"
	"billingService/backend/internal/subscription"
	"log"
	"os"

	_ "billingService/backend/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading environment variables.")
	}

	// Translate DB-level errors (like unique constraint violations) into typed errors like gorm.ErrDuplicatedKey
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{TranslateError: true})
	if err != nil {
		log.Println(err)
	}

	route := gin.Default()

	// CORS — allow requests from any local origin (frontend dev server)
	route.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Swagger UI endpoint
	// https://github.com/swaggo/swag/blob/master/README.md#declarative-comments-format
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	plansRepository := plans.NewRepository(db)
	plansProvider := plans.NewProvider(plansRepository)
	plansGroup := route.Group("/api/v1/plans")
	{
		plansGroup.GET("/", plansProvider.GetPlans)
		plansGroup.GET("/:id", plansProvider.GetPlanByID)
		plansGroup.POST("/create/", plansProvider.CreatePlan)
		plansGroup.PATCH("/update/:id", plansProvider.UpdatePlanByID)
		plansGroup.PATCH("/deprecate/:id", plansProvider.DeprecatePlanByID)
	}

	accountsRepository := accounts.NewAccountsRepository(db)
	accountsProvider := accounts.NewProvider(accountsRepository)
	accountsGroup := route.Group("/api/v1/accounts")
	{
		accountsGroup.POST("/register", accountsProvider.Register)
		accountsGroup.POST("/login", accountsProvider.Login)

		// Protected routes - require valid JWT token in Authorization header
		protected := accountsGroup.Group("").Use(middleware.RequireAuth())
		{
			protected.GET("/", accountsProvider.GetAccounts)
			protected.GET("/me", accountsProvider.GetMe)
		}
	}

	invoiceRepository := invoice.NewRepository(db)
	invoiceProvider := invoice.NewProvider(invoiceRepository)
	invoiceGroup := route.Group("/api/v1/invoices")
	{
		invoiceGroup.GET("/:id", invoiceProvider.GetInvoiceByID)
	}

	subscriptionRepository := subscription.NewRepository(db)
	subscriptionProvider := subscription.NewProvider(subscriptionRepository)
	subscriptionGroup := route.Group("/api/v1/subscriptions")
	{
		subscriptionGroup.GET("/:id", subscriptionProvider.GetSubscriptionByID)
	}

	if err := route.Run(); err != nil {
		log.Fatal(err)
	}

}
