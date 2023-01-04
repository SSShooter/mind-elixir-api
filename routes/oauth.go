package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GithubResp struct {
	Error             string
	Error_description string
	Error_uri         string

	Access_token string
	Token_type   string
	Scope        string
}

type UserData struct {
	Id        int    `json:"id"`
	NodeId    string `json:"node_id"`
	AvatarUrl string `json:"avatar_url"`
	Login     string `json:"login"`
	Name      string `json:"Name"`
	Company   string `json:"company"`
	Location  string `json:"location"`
	Email     string `json:"Email"`
	Bio       string `json:"bio"`
}

func getToken(url string) (GithubResp, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Add("accept", `application/json`)
	resp, err := client.Do(req)
	if err != nil {
		return GithubResp{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GithubResp{}, err
	}
	var data GithubResp
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("error:", err)
		return GithubResp{}, err
	}
	fmt.Printf("%s", data.Access_token)
	return data, nil
}

func fetchUserData(token string) (UserData, error) {
	client := &http.Client{}
	dataReq, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	dataReq.Header.Add("Authorization", "token "+token)
	dataResp, err := client.Do(dataReq)
	if err != nil {
		return UserData{}, err
	}

	var data UserData

	json.NewDecoder(dataResp.Body).Decode(&data)
	fmt.Printf("\ndata:%+v\n", data)
	return data, nil
}

func updateUserData(coll *mongo.Collection, data UserData) error {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"id", data.Id}}
	update := bson.D{{"$set", data}}
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

// @Summary githubAuth
// @Schemes
// @Description Github Oauth Callback
// @Param code query string true "code"
// @Tags auth
// @Success 200
// @Router /oauth/redirect [get]
func githubAuth(userColl *mongo.Collection) func(c *gin.Context) {
	return func(c *gin.Context) {
		code, _ := c.GetQuery("code")
		clientId := os.Getenv("CLIENT_ID")
		clientSecret := os.Getenv("CLIENT_SECRET")
		redirectDomain := os.Getenv("REDIRECT_DOMAIN")
		url := "https://github.com/login/oauth/access_token?client_id=" + clientId + "&client_secret=" + clientSecret + "&code=" + code
		data, err := getToken(url)
		if err != nil {
			c.JSON(200, gin.H{"error": "Can not connect to GitHub"})
			return
		}
		if data.Access_token == "" {
			c.JSON(200, gin.H{"error": "Can not get token"})
			return
		}
		userData, err := fetchUserData(data.Access_token)
		session := sessions.Default(c)
		session.Set("loginId", userData.Id)
		session.Save()
		err = updateUserData(userColl, userData)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.Redirect(301, redirectDomain)
	}
}

func AddOauthRoutes(rg *gin.RouterGroup, userColl *mongo.Collection) {
	rg.GET("/redirect", githubAuth(userColl))
}
