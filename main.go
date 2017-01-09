package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"os"
	"text/tabwriter"
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
			Name:   "ls",
			Usage:  "list all mappings",
			Action: ListMappingsAction,
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
				cli.BoolFlag{
					Name:  "permanent,p",
					Usage: "Redirect is permanent (301)",
				},
				cli.StringFlag{
					Name:  "comment,c",
					Usage: "Description of this redirection",
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

	defer rt.Close()

	go func() {
		serveManager(rt)
	}()

	return serve(rt)
}

func ListMappingsAction(c *cli.Context) error {
	cfg, err := GetConfig()
	if err != nil {
		return err
	}

	client := NewManagementClient(cfg)
	mappings, err := client.GetMappings()
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)
	fmt.Fprintln(w, "KEY\tDESTINATION\tPERMANENT\tCOMMENT")
	for _, m := range mappings {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", m.Key, m.Destination, m.Permanent, m.Comment)
	}
	w.Flush()

	return nil
}

func AddMappingAction(c *cli.Context) error {
	m := &Mapping{
		Key:         c.String("key"),
		Destination: c.String("dest"),
		Permanent:   c.Bool("permenant"),
		Comment:     c.String("comment"),
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

	fmt.Printf("Added %v\n", m)
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

	fmt.Printf("Removed %v\n", m.Key)
	return nil
}
