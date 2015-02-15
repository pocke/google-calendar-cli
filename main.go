package main

import (
	"fmt"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

func main() {
	config := &oauth2.Config{
		ClientID:     "406595811286-kntrl98g7aujgsesob2oi6sfelsqlh4n.apps.googleusercontent.com",
		ClientSecret: "ypjZIJwF1rnfkn-MqIzQPa3i",
		Endpoint:     google.Endpoint,
		Scopes:       []string{calendar.CalendarScope},
	}
	ctx := context.Background()

	c := newOAuthClient(ctx, config)

	svc, err := calendar.New(c)
	if err != nil {
		panic(err)
	}

	List(svc)
}

// XXX: Type of event
// func AddEvent(svc *calendar.Service, cal string, event Event) {
// }

func List(svc *calendar.Service) {
	list, err := svc.CalendarList.List().Do()
	if err != nil {
		panic(err)
	}
	for _, i := range list.Items {
		fmt.Println(i.Summary)
	}
}
