package commands

import (
	"fmt"
	"github.com/KnightHacks/knighthacks_cli/api"
	"github.com/KnightHacks/knighthacks_cli/config"
	"github.com/urfave/cli/v2"
	"log"
)

func GetUserCommand(a *api.Api, c *config.Config) *cli.Command {
	return &cli.Command{
		Name:    "user",
		Aliases: []string{"users", "u"},
		Usage:   "options relating to users",
		Subcommands: cli.Commands{
			&cli.Command{
				Name: "me",
				Action: func(context *cli.Context) error {
					configPath := context.Path("config")
					err := c.Load(configPath)
					if err != nil {
						return err
					}
					accessToken := c.Auth.Tokens.Access
					if len(accessToken) == 0 {
						return fmt.Errorf("you must first login to execute this command")
					}
					me, err := a.Me(c)
					if err != nil {
						return err
					}
					log.Printf("You are = %v\n", me)
					return nil
				},
			},
		},
	}
}
