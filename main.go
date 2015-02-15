package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
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

	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:   "list",
			Action: toCommandFunc(svc, List),
		},
		{
			Name: "event",
			Subcommands: []cli.Command{
				{
					Name:   "add",
					Action: toCommandFunc(svc, AddEvent),
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "calendar-name, c",
							Value: "",
						},
						cli.StringFlag{
							Name:  "calendar-id, i",
							Value: "",
						},
						cli.StringFlag{
							Name:  "template, t",
							Value: "",
						},
						cli.StringFlag{
							Name: "Week-day, w",
						},
						cli.IntFlag{
							Name:  "next-week, n",
							Value: 0,
						},
					},
				},
			},
		},
	}
	app.Run(os.Args)
}

func AddEvent(svc *calendar.Service, ctx *cli.Context) {
}

func List(svc *calendar.Service, _ *cli.Context) {
	list, err := svc.CalendarList.List().Do()
	if err != nil {
		panic(err)
	}
	for _, i := range list.Items {
		fmt.Println(i.Summary)
	}
}

func toCommandFunc(svc *calendar.Service, f func(*calendar.Service, *cli.Context)) func(*cli.Context) {
	return func(c *cli.Context) {
		f(svc, c)
	}
}

func initCacheDir() (string, string) {
	path_tail := "/google-calendar-cli"
	var path string
	if cache_dir := os.Getenv("XDG_CACHE_HOME"); cache_dir != "" {
		path = cache_dir + path_tail
	} else {
		path = os.Getenv("HOME") + "/.cache" + path_tail
	}

	if _, err := os.Stat(path); err != nil {
		if err := os.MkdirAll(path, 0755); err != nil {
			panic(err)
		}
	}

	eventTemplatesPath := path + "/event_templates"
	if _, err := os.Stat(eventTemplatesPath); err != nil {
		if err := os.MkdirAll(eventTemplatesPath, 0755); err != nil {
			panic(err)
		}
	}

	return path, eventTemplatesPath
}

var CacheDirPath, EventTemplatesPath = initCacheDir()
