package main

import (
	"fmt"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func main() {
	config := oauth1.NewConfig(os.Getenv("DELPHIS_TWITTER_CONSUMER_KEY"), os.Getenv("DELPHIS_TWITTER_CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"))

	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	ids, resp, err := client.Followers.IDs(&twitter.FollowerIDParams{ScreenName: "nedrockson"})

	if err != nil {
		fmt.Printf("err: %+v\n", err)
	}
	fmt.Printf("ids: %+v\n", ids)
	fmt.Printf("headers: %+v\n", resp.Header)
}
