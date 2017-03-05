package gotweet
import (
	"encoding/json"
	"github.com/mrjones/oauth"
	"io/ioutil"
)

type App struct {
	consumer    *oauth.Consumer
	accessToken *oauth.AccessToken
}

type APIkeys struct{
	APIKey string
	APISecret string
	AccessToken string
	AccessTokenSecret string
}
type Twitter struct {
	Users TwitterUsers
}
type TwitterUsers struct {
	contributors_enabled bool
	created_at string

}

func TwitterServiceProvider() oauth.ServiceProvider {
	return oauth.ServiceProvider{RequestTokenUrl: "http://api.twitter.com/oauth/request_token",
		AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
		AccessTokenUrl:    "https://api.twitter.com/oauth/access_token"}
}
func Newapp(a APIkeys, sv oauth.ServiceProvider) *App {
	app := new(App)
	app.consumer = oauth.NewConsumer(a.APIKey, a.APISecret, sv)
	app.accessToken = &oauth.AccessToken{Token:a.AccessToken, Secret: a.AccessTokenSecret}
	return app
}
func (t *App) Get(url string, params map[string]string) (interface{}, error) {
	response, err := t.consumer.Get(url, params, t.accessToken)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// decode
	var result interface{}
	err = json.Unmarshal(b, &result)
	return result, err
}

func (t *App) Post(url string, params map[string]string) (interface{}, error) {
	response, err := t.consumer.Post(url, params, t.accessToken)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// decode
	var result interface{}
	err = json.Unmarshal(b, &result)
	return result, err
}