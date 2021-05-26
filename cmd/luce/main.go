package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	c := cli.NewApp()
	c.Name = "luce"
	c.Usage = "A collection of tools for luce projects"

	c.Commands = []cli.Command{
		{
			Name: "gen",
			Subcommands: []cli.Command{
				{
					Name:   "key",
					Action: keyCmd,
					Usage:  "Generate a random key",
				},
				{
					Name:   "rand",
					Action: randCmd,
					Usage:  "Generate a random number",
					Flags: []cli.Flag{
						&cli.IntFlag{
							Name:  "b",
							Value: 0,
							Usage: "Set size by bits",
						},
						&cli.Int64Flag{
							Name:  "n",
							Value: 100,
							Usage: "Set max size",
						},
					},
				},
				{
					Name:   "randbase64",
					Action: randBase64,
					Usage:  "Generate a random base64 value",
					Flags: []cli.Flag{
						&cli.IntFlag{
							Name:  "b",
							Value: 64,
							Usage: "Set the number of bytes",
						},
					},
				},
			},
		},
		{
			Name:    "socketclient",
			Action:  socketclient,
			Aliases: []string{"sc"},
			Usage:   "Connect to a unixsocket",
		},
	}

	c.Run(os.Args)
}
