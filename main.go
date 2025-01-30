package main

import (
	"context"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "GoRecipe/docs" // Required for Swagger docs

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"GoRecipe/controllers"
	"GoRecipe/repository"
)

// @title Go Recipe API
// @version 1.0
// @description This is a sample API for managing recipes.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host http://go-kitchen-lb-1913677782.ap-northeast-1.elb.amazonaws.com
// @BasePath /

func main() {
	// Optionally configure logrus (set format, level, etc.).
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)

	// Get MongoDB URI from environment variable
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		logrus.Fatal("MONGODB_URI environment variable is not set")
	}

	// Initialize MongoDB client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to Mongo")
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		logrus.WithError(err).Fatal("Failed to ping Mongo")
	}

	logrus.Info("Connected to MongoDB")

	// Choose your database
	db := client.Database("recipe_db")

	// Initialize repository
	recipeRepo := repository.NewRecipeRepository(db, "recipes")

	// Initialize controllers
	recipeCtrl := controllers.NewRecipeController(recipeRepo)

	// Setup Gin router
	r := gin.Default()

	// 🔥 Enable CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins (change this for security)
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Version v1 routes
	v1 := r.Group("/v1")
	{
		v1.POST("/recipes", recipeCtrl.CreateRecipe)
		v1.GET("/recipes", recipeCtrl.GetAllRecipes)
		v1.GET("/recipes/:id", recipeCtrl.GetRecipeByID)
		v1.PUT("/recipes/:id", recipeCtrl.UpdateRecipe)
		v1.DELETE("/recipes/:id", recipeCtrl.DeleteRecipe)
	}

	// Run the server
	logrus.Info("Starting server on port 8080...")
	if err := r.Run(":8080"); err != nil {
		logrus.WithError(err).Fatal("Failed to run server")
	}
}
