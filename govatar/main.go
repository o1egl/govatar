package main

import (
	"fmt"
	"log"
	"os"

	"github.com/o1egl/govatar"
	"github.com/urfave/cli"
)

var version = "dev"

func main() {
	app := cli.NewApp()
	app.Name = "govatar"
	app.Usage = "Avatar generator service."
	app.Version = version
	app.Author = "Oleg Lobanov"
	app.Commands = []cli.Command{
		{
			Name:      "generate",
			ArgsUsage: "<(male|m)|(female|f)>",
			Aliases:   []string{"g"},
			Usage:     "Generates random avatar",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output,o",
					Value: "avatar.png",
					Usage: "Output file name",
				},
				cli.StringFlag{
					Name:  "username,u",
					Value: "",
					Usage: "Username",
				},
			},
			Action: func(c *cli.Context) {
				var g govatar.Gender
				var err error
				switch c.Args().First() {
				case "male", "m":
					g = govatar.MALE
				case "female", "f":
					g = govatar.FEMALE
				default:
					fmt.Println("Incorrect gender param. Run `govatar help generate`")
					os.Exit(1)
				}

				username := c.String("username")
				if username != "" {
					err = govatar.GenerateFileForUsername(g, username, c.String("output"))
				} else {
					err = govatar.GenerateFile(g, c.String("output"))
				}
				if err != nil {
					log.Fatal(err)
				}
			},
		},
	}
	app.Run(os.Args)
}
