package main

import (
	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	api := Api{Client: &http.Client{Timeout: time.Second * 10}}

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "endpoint",
				Value:       "http://localhost:4000/",
				DefaultText: "http://localhost:4000/",
				Usage:       "url to backend endpoint",
				Destination: &api.Endpoint,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "auth",
				Usage: "options relating to authentication",
				Subcommands: []*cli.Command{
					{
						Name:  "login",
						Usage: "uses kh login flow to login the user",
						Action: func(context *cli.Context) error {
							log.Println(api.Endpoint)
							// TODO: Implement provider dynamic-ability
							provider := "GITHUB"
							link, err := api.GetAuthRedirectLink(provider)
							if err != nil {
								return err
							}

							log.Printf("Opening %s in browser", link)
							err = browser.OpenURL(link)
							if err != nil {
								return err
							}
							code := RunRedirectServer()

							loginPayload, err := api.Login(provider, code)
							exists := loginPayload.AccountExists
							log.Printf("AccountExists=%v\n", exists)
							if exists {
								log.Printf("user=%v\n", *loginPayload.User)
							}
							return nil
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("Unable to run CLI, %s\n", err)
		return
	}

}
