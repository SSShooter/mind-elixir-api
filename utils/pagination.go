package utils

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Pagination struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"pageSize,default=10"`
}

func GetPaginatedResults(c *gin.Context, collection *mongo.Collection, query bson.M) (gin.H, error) {
	var pagination Pagination
	if err := c.ShouldBindQuery(&pagination); err != nil {
		return nil, err
	}
	options := options.Find()
	options.SetSkip(int64((pagination.Page - 1) * pagination.PageSize))
	options.SetLimit(int64(pagination.PageSize))
	options.SetSort(bson.D{{Key: "_id", Value: -1}})

	cursor, err := collection.Find(context.Background(), query, options)
	if err != nil {
		return nil, err
	}
	total, err := collection.CountDocuments(context.Background(), query)

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	// var results []bson.M
	// for cursor.Next(context.Background()) {
	// 	var result bson.M
	// 	err := cursor.Decode(&result)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	results = append(results, result)
	// }

	data := gin.H{
		"data":     results,
		"page":     pagination.Page,
		"pageSize": pagination.PageSize,
		"total":    total,
	}
	return data, nil
}
