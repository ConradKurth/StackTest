package main

import (
	"github.com/conradkurth/twitter/tokens"
	"github.com/conradkurth/twitter/tweets"
	"github.com/plimble/ace"
	"log"
	"os"
)

var t *tweets.Twitter

func Auth(c *ace.C) {
	// t := c.Request.Header.Get("Authorization")
	// if err := tokens.VerifyJWT(t); err != nil {
	// 	c.JSON(400, map[string]string{"error": "Invalid token"})
	// 	c.Abort()
	// 	return
	// }
	c.Next()
}

func RegisterUser(c *ace.C) {
	t := tokens.GenerateToken()
	c.JSON(200, map[string]string{"token": t})
}

func GetTweets(c *ace.C) {
	name := c.MustQueryString("name", "")
	id := c.MustQueryString("id", "")

	if name == "" {
		c.JSON(400, map[string]string{"error": "No username provided"})
		return
	}
	t, err := t.GetTweets(name, id)
	if err != nil {
		c.JSON(500, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(200, t)
}

func main() {

	t = tweets.NewTwitter(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))

	a := ace.New()
	a.GET("/register", RegisterUser)
	a.GET("/tweets", Auth, GetTweets)

	a.Run(GetPort())
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("[-] No PORT environment variable detected. Setting to ", port)
	}
	return ":" + port
}
