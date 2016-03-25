package main

import (
	"fmt"
	"log"
	"os"

	"github.com/JDevlieghere/GoTV/kat"
	"github.com/codegangsta/cli"
)

func downloadKat(episode *Episode, config Configuration, ch chan<- *Result) {
	ch <- download(episode, config.Quality, config.Directory, kat.GetUrl)
}

func run(config Configuration, verbose bool) {

	episodeCh := make(chan *Episode)
	resultCh := make(chan *Result)
	downloads := 0

	for _, title := range config.Series {
		go fetchLastEpisode(title, episodeCh)
	}

	queuedEpisode := config.dequeue()
	for queuedEpisode != nil {
		go downloadKat(queuedEpisode, config, resultCh)
		downloads++
		queuedEpisode = config.dequeue()
	}

	// Download new episodes
	for _, title := range config.Series {
		episode := <-episodeCh
		if episode != nil {
			go downloadKat(episode, config, resultCh)
			downloads++
			if verbose {
				log.Printf("Downloading %s", episode)
			}
		} else if verbose {
			log.Printf("No new episode found for %s", title)
		}
	}

	for i := 0; i < downloads; i++ {
		result := <-resultCh
		if result.Err != nil {
			config.enqueue(&result.Episode)
			if verbose {
				log.Print(result.Err)
			}
		}
	}

	config.save()
}

func main() {

	cfg := getConfig()
	app := cli.NewApp()

	app.Name = "GoTV"
	app.Usage = "Automatically download TV shows"
	app.Author = "Jonas Devlieghere"
	app.Email = "info@jonasdevlieghere.com"
	app.Version = "1.1.0"

	app.Commands = []cli.Command{
		{
			Name: "run",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "verbose, v",
					Usage: "show detailed output",
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
		{
			Name:  "clean",
			Flags: []cli.Flag{},
			Usage: "Empty download queue",
			Action: func(c *cli.Context) {
				var empty []Episode
				cfg.Queue = empty
				cfg.save()
			},
		},
	}

	app.Run(os.Args)
}
