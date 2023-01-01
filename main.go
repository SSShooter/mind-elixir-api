package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/SSShooter/mind-elixir-backend-go/middlewares"
	"github.com/SSShooter/mind-elixir-backend-go/routes"
)

func connectDatabase() (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	url := os.Getenv("MONGODB_URL")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}
	db := client.Database("test")
	return db, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	db, err := connectDatabase()
	if err != nil {
		log.Fatal("Db connect fail")
	}

	r := gin.Default()

	sessionSecret := os.Getenv("SESSION_SECRET")
	store := cookie.NewStore([]byte(sessionSecret))
	r.Use(sessions.Sessions("mindelixir", store))

	AllowOrigin := os.Getenv("ALLOW_ORIGIN")
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{AllowOrigin},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin", "content-type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	oauth := r.Group("/oauth")
	routes.AddOauthRoutes(oauth, db.Collection("users"))

	api := r.Group("/api")
	api.Use(middlewares.Auth())
	user := api.Group("/user")
	routes.GetUserData(user, db.Collection("users"))
	mapr := api.Group("/map")
	routes.AddMapRoutes(mapr, db.Collection("maps"))

	free := r.Group("/api")
	public := free.Group("/public")

	routes.AddPublicMapRoutes(public, db.Collection("maps"))

	r.GET("/login", func(c *gin.Context) {
		loginUrl := os.Getenv("LOGIN_URL")
		c.Redirect(301, loginUrl)
	})
	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("loginId", nil)
		session.Save()
		c.JSON(200, gin.H{"msg": "logout"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "7001"
	}
	r.Run(":" + port)
}
