package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/SSShooter/mind-elixir-backend-go/docs"
	"github.com/SSShooter/mind-elixir-backend-go/middlewares"
	"github.com/SSShooter/mind-elixir-backend-go/routes"
	"github.com/SSShooter/mind-elixir-backend-go/routes/oauth"
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

// @Summary login
// @Schemes
// @Description It will redirect to GitHub login page
// @Tags auth
// @Success 200
// @Router /login [get]

// @Summary logout
// @Schemes
// @Description Clear session
// @Tags auth
// @Success 200
// @Router /logout [get]
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

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	sessionSecret := os.Getenv("SESSION_SECRET")
	store := cookie.NewStore([]byte(sessionSecret))
	store.Options(sessions.Options{
		MaxAge:   60 * 60 * 24 * 7,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
	r.Use(sessions.Sessions("mindelixir", store))

	AllowOrigin := os.Getenv("FRONTEND_URL")
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{AllowOrigin},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin", "content-type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	oauthGroup := r.Group("/oauth")
	oauth.AddOauthRoutes(oauthGroup, db.Collection("users"))

	api := r.Group("/api")
	api.Use(middlewares.Auth())
	user := api.Group("/user")
	routes.AddUserRoutes(user, db.Collection("users"))
	mapr := api.Group("/map")
	routes.AddMapRoutes(mapr, db.Collection("maps"))

	free := r.Group("/api")
	public := free.Group("/public")

	routes.AddPublicMapRoutes(public, db.Collection("maps"))

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
