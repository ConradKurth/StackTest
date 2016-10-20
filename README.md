# StackTest

This is a simple program that demonstrates some basic security and fetching data from a third party API.

To get a new token, do a GET request at /register. 

With each request to the /tweets endpoint, set the header Authorization with the returned token. The token will expire after an hour


The 25 latest tweets will be returned containing the time of post, text and an array of users mentioned in the tweet. 

