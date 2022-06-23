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
					Name:   "r32",
					Action: rand32,
					Usage:  "Generate a random uint32",
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
			Aliases: []string{"sc"},
			Action:  socketclient,
			Usage:   "Connect to a unixsocket",
		},
		{
			Name:    "filevars",
			Aliases: []string{"fv"},
			Action:  filevars,
			Usage:   "generates a go file with input files set to variables",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "o",
					Usage: "Set the name of the output file",
					Value: "filevars",
				},
				&cli.StringFlag{
					Name:  "p",
					Usage: "package name",
				},
			},
		},
	}

	c.Run(os.Args)
}
