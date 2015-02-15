package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

func newOAuthClient(ctx context.Context, config *oauth2.Config) *http.Client {
	path_tail := "/google-calendar-cli/token"
	var path string
	if cache_dir := os.Getenv("XDG_CACHE_HOME"); cache_dir != "" {
		path = cache_dir + path_tail
	} else {
		path = os.Getenv("HOME") + "/.cache" + path_tail
	}

	t, err := tokenFromFile(path)
	if err != nil {
		fmt.Println(err)
		t, err = tokenFromWeb(ctx, config)
		if err != nil {
			panic(err)
		}
		err := saveToken(path, t)
		if err != nil {
			panic(err)
		}
	}

	return config.Client(ctx, t)
}

func tokenFromFile(path string) (*oauth2.Token, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	t := new(oauth2.Token)
	err = json.NewDecoder(f).Decode(t)
	return t, err
}

func saveToken(path string, token *oauth2.Token) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(token)
}

// Ref: https://github.com/google/google-api-go-client/blob/master/examples/calendar.go
// Copyright 2014 The Go Authors. All rights reserved.
func tokenFromWeb(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
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

	return config.Exchange(ctx, code)
}
