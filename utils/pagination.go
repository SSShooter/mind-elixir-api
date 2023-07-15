package utils

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 定义一个结构体来封装分页参数
type Pagination struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"pageSize,default=10"`
}

// 定义一个函数来处理分页请求
func GetPaginatedResults(c *gin.Context, collection *mongo.Collection, query bson.M) (gin.H, error) {
	var pagination Pagination
	if err := c.ShouldBindQuery(&pagination); err != nil {
		return nil, err
	}
	println(pagination.Page, pagination.PageSize)
	// 设置查询选项，包括 skip 和 limit
	options := options.Find()
	options.SetSkip(int64((pagination.Page - 1) * pagination.PageSize))
	options.SetLimit(int64(pagination.PageSize))

	cursor, err := collection.Find(context.Background(), query, options)
	if err != nil {
		return nil, err
	}
	total, err := collection.CountDocuments(context.Background(), query)
	// 遍历结果集并将数据存入一个 slice 中
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
