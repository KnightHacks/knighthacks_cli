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
							link, err := api.GetAuthRedirectLink("GITHUB")
							if err != nil {
								return err
							}

							log.Printf("Opening %s in browser", link)
							err = browser.OpenURL(link)
							if err != nil {
								return err
							}
							code := RunRedirectServer()

							api.Login()
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
