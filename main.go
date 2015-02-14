package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

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

// Ref: https://github.com/google/google-api-go-client/blob/master/examples/calendar.go
// Copyright 2014 The Go Authors. All rights reserved.
// TODO: save token
func newOAuthClient(ctx context.Context, config *oauth2.Config) *http.Client {
	ch := make(chan string)
	randState := fmt.Sprintf("st%d", time.Now().UnixNano())
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/favicon.ico" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.FormValue("state") != randState {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if code := r.FormValue("code"); code != "" {
			w.Write([]byte("<h1>Success</h1>Authorized."))
			w.(http.Flusher).Flush()
			ch <- code
			return
		}
		// no code
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	config.RedirectURL = ts.URL
	authURL := config.AuthCodeURL(randState)
	fmt.Println("Access to %s", authURL)
	code := <-ch

	token, err := config.Exchange(ctx, code)
	if err != nil {
		panic(err)
	}
	return config.Client(ctx, token)
}
