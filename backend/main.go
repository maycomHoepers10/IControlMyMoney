package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"

	database "home_money/database"
	handlers "home_money/handlers"
	"home_money/middlewares"
)

func main() {
	viper.SetConfigFile("config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config file: %s", err)
	}

	// Obter a chave secreta do JWT do arquivo de configuração
	jwtSecret := viper.GetString("jwt.secret")

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}

	userHandler := handlers.NewUserHandler(db, []byte(jwtSecret))
	categoryHandler := handlers.NewCategoryHandler(db)
	importCSV := handlers.NewImportCSVHandler(db, categoryHandler)
	transactionHandler := handlers.NewTransactionHandler(db)
	accountHandler := handlers.NewAccountHandler(db)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Carregar a configuração do JWT
	jwtConfig, err := middlewares.LoadJWTConfig()
	if err != nil {
		log.Fatalf("failed to load JWT config: %s", err)
	}

	// Rota para login (sem o middleware JWT)
	e.POST("/users/login", userHandler.Login)
	e.POST("/users/register", userHandler.RegisterUser)
	e.POST("/validate-token", userHandler.ValidateToken)

	// Grupo de rotas protegidas pelo JWT
	authGroup := e.Group("")
	authGroup.Use(middlewares.JWTMiddleware(jwtConfig))
	authGroup.Use(middlewares.ExtractUserIDMiddleware(db))

	//Para fazer as operações internas, deve usar o jwt

	//Categorias
	authGroup.POST("/category", categoryHandler.CreateCategory)
	authGroup.PUT("/category/:id", categoryHandler.UpdateCategory)
	authGroup.DELETE("/category/:id", categoryHandler.DeleteCategory)
	authGroup.GET("/categories", categoryHandler.ListCategories)
	authGroup.GET("/category/:id", categoryHandler.GetCategory)

	//Transações
	authGroup.POST("/transaction", transactionHandler.CreateTransaction)
	authGroup.PUT("/transaction/:id", transactionHandler.UpdateTransaction)
	authGroup.DELETE("/transaction/:id", transactionHandler.DeleteTransaction)
	authGroup.GET("/transactions", transactionHandler.ListTransactions)
	authGroup.GET("/transaction/:id", transactionHandler.GetTransaction)
	authGroup.POST("/transactions/approve", transactionHandler.ApproveTransaction)

	//Contas
	authGroup.POST("/financialAccount", accountHandler.CreateAccount)
	authGroup.PUT("/financialAccount/:id", accountHandler.UpdateAccount)
	authGroup.DELETE("/financialAccount/:id", accountHandler.DeleteAccount)
	authGroup.GET("/financialAccounts", accountHandler.ListAccounts)
	authGroup.GET("/financialAccount/:id", accountHandler.GetAccount)

	//Import CSVfinancial accounts
	authGroup.POST("/upload/csv", importCSV.UploadCSV)

	address := fmt.Sprintf(":%s", "54165")
	e.Start("0.0.0.0" + address)
}
