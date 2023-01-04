package routes

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary getUserData
// @Schemes
// @Description getUserData
// @Tags user
// @Router /api/user [get]
func getUserData(userColl *mongo.Collection) func(ctx *gin.Context) {
	return func(c *gin.Context) {
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
	}
}

func AddUserRoutes(rg *gin.RouterGroup, userColl *mongo.Collection) {
	rg.GET("", getUserData(userColl))
}
