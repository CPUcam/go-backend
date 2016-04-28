package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/revel/revel"
	"github.com/revel/samples/facebook-oauth2/app/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

type App struct {
	*revel.Controller
}

var FACEBOOK = &oauth2.Config{
	ClientID:     "1558186451167627",
	ClientSecret: "5fa1d0df66447e984e09ce64ff8edffe",
	Scopes:       []string{},
	Endpoint:     facebook.Endpoint,
	RedirectURL:  "http://loisant.org:9000/Application/Auth",
}

func (c App) Index() revel.Result {
	return c.Render()
}

type FileInfo struct {
	ContentType string
	Filename    string
	RealFormat  string `json:",omitempty"`
	Resolution  string `json:",omitempty"`
	Size        int
	Status      string `json:",omitempty"`
}

func (c App) Login() revel.Result {
	u := c.connected()
	me := map[string]interface{}{}
	if u != nil && u.AccessToken != "" {
		resp, _ := http.Get("https://graph.facebook.com/me?access_token=" +
			url.QueryEscape(u.AccessToken))
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&me); err != nil {
			revel.ERROR.Println(err)
		}
		revel.INFO.Println(me)
	}

	authUrl := FACEBOOK.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return c.Render(me, authUrl)
}

func (c App) Auth(code string) revel.Result {

	tok, err := FACEBOOK.Exchange(oauth2.NoContext, code)
	if err != nil {
		revel.ERROR.Println(err)
		return c.Redirect(App.Login)
	}

	user := c.connected()
	user.AccessToken = tok.AccessToken
	return c.Redirect(App.Login)
}

func setuser(c *revel.Controller) revel.Result {
	var user *models.User
	if _, ok := c.Session["uid"]; ok {
		uid, _ := strconv.ParseInt(c.Session["uid"], 10, 0)
		user = models.GetUser(int(uid))
	}
	if user == nil {
		user = models.NewUser()
		c.Session["uid"] = fmt.Sprintf("%d", user.Uid)
	}
	c.RenderArgs["user"] = user
	return nil
}

func init() {
	revel.InterceptFunc(setuser, revel.BEFORE, &App{})
}

func (c App) connected() *models.User {
	return c.RenderArgs["user"].(*models.User)
}

func (c App) SignUp() revel.Result {
	return c.Render()
}

func (c App) Upload() revel.Result {
	return c.Render()
}

func (c App) Student() revel.Result {
	return c.Render()
}

func (c App) Viewer() revel.Result {
	return c.Render()
}
