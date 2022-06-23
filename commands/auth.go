package commands

import (
	"github.com/KnightHacks/knighthacks_cli/api"
	"github.com/KnightHacks/knighthacks_cli/config"
	"github.com/KnightHacks/knighthacks_cli/model"
	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
	"log"
	"strings"
)

func GetAuthCommand(a *api.Api, c *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "auth",
		Usage: "options relating to authentication",
		Subcommands: []*cli.Command{
			{
				Name:  "login",
				Usage: "uses kh login flow to login the user",
				Action: func(context *cli.Context) error {
					log.Println(a.Endpoint)
					// TODO: Implement provider dynamic-ability
					provider := "GITHUB"
					link, state, err := a.GetAuthRedirectLink(provider)
					if err != nil {
						return err
					}

					log.Printf("Opening %s in browser", link)
					err = browser.OpenURL(link)
					if err != nil {
						return err
					}
					code := api.RunRedirectServer(context.Context)

					loginPayload, err := a.Login(provider, code, state)
					if err != nil {
						return err
					}
					exists := loginPayload.AccountExists
					log.Printf("AccountExists=%v\n", exists)
					if exists {
						log.Printf("user=%v\n", *loginPayload.User)
						configPath := context.Path("config")
						err = c.Load(configPath)
						if err != nil {
							return err
						}
						c.Auth.Tokens.Access = *loginPayload.AccessToken
						c.Auth.Tokens.Refresh = *loginPayload.RefreshToken
						c.Auth.UserID = loginPayload.User.ID

						err = c.Save(configPath)
						if err != nil {
							return err
						}
						log.Printf("access=%s refresh=%s\n", c.Auth.Tokens.Access, c.Auth.Tokens.Refresh)
					}
					return nil
				},
			},
			{
				Name:  "register",
				Usage: "uses kh login flow to register a user",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "first-name",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "last-name",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "email",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "phone",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "pronouns",
						Usage:    "subjective/objective",
						Required: false,
					},
					&cli.IntFlag{
						Name:     "age",
						Required: false,
					},
				},
				Action: func(context *cli.Context) error {
					log.Println(a.Endpoint)
					// TODO: Implement provider dynamic-ability
					provider := "GITHUB"
					link, state, err := a.GetAuthRedirectLink(provider)
					if err != nil {
						return err
					}

					log.Printf("Opening %s in browser", link)
					err = browser.OpenURL(link)
					if err != nil {
						return err
					}
					code := api.RunRedirectServer(context.Context)

					loginPayload, err := a.Login(provider, code, state)
					if err != nil {
						return err
					}
					exists := loginPayload.AccountExists
					log.Printf("AccountExists=%v\n", exists)
					if exists {
						log.Printf("user=%v\n", *loginPayload.User)
					} else {
						log.Println("Registering account now...")

						log.Printf("loginPayload=%v\n", *loginPayload)

						log.Printf("EncryptedOAuthAccessToken=%s\n", *loginPayload.EncryptedOAuthAccessToken)

						user := GetNewUserFromFlags(context)
						registrationPayload, err := a.Register(provider, *loginPayload.EncryptedOAuthAccessToken, user)
						if err != nil {
							return err
						}
						log.Printf("Created user with ID=%s", registrationPayload.User.ID)
						log.Printf("AccessToken=%s\n", registrationPayload.AccessToken)
						log.Printf("RefreshToken=%s\n", registrationPayload.RefreshToken)
					}
					return nil
				},
			},
		},
	}
}

func GetNewUserFromFlags(context *cli.Context) model.NewUser {
	user := model.NewUser{
		FirstName:   context.String("first-name"),
		LastName:    context.String("last-name"),
		Email:       context.String("email"),
		PhoneNumber: context.String("phone"),
	}
	age := context.Int("age")
	if age != 0 {
		user.Age = &age
	}
	pronounsString := context.String("pronouns")
	if len(pronounsString) > 0 {
		pronounSlice := strings.Split(pronounsString, "/")
		if len(pronounSlice) == 2 {
			user.Pronouns = &model.PronounsInput{
				Subjective: pronounSlice[0],
				Objective:  pronounSlice[1],
			}
		} else {
			log.Fatalf("Incorrectly enter pronouns, %s is not valid, should be similar to he/him, she/her, they/them, etc\n", pronounsString)
		}
	}
	return user
}
