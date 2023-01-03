package routes

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddPublicMapRoutes(rg *gin.RouterGroup, mapColl *mongo.Collection) {
	rg.GET("", func(c *gin.Context) {
		cursor, err := mapColl.Find(
			context.TODO(),
			bson.D{{"public", true}},
		)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		var results []bson.M
		if err = cursor.All(context.TODO(), &results); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		c.JSON(200, gin.H{"data": results})
	})

	rg.GET("/:id", func(c *gin.Context) {
		id, _ := primitive.ObjectIDFromHex(c.Param("id"))
		var result bson.M
		err := mapColl.FindOne(
			context.TODO(),
			bson.D{{"_id", id}, {"public", true}},
		).Decode(&result)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": result})
	})

}
