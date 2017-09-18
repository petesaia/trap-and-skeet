package main

import (
	"os"

	"github.com/urfave/cli"
)

// Identifier for containers that are built. This allows for control.
const Identifier string = "_TaS_"

// TargetResponse is what's returned to the user once they spawn a new target
// that they want to abuse..
type TargetResponse struct {
	IPAddress string `json:"ipAddress"`
	Port      int    `json:"port"`
	Example   string `json:"example"`
}

// The CLI prompt.
func main() {
	app := cli.NewApp()
	app.Name = "Trap and Skeet"
	app.Usage = "Use to create containers that are meant to be abused."
	app.Version = "1.0.0"
	app.Copyright = "(c) 2017 Lev Interactive"
	app.Commands = []cli.Command{
		{
			Name:    "fire",
			Aliases: []string{"create"},
			Usage:   "This will fire off an ephemeral container.",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "docker-image", Value: "krlmlr/debian-ssh"},
				cli.StringFlag{Name: "ssh-user", Value: "root"},
				cli.StringFlag{Name: "ssh-password", Value: ""},
				cli.StringFlag{Name: "ssh-private-key-file", Value: ""},
				cli.StringFlag{Name: "ssh-port", Value: "2020"},
			},
			Action: func(c *cli.Context) error {
				Spawn(&ConfigureInput{
					User:       c.String("ssh-user"),
					Port:       c.String("ssh-port"),
					Pass:       c.String("ssh-password"),
					Image:      c.String("docker-image"),
					SSHKeyFile: c.String("ssh-private-key-file"),
				})
				return nil
			},
		},
		{
			Name:    "destroy",
			Aliases: []string{"remove"},
			Usage:   "This will destroy all containers were created.",
			Action: func(c *cli.Context) error {
				Remove()
				return nil
			},
		},
	}

	app.Run(os.Args)
}
