package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"os"
	"text/tabwriter"
)

const (
	PACKAGE_NAME    = "redirector"
	PACKAGE_VERSION = "1.1.1"
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
			Name:   "export",
			Usage:  "export all mappings to a JSON document",
			Action: ExportMappingsAction,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "pretty,p",
					Usage: "print JSON with human-readable whitespace",
				},
			},
		},
		{
			Name:   "import",
			Usage:  "import mappings from a JSON document",
			Action: ImportMappingsAction,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file,f",
					Usage: "JSON file path",
				},
				cli.BoolFlag{
					Name:  "clear,x",
					Usage: "Clear all existing mappings before import",
				},
				cli.StringFlag{
					Name:  "comment,c",
					Usage: "Overwrite the comment for all imported mappings",
				},
			},
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
				cli.BoolFlag{
					Name:  "all,a",
					Usage: "permanently delete all mappings",
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

func ExportMappingsAction(c *cli.Context) error {
	cfg, err := GetConfig()
	if err != nil {
		return err
	}

	client := NewManagementClient(cfg)
	mappings, err := client.GetMappings()
	if err != nil {
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	if c.Bool("pretty") {
		enc.SetIndent("", "  ")
	}

	if err := enc.Encode(mappings); err != nil {
		return err
	}

	return nil
}

func ImportMappingsAction(c *cli.Context) error {
	r := os.Stdin
	if name := c.String("file"); name != "" {
		f, err := os.Open(name)
		if err != nil {
			return err
		}

		r = f
	}

	mappings := make([]*Mapping, 0)
	dec := json.NewDecoder(r)
	if err := dec.Decode(&mappings); err != nil {
		if serr, ok := err.(*json.SyntaxError); ok {
			return fmt.Errorf("Syntax error: %v at %v", err, serr.Offset)
		}

		return err
	}

	if len(mappings) == 0 {
		return fmt.Errorf("No mappings found in import document")
	}

	if comment := c.String("comment"); comment != "" {
		for _, m := range mappings {
			if m.Comment == "" {
				m.Comment = comment
			}
		}
	}

	for i, m := range mappings {
		if err := m.Validate(); err != nil {
			return fmt.Errorf("Validation error in mapping %v (key: %v): %v\n%v", i+1, m.Key, err)
		}
	}

	cfg, err := GetConfig()
	if err != nil {
		return err
	}

	client := NewManagementClient(cfg)
	if c.Bool("clear") {
		if err := client.RemoveAllMappings(); err != nil {
			return fmt.Errorf("Error removing existing mappings: %v", err)
		}
		// TODO: 404s will occur here until mappings are reimported
	}

	if err := client.AddMappings(mappings); err != nil {
		return fmt.Errorf("Error adding mappings: %v", err)
	}

	fmt.Printf("Imported %v mappings\n", len(mappings))
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
	cfg, err := GetConfig()
	if err != nil {
		return err
	}

	client := NewManagementClient(cfg)

	if c.Bool("all") {
		if err := client.RemoveAllMappings(); err != nil {
			return err
		}

		fmt.Printf("Removed all mappings\n")
	} else {
		// remove one
		m := &Mapping{
			Key: c.String("key"),
		}

		if m.Key == "" {
			return fmt.Errorf("Key not specified")
		}

		if err := client.RemoveMapping(m); err != nil {
			return err
		}

		fmt.Printf("Removed %v\n", m.Key)
	}

	return nil
}
