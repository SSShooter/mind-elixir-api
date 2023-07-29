package oauth

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/SSShooter/mind-elixir-backend-go/models"
)

func ConvertToUser(u interface{}) (models.User, error) {
	switch v := u.(type) {
	case GoogleUserData:
		return models.User{
			Id:     v.Id,
			Name:   v.Name,
			Email:  v.Email,
			Avatar: v.Pictrue,
			From:   "Google",
		}, nil
	case GithubUserdata:
		return models.User{
			Id:     strconv.Itoa(v.Id),
			Name:   v.Name,
			Email:  v.Email,
			Avatar: v.AvatarUrl,
			From:   "Github",
		}, nil
	default:
		return models.User{}, fmt.Errorf("unsupported user type")
	}
}

func updateUserData(coll *mongo.Collection, origin interface{}) error {
	data, _ := ConvertToUser(origin)
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.M{"id": data.Id}
	update := bson.M{"$set": data}
	var updatedDocument bson.M
	err := coll.FindOneAndUpdate(
		context.TODO(),
		filter,
		update,
		opts,
	).Decode(&updatedDocument)
	if err != nil {
		return err
	}
	return nil
}

func AddOauthRoutes(rg *gin.RouterGroup, userColl *mongo.Collection) {
	rg.GET("/github/login", githubLogin)
	rg.GET("/github/redirect", githubAuth(userColl))
	rg.GET("/google/login", googleLogin)
	rg.GET("/google/redirect", googleAuth(userColl))
}
