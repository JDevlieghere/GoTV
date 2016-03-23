package main

import (
	"fmt"
	"log"
	"os"

	"github.com/JDevlieghere/GoTV/config"
	"github.com/JDevlieghere/GoTV/core"
	"github.com/JDevlieghere/GoTV/kat"

	"github.com/codegangsta/cli"
)

func downloadKat(episode *core.Episode, dir string, ch chan<- error) {
	ch <- core.Download(episode, dir, kat.GetUrl)
}

func run(config config.Configuration, verbose bool) {

	episodeCh := make(chan *core.Episode)
	errorCh := make(chan error)

	for _, title := range config.Series {
		go core.FetchLastEpisode(title, episodeCh)
	}

	downloads := 0
	for _, title := range config.Series {
		episode := <-episodeCh
		if episode != nil {
			go downloadKat(episode, config.Directory, errorCh)
			downloads++
			if verbose {
				log.Printf("Downloading %s", episode)
			}
		} else if verbose {
			log.Printf("No new episode found for %s", title)
		}
	}

	for i := 0; i < downloads; i++ {
		err := <-errorCh
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {

	cfg := config.Get()
	app := cli.NewApp()

	app.Name = "GoTV"
	app.Usage = "Automatically download TV shows"
	app.Author = "Jonas Devlieghere"
	app.Email = "info@jonasdevlieghere.com"
	app.Version = "1.0.0"

	app.Commands = []cli.Command{
		{
			Name: "run",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "verbose, v",
					Usage: "show detailled output",
				},
			},
			Usage: "run GoTV",
			Action: func(c *cli.Context) {
				run(cfg, c.Bool("verbose"))
			},
		},
		{
			Name:  "info",
			Flags: []cli.Flag{},
			Usage: "Show configuration info",
			Action: func(c *cli.Context) {
				fmt.Println(cfg)
			},
		},
	}

	app.Run(os.Args)
}
