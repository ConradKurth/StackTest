# Twitter

The design of the program is fairly simple. I have a package that will manage interaction with the Twitter API. 

The user will pass in the consumer key and secret key to the object to be created. It will then fetch a token from Twitter
to be used in each subsiquent request. 

With the returned Twitter struct, you can request various users tweets which will be formatted a little upon return. If any errors are encountered they will be handled and returned. 

As for the security the only endpoint that need authorization is the /tweets endpoint which will need a JWT that can be returned from the /register endpoint

With each request to the /tweets endpoint, set the header Authorization with the returned token. The token will expire after an hour

The 25 latest tweets will be returned containing the time of post, text and an array of users mentioned in the tweet. 


To set up this enviornment locally, you will need to set up 3 environment variables. The 'SECRET' which is some random password. Then you will need to get your consumer secret/key from Twitter when you create an app. The respective variables are 'CONSUMER_KEY' and 'CONSUMER_SECRET'


To get the project just do: go get github.com/conradkurth/twitter . The run main.go to start a sever locally. 

To get a project up and running on heroku, follow https://devcenter.heroku.com/articles/getting-started-with-go#set-up and 
instead of pulling down their project, pull down this one.
