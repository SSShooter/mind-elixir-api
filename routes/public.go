package routes

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/SSShooter/mind-elixir-backend-go/utils"
)

// @Summary getAllPublicMaps
// @Schemes
// @Description getAllPublicMaps
// @Tags public
// @Router /api/public [get]
// @Param name query string false "Map Name"
func getAllPublicMaps(mapColl *mongo.Collection) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		name := c.Query("name")
		query := bson.M{"public": true, "name": bson.M{"$regex": name, "$options": "i"}}
		results, err := utils.GetPaginatedResults(c, mapColl, query)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, results)
	}
}

// @Summary getPublicMap
// @Schemes
// @Description getPublicMap
// @Tags public
// @Param id path string true "Map ID"
// @Router /api/public/{id} [get]
func getPublicMap(mapColl *mongo.Collection) func(ctx *gin.Context) {
	return func(c *gin.Context) {
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
	}
}

// TODO add pw protection

func AddPublicMapRoutes(rg *gin.RouterGroup, mapColl *mongo.Collection) {
	rg.GET("", getAllPublicMaps(mapColl))
	rg.GET("/:id", getPublicMap(mapColl))
}
