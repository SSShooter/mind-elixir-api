package routes

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/SSShooter/mind-elixir-backend-go/models"
)

func AddMapRoutes(rg *gin.RouterGroup, mapColl *mongo.Collection) {
	rg.GET("/:id", func(c *gin.Context) {
		id, _ := primitive.ObjectIDFromHex(c.Param("id"))
		loginId := c.MustGet("loginId").(int)
		var result bson.M
		err := mapColl.FindOne(
			context.TODO(),
			bson.D{{"_id", id}, {"author", loginId}},
		).Decode(&result)
		if err != nil {
			c.JSON(200, gin.H{"error": err})
			return
		}
		c.JSON(200, gin.H{"data": result})
	})

	rg.PATCH("/:id", func(c *gin.Context) {
		id, _ := primitive.ObjectIDFromHex(c.Param("id"))
		var data map[string]interface{}
		c.ShouldBind(&data)
		loginId := c.MustGet("loginId").(int)
		var result bson.M
		update := bson.D{{"$set", data}}
		err := mapColl.FindOneAndUpdate(
			context.TODO(),
			bson.D{{"_id", id}, {"author", loginId}},
			update,
		).Decode(&result)
		if err != nil {
			c.JSON(200, gin.H{"error": err, "result": result})
			return
		}
		c.JSON(200, gin.H{"data": result})
	})

	rg.DELETE("/:id", func(c *gin.Context) {
		id, _ := primitive.ObjectIDFromHex(c.Param("id"))
		loginId := c.MustGet("loginId").(int)
		var result bson.M
		err := mapColl.FindOneAndDelete(
			context.TODO(),
			bson.D{{"_id", id}, {"author", loginId}},
		).Decode(&result)
		if err != nil {
			c.JSON(200, gin.H{"error": err})
			return
		}
		c.JSON(200, gin.H{"data": result})
	})

	rg.GET("", func(c *gin.Context) {
		loginId := c.MustGet("loginId").(int)
		cursor, err := mapColl.Find(
			context.TODO(),
			bson.D{{"author", loginId}},
		)
		if err != nil {
			c.JSON(200, gin.H{"error": err})
			return
		}
		var results []bson.M
		if err = cursor.All(context.TODO(), &results); err != nil {
			log.Fatal(err)
		}
		c.JSON(200, gin.H{"data": results})
	})

	rg.POST("", func(c *gin.Context) {
		var mapData *models.Map
		c.ShouldBind(&mapData)
		loginId := c.MustGet("loginId").(int)
		mapData.Author = loginId
		fmt.Printf("id:%s", loginId)

		res, err := mapColl.InsertOne(context.TODO(), mapData)
		if err != nil {
			c.JSON(200, gin.H{"error": err})
			return
		}
		c.JSON(200, gin.H{"data": gin.H{"_id": res.InsertedID}})
	})
}
