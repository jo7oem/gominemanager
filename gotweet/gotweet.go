package gotweet

import (
	"encoding/json"
	"github.com/mrjones/oauth"
	"io/ioutil"
	"reflect"
	"unicode/utf8"
)

type App struct {
	consumer    *oauth.Consumer
	accessToken *oauth.AccessToken
}
type APIkeys struct {
	APIKey            string
	APISecret         string
	AccessToken       string
	AccessTokenSecret string
}
type objects map[string]string

type Users struct {
	Contributors_enabled               bool
	Created_at                         string
	Default_profile                    bool
	Default_profile_image              bool
	Description                        string
	Entities                           Entities
	Favourites_count                   int
	Follow_request_sent                interface{}
	Following                          interface{}
	Followers_count                    int
	Friends_count                      int
	Geo_enabled                        bool
	Id                                 int64
	Id_str                             string
	Is_translator                      bool
	Lang                               string
	Listed_count                       int
	Location                           string
	Name                               string
	Notifications                      bool
	Profile_background_color           string
	Profile_background_image_url       string
	Profile_background_image_url_https string
	Profile_background_tile            bool
	Profile_banner_url                 string
	Profile_image_url                  string
	Profile_image_url_https            string
	Profile_link_color                 string
	Profile_sidebar_border_color       string
	Profile_sidebar_fill_color         string
	Profile_text_color                 string
	Profile_use_background_image       bool
	Protected                          bool
	Screen_name                        string
	Show_all_inline_media              bool
	Status                             objects //Tweets
	Statuses_count                     int
	Time_zone                          string
	Url                                string
	Utc_offset                         int
	Verified                           bool
	Withheld_in_countries              string
	Withheld_scope                     string
}
type Entities struct {
	Hashtags []struct {
		Indices []int
		Text    string
	}
	Media []struct {
		Display_url     string
		Expanded_url    string
		Id              int64
		Id_str          string
		Indices         []int
		Media_url       string
		Media_url_https string
		Sizes           struct {
			Thumb  Size
			Large  Size
			Medium Size
			Small  Size
		}
		Source_status_id     int64
		Source_status_id_str string
		Type                 string
		Url                  string
	}
	Urls []struct {
		Display_url  string
		Expanded_url string
		Idices       []int
		Url          string
	}
	User_mentions []struct {
		Id          int64
		Id_str      string
		Indices     []int
		Name        string
		Screen_name string
	}
}
type Size struct {
	H      int
	Resize string
	W      int
}
type Tweets struct {
	Contributors              string `json:"contributors"` //非推奨
	Coordinates               interface{}
	Created_at                string
	Current_user_retweet      objects
	Entities                  Entities
	Favorite_count            int
	Favorited                 bool
	Filter_level              string
	Geo                       objects
	Id                        int64
	Id_str                    string `json:"id_str"`
	In_reply_to_screen_name   string
	In_reply_to_status_id     int64
	In_reply_to_status_id_str string
	In_reply_to_user_id       int64
	In_reply_to_user_id_str   string
	Lang                      string
	Place                     struct {
		Attributes   objects
		Bounding_box objects
		Country      string
		Country_code string
		Full_name    string
		Id           string
		Name         string
		Place_type   string
		Url          string
	}
	Possibly_sensitive    bool
	Quoted_status_id      int64
	Quoted_status_id_str  string
	Quoted_status         objects
	Scopes                objects
	Retweet_count         int
	Retweeted             bool
	Retweeted_status      objects
	Source                string
	Text                  string
	Truncated             bool
	User                  Users
	Withheld_copyright    bool
	Withheld_in_countries []string
	Withheld_scope        string
}
type Coordinates map[string]string

func twitterServiceProvider() oauth.ServiceProvider {
	return oauth.ServiceProvider{RequestTokenUrl: "http://api.twitter.com/oauth/request_token",
		AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
		AccessTokenUrl:    "https://api.twitter.com/oauth/access_token"}
}
func Newapp(a APIkeys) *App {
	app := new(App)
	app.consumer = oauth.NewConsumer(a.APIKey, a.APISecret, twitterServiceProvider())
	app.accessToken = &oauth.AccessToken{Token: a.AccessToken, Secret: a.AccessTokenSecret}
	return app
}
func (t *App) Get(url string, params map[string]string, result interface{}) (interface{}, error) {
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
	var result interface{}
	// decoded
	err = json.Unmarshal(b, &result)
	return result, err
}
func (t *App) Tweet(params map[string]string) (Tweets, error) {
	result := Tweets{}
	url := "https://api.twitter.com/1.1/statuses/update.json"
	response, err := t.consumer.Post(url, params, t.accessToken)
	if err != nil {
		return result, err
	}
	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(b, &result)
	return result, err

}

func SliceStrlen(slice []string) []int {
	result := make([]int, len(slice))
	for n, s := range slice {
		result[n] = utf8.RuneCountInString(s)
	}
	return result
}
func SliceFindfunc(t interface{}, f func(interface{}) bool) []int {
	result := []int{}
	buf := SliceInterface(t)
	for n, r := range buf {
		if f(r) {
			result = append(result, n)
		} else {
			continue
		}
	}
	return result
}
func SliceInterface(t interface{}) []interface{} {
	buf := []interface{}{}
	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(t)
		for i := 0; i < s.Len(); i++ {
			b := s.Index(i)
			buf = append(buf, b.Interface())
		}
	case reflect.String:
		s := reflect.ValueOf(t)
		for i := 0; i < s.Len(); i++ {
			b := s.Index(i)
			buf = append(buf, b.Interface())
		}
	}
	return buf
}
func ForceSplitStringN(limit int, s string) []string {
	sr := []rune(s)
	wc := len(sr)
	result := []string{}
	dosp := wc / limit
	for i := 0; i <= dosp; i++ {
		if i == dosp {
			result = append(result, string(sr[limit*i:]))
			break
		}
		result = append(result, string(sr[limit*i:limit*(i+1)]))
	}
	return result
}
func SplitRunelimit(limit int, s string, sep rune) []string {
	rs := []rune(s)
	result := []string{}
	f := func(i interface{}) bool {
		irs := i.(rune)
		if irs == sep {
			return true
		}
		return false
	}
	splitbuf := isuseindexl(SliceFindfunc(rs, f), limit)
	starti := 0
	for _, idx := range splitbuf {
		result = append(result, string(rs[starti:idx]))
		starti = idx
	}
	result = append(result, ForceSplitStringN(limit, string(rs[starti:]))...)
	return result
}
func isuseindexl(il []int, limit int) []int {
	bfri := 0
	starti := 0
	result := []int{}
	for _, i := range il {
		nwc := i - starti
		if nwc < limit {
			bfri = i
			continue
		} else {
			if (starti != bfri) && (starti-bfri <= limit) {
				result = append(result, bfri)
				starti = bfri
			}
			if i-bfri <= limit {
				starti = i
				bfri = i
				continue
			}
			for j := 1; j <= (i-bfri)/limit; j++ {
				result = append(result, starti+limit*j)
			}
			result = append(result, i)
			starti = i
			bfri = i
			continue
		}
	}
	return result
}
func SplitRuneslimit(limit int, s string, seps []rune) []string {
	rs := []rune(s)
	result := []string{}
	splitbuf := [][]int{}
	for _, sep := range seps {
		f := func(i interface{}) bool {
			irs := i.(rune)
			if irs == sep {
				return true
			}
			return false
		}
		splitbuf = append(splitbuf, isuseindexl(SliceFindfunc(rs, f), limit))
	}
	/*starti := 0
	for _, idx := range splitbuf {
		result = append(result, string(rs[starti:idx]))
		starti = idx
	}
	result = append(result, ForceSplitStringN(limit, string(rs[starti:]))...)
	*/return result
}
