// @title           Billing Service API
// @version         1.0
// @description     Subscription billing service for managing plans, subscriptions, and payment methods.
// @host            localhost:8080
// @BasePath        /

package main

import (
	"billingService/backend/internal/accounts"
	"billingService/backend/internal/middleware"
	"billingService/backend/internal/plans"
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
		log.Fatal("Error loading environment variables.")
	}

	// Translate DB-level errors (like unique constraint violations) into typed errors like gorm.ErrDuplicatedKey
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{TranslateError: true})
	if err != nil {
		log.Fatal(err)
	}

	route := gin.Default()

	// Swagger UI endpoint
	// https://github.com/swaggo/swag/blob/master/README.md#declarative-comments-format
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	plansRepository := plans.NewPlansRepository(db)
	plansProvider := plans.NewProvider(plansRepository)
	plansGroup := route.Group("/api/v1/plans")
	{
		plansGroup.GET("/", plansProvider.GetPlans)
		plansGroup.GET("/:id", plansProvider.GetPlanById)
		plansGroup.POST("/create/", plansProvider.CreatePlan)
		plansGroup.PATCH("/update/:id", plansProvider.UpdatePlanById)
		plansGroup.PATCH("/update/status/:id", plansProvider.UpdatePlanStatusById)
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

	route.Run()
}
