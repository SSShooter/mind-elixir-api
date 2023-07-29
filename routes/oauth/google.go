package oauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUserData struct {
	Id            string `json:"id"`
	FamilyName    string `json:"family_name"`
	GivenName     string `json:"given_name"`
	Name          string `json:"Name"`
	Locale        string `json:"locale"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Pictrue       string `json:"picture"`
}

func getConf() *oauth2.Config {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	backendUrl := os.Getenv("BACKEND_URL")

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  backendUrl + "/oauth/google/redirect",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
			"openid",
		},
		Endpoint: google.Endpoint,
	}
	return conf
}

func googleLogin(c *gin.Context) {
	url := getConf().AuthCodeURL("state")
	c.Redirect(301, url)
}

// @Summary googleAuth
// @Schemes
// @Description Github Oauth Callback
// @Param code query string true "code"
// @Tags auth
// @Success 200
// @Router /oauth/google/redirect [get]
func googleAuth(userColl *mongo.Collection) func(c *gin.Context) {
	return func(c *gin.Context) {
		var code string
		conf := getConf()
		code, _ = c.GetQuery("code")
		// Handle the exchange code to initiate a transport.
		tok, err := conf.Exchange(oauth2.NoContext, code)
		if err != nil {
			c.JSON(400, err.Error())
			return
		}
		fmt.Println("tok", tok)
		client := conf.Client(oauth2.NoContext, tok)
		fmt.Println("client!")
		resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + tok.AccessToken)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}

		var data GoogleUserData
		// print data
		err = json.Unmarshal(body, &data)
		fmt.Println("data", data)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		session := sessions.Default(c)
		session.Set("loginId", data.Id)
		session.Save()
		err = updateUserData(userColl, data)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		frontendUrl := os.Getenv("FRONTEND_URL")
		c.Redirect(301, frontendUrl)
	}
}
