package main

import (
	"github.com/KnightHacks/knighthacks_cli/model"
	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"os"
	"strings"
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
							code := RunRedirectServer(context.Context)

							loginPayload, err := api.Login(provider, code)
							if err != nil {
								return err
							}
							exists := loginPayload.AccountExists
							log.Printf("AccountExists=%v\n", exists)
							if exists {
								log.Printf("user=%v\n", *loginPayload.User)
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
							code := RunRedirectServer(context.Context)

							log.Println("Registering account now...")

							user := GetNewUserFromFlags(context)
							userId, err := api.Register(provider, code, user)
							if err != nil {
								return err
							}
							log.Printf("Created user with ID=%s", userId)

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
