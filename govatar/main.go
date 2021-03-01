package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/o1egl/govatar"
)

var version = "dev"

func main() {
	app := cli.NewApp()
	app.Name = "govatar"
	app.Usage = "Avatar generator service."
	app.Version = version
	app.Authors = []*cli.Author{
		{
			Name: "Oleg Lobanov",
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:      "generate",
			ArgsUsage: "<(male|m)|(female|f)>",
			Aliases:   []string{"g"},
			Usage:     "Generates random avatar",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "output,o",
					Value: "avatar.png",
					Usage: "Output file name",
				},
				&cli.StringFlag{
					Name:  "username,u",
					Value: "",
					Usage: "Username",
				},
			},
			Action: func(c *cli.Context) error {
				var g govatar.Gender
				var err error
				switch c.Args().First() {
				case "male", "m":
					g = govatar.MALE
				case "female", "f":
					g = govatar.FEMALE
				default:
					return fmt.Errorf("incorrect gender param. Run `govatar help generate`")
				}

				username := c.String("username")
				if username != "" {
					err = govatar.GenerateFileForUsername(g, username, c.String("output"))
				} else {
					err = govatar.GenerateFile(g, c.String("output"))
				}
				return err
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
