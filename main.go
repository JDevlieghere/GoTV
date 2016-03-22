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

func run(config config.Configuration) {

	episodeCh := make(chan *core.Episode)
	errorCh := make(chan error)

	series, err := core.ReadSeries(config.File)
	if err != nil {
		log.Fatal(err)
	}

	for _, title := range series {
		go core.FetchLastEpisode(title, episodeCh)
	}

	downloads := 0
	for i := 0; i < len(series); i++ {
		episode := <-episodeCh
		if episode != nil {
			go downloadKat(episode, config.Directory, errorCh)
			downloads++
		}
	}

	for i := 0; i < downloads; i++ {
		err = <-errorCh
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {

	app := cli.NewApp()

	app.Name = "GoTV"
	app.Usage = "Automatically download TV shows"
	app.Author = "Jonas Devlieghere"
	app.Version = "1.0.0"

	app.EnableBashCompletion = true

	app.Commands = []cli.Command{
		{
			Name:  "run",
			Flags: []cli.Flag{},
			Usage: "run GoTV",
			Action: func(c *cli.Context) {
				run(config.GetConfig())
			},
		},
		{
			Name: "config",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Usage: "File containing one tv show per line",
				},
				cli.StringFlag{
					Name:  "dir, d",
					Usage: "Path to download directory",
				},
			},
			Usage: "Change configuration",
			Action: func(c *cli.Context) {
				if c.IsSet("file") {
					file := c.String("file")
					config.SetFile(file)
				}
				if c.IsSet("dir") {
					dir := c.String("dir")
					config.SetDirectory(dir)
				}
				fmt.Printf("%s", config.GetConfig())
			},
		},
	}

	app.Run(os.Args)
}
