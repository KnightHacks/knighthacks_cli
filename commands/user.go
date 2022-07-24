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
					if err := c.Load(configPath); err != nil {
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
			&cli.Command{
				Name:        "delete",
				Description: "deletes your currently logged in account",
				Action: func(context *cli.Context) error {
					configPath := context.Path("config")
					if err := c.Load(configPath); err != nil {
						return err
					}
					accessToken := c.Auth.Tokens.Access
					if len(accessToken) == 0 {
						return fmt.Errorf("you must first login to execute this command")
					}
					me, err := a.Delete(c, c.Auth.UserID)
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
