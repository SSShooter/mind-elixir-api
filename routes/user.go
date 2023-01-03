package routes

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserData(rg *gin.RouterGroup, userColl *mongo.Collection) {
	rg.GET("", func(c *gin.Context) {
		loginId := c.MustGet("loginId").(int)
		var result bson.M
		err := userColl.FindOne(
			context.TODO(),
			bson.D{{"id", loginId}},
		).Decode(&result)
		if err != nil {
			c.JSON(400, gin.H{"error": "no data"})
			return
		}
		c.JSON(200, gin.H{"data": result})
	})
}
