package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/coreos/pkg/flagutil"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/sirupsen/logrus"
)

func main() {
	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	consumerKey := flags.String("consumer-key", "", "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", "", "Twitter Consumer Secret")
	accessToken := flags.String("access-token", "", "Twitter Access Token")
	accessSecret := flags.String("access-secret", "", "Twitter Access Secret")
	flags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(flags, "TWITTER")

	fmt.Printf("consumerKey: %s, consumerSecret: %s, accessToken: %s, accessSecret: %s", *consumerKey, *consumerSecret, *accessToken, *accessSecret)
	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}

	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// user show
	userShowParams := &twitter.UserShowParams{ScreenName: "nedrockson"}
	user, _, _ := client.Users.Show(userShowParams)
	//fmt.Printf("USERS SHOW:\n%+v\n", user)

	dm, resp, err := client.DirectMessages.EventsNew(&twitter.DirectMessageEventsNewParams{
		Event: &twitter.DirectMessageEvent{
			Type: "message_create",
			Message: &twitter.DirectMessageEventMessage{
				Target: &twitter.DirectMessageTarget{
					RecipientID: user.IDStr,
				},
				Data: &twitter.DirectMessageData{
					Text: "this is a DM",
				},
			},
		},
	})

	logrus.Infof("dm: %+v, resp: %+v, err: %+v", dm, resp, err)
}
