package main

import (
	"github.com/KnightHacks/knighthacks_cli/api"
	"github.com/KnightHacks/knighthacks_cli/commands"
	"github.com/KnightHacks/knighthacks_cli/config"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	a := &api.Api{Client: &http.Client{Timeout: time.Second * 10}}

	c := &config.Config{}
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "endpoint",
				Value:       "http://localhost:4000/",
				DefaultText: "http://localhost:4000/",
				Usage:       "url to backend endpoint",
				Destination: &a.Endpoint,
			},
			&cli.PathFlag{Name: "config", Value: "config.yaml"},
		},
		Commands: getCommands(a, c),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("Unable to run CLI, %s\n", err)
		return
	}
}

func getCommands(api *api.Api, c *config.Config) []*cli.Command {
	return []*cli.Command{
		commands.GetAuthCommand(api, c),
		commands.GetUserCommand(api, c),
	}
}
