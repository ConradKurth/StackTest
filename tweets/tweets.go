package tweets

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/conradkurth/twitter/cache"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Twitter struct {
	consumerKey    string
	consumerSecret string
	bearer         string
}

const (
	baseUrl = "https://api.twitter.com/1.1/statuses/user_timeline.json"
)

type Mentions struct {
	Name string `json:"screen_name"`
}

type TwitterResp struct {
	Created string `json:"created_at"`
	Text    string `json:"text"`
	//Id       int64  `json:"id"`
	IdStr    string `json:"id_str"`
	Entities struct {
		M []Mentions `json:"user_mentions"`
	} `json:"entities"`
}

type FResp struct {
	Created  string   `json:"created"`
	Text     string   `json:"text"`
	Mentions []string `json:"mentions"`
}

type Pagination struct {
	Resp []FResp `json:"tweets"`
	//Id    int64   `json:"id"`
	IdStr string `json:"id_str"`
}

var t *Twitter

func NewTwitter(key, secret string) *Twitter {
	t := &Twitter{consumerKey: key, consumerSecret: secret}
	err := t.getBearerToken()
	if err != nil {
		log.Fatal("Error getting bearer token", err)
	}
	return t
}

func (t *Twitter) getBearerToken() error {
	e := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s",
		url.QueryEscape(t.consumerKey),
		url.QueryEscape(t.consumerSecret))))

	reqBody := bytes.NewBuffer([]byte(`grant_type=client_credentials`))
	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", reqBody)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", e))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Add("Content-Length", strconv.Itoa(reqBody.Len()))

	//Issue the request and get the bearer token from the JSON you get back
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	type BToken struct {
		AccessToken string `json:"access_token"`
	}
	var b BToken
	err = json.Unmarshal(respBody, &b)
	if err != nil {
		return err
	}
	t.bearer = b.AccessToken
	return nil
}

func (t *Twitter) getRequest(u string, r interface{}) error {

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.bearer))

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotModified {
		return errors.New("api error, response code: " + strconv.Itoa(resp.StatusCode))
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Poorly formed respone from twitter")
	}

	//fmt.Println(string(b))

	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}
	return nil
}

func (t *Twitter) GetTweets(name, maxId string) (Pagination, error) {

	i := cache.GetItem(name + maxId)
	if i != nil {
		return i.(Pagination), nil
	}

	u := baseUrl + "?" + "screen_name=" + name + "&count=" + "25"
	if maxId != "" {
		u += "&max_id=" + maxId
	}

	var r []TwitterResp
	err := t.getRequest(u, &r)
	if err != nil {
		return Pagination{}, err
	}

	resp := Pagination{IdStr: r[len(r)-1].IdStr}

	f := make([]FResp, 0, len(r))
	for _, v := range r {
		t := FResp{Created: v.Created, Text: v.Text}
		for _, m := range v.Entities.M {
			t.Mentions = append(t.Mentions, m.Name)
		}
		f = append(f, t)
	}

	resp.Resp = f
	cache.AddItem(name, resp)
	return resp, nil
}
