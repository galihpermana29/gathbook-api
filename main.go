package main

import (
	"learn-golang/controllers"
	"learn-golang/initializers"
	"learn-golang/middlewares"
	"learn-golang/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func loadDatabase() {
	initializers.Connect()
	initializers.DB.AutoMigrate(
		&models.User{},
		&models.Book{},
		&models.Topic{},
		&models.SubTopic{},
		&models.Content{},
		&models.Buyer{},
	)
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func routes() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", " Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World!",
		})
	})

	router.POST("/auth/signup", controllers.CreateUser)
	router.POST("/auth/login", controllers.Login)

	router.GET("/user/:id", middlewares.CheckAuth, controllers.GetUserById)
	router.POST("/book", middlewares.CheckAuth, controllers.CreateBook)
	router.GET("/book", middlewares.CheckAuth, controllers.GetBooks)
	router.GET("/book/:id", middlewares.CheckAuth, controllers.GetBookByID)
	router.PUT("/book/:id", middlewares.CheckAuth, controllers.UpdateBook)

	router.POST("/images", middlewares.CheckAuth, controllers.UploadImages)
	router.GET("/image/:filename", controllers.ServeImage)

	router.POST("/buy", middlewares.CheckAuth, controllers.BuyBook)
	router.Run()
}

func main() {
	loadEnv()
	loadDatabase()
	routes()
}
