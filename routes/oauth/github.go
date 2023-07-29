package oauth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type GithubResp struct {
	Error             string
	Error_description string
	Error_uri         string

	Access_token string
	Token_type   string
	Scope        string
}

type GithubUserdata struct {
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
		return GithubResp{}, err
	}
	return data, nil
}

func fetchUserData(token string) (GithubUserdata, error) {
	client := &http.Client{}
	dataReq, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	dataReq.Header.Add("Authorization", "token "+token)
	dataResp, err := client.Do(dataReq)
	if err != nil {
		return GithubUserdata{}, err
	}

	var data GithubUserdata

	json.NewDecoder(dataResp.Body).Decode(&data)
	return data, nil
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
		session.Set("loginId", strconv.Itoa(userData.Id))
		session.Save()
		err = updateUserData(userColl, userData)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		frontendUrl := os.Getenv("FRONTEND_URL") // 写在外面会读不到，可能是因为 godotenv 还没把他读进来
		c.Redirect(301, frontendUrl)
	}
}

func githubLogin(c *gin.Context) {
	clientId := os.Getenv("CLIENT_ID")
	backendUrl := os.Getenv("BACKEND_URL")
	loginUrl := "https://github.com/login/oauth/authorize?client_id=" + clientId + "&redirect_uri=" + backendUrl + "/oauth/github/redirect"
	c.Redirect(301, loginUrl)
}
