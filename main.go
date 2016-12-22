package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"os"
)

const (
	PACKAGE_NAME    = "redirector"
	PACKAGE_VERSION = "1.0.0"
)

func main() {
	app := cli.NewApp()
	app.Name = PACKAGE_NAME
	app.Version = PACKAGE_VERSION
	app.Usage = "Simple HTTP server to map old URLs to new URLs"
	app.Commands = []cli.Command{
		{
			Name:   "serve",
			Usage:  "start HTTP server",
			Action: ServeAction,
		},
		{
			Name:   "add",
			Usage:  "add a mapping",
			Action: AddMappingAction,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "key,k",
					Usage: "key that identifies this redirect",
				},
				cli.StringFlag{
					Name:  "dest,d",
					Usage: "URL to redirect to",
				},
			},
		},
		{
			Name:   "rm",
			Usage:  "remove a mapping",
			Action: RemoveMappingAction,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "key,k",
					Usage: "key that identifies this redirect",
				},
			},
		},
	}

	app.Run(os.Args)
}

func ServeAction(c *cli.Context) error {
	rt, err := NewRuntime()
	if err != nil {
		return err
	}

	defer rt.Database.Close()

	go func() {
		serveManager(rt)
	}()

	return serve(rt)
}

func AddMappingAction(c *cli.Context) error {
	m := &Mapping{
		Key:         c.String("key"),
		Destination: c.String("dest"),
	}

	if m.Key == "" {
		return fmt.Errorf("Key not specified")
	}

	if m.Destination == "" {
		return fmt.Errorf("Destination URL not specified")
	}

	cfg, err := GetConfig()
	if err != nil {
		return err
	}

	client := NewManagementClient(cfg)
	if err := client.AddMapping(m); err != nil {
		return err
	}

	fmt.Printf("OK\n")
	return nil
}

func RemoveMappingAction(c *cli.Context) error {
	m := &Mapping{
		Key: c.String("key"),
	}

	if m.Key == "" {
		return fmt.Errorf("Key not specified")
	}

	cfg, err := GetConfig()
	if err != nil {
		return err
	}

	client := NewManagementClient(cfg)
	if err := client.RemoveMapping(m); err != nil {
		return err
	}

	fmt.Printf("OK\n")
	return nil
}
