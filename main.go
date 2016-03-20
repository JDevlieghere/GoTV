package main

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	API_KEY     string = ""
	SERIES_URL  string = "http://thetvdb.com/api/GetSeries.php?seriesname=%v"
	EPISODE_URL string = "http://thetvdb.com/api/GetEpisodeByAirDate.php?apikey=%s&seriesid=%d&airdate=%v"
)

type TvDBEpisodeQuery struct {
	Episode TvDBEpisode
}

type TvDBEpisode struct {
	SeasonNumber  int
	EpisodeNumber int
	EpisodeName   string
}

type TvDBSeriesQuery struct {
	Series []TvDBSeries
}

type TvDBSeries struct {
	SeriesId   int `xml:"seriesid"`
	SeriesName string
}

type Episode struct {
	Show    string
	Season  int
	Episode int
	Title   string
}

func getXml(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	} else {
		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return contents, nil
	}
}

func getSeries(name string) (TvDBSeries, error) {
	seriesUrl := fmt.Sprintf(SERIES_URL, url.QueryEscape(name))
	xmlData, err := getXml(seriesUrl)
	if err != nil {
		return TvDBSeries{}, err
	}

	var q TvDBSeriesQuery
	err = xml.Unmarshal(xmlData, &q)
	if err != nil {
		return TvDBSeries{}, err
	}

	if len(q.Series) <= 0 {
		return TvDBSeries{}, errors.New("Could not find series with name " + name)
	} else {
		return q.Series[0], nil
	}
}

func getEpisode(series TvDBSeries, date string) (Episode, error) {
	episodeUrl := fmt.Sprintf(EPISODE_URL, url.QueryEscape(API_KEY), series.SeriesId, url.QueryEscape(date))
	xmlData, err := getXml(episodeUrl)
	if err != nil {
		return Episode{}, err
	}

	if bytes.Contains(xmlData, []byte("Error")) {
		return Episode{}, errors.New("No episode found for " + series.SeriesName + " on " + date)
	}

	var q TvDBEpisodeQuery
	err = xml.Unmarshal(xmlData, &q)
	if err != nil {
		return Episode{}, err
	}

	return Episode{
		Show:    series.SeriesName,
		Season:  q.Episode.SeasonNumber,
		Episode: q.Episode.EpisodeNumber,
		Title:   q.Episode.EpisodeName,
	}, nil
}

func lastEpisode(name string) (Episode, error) {
	series, err := getSeries(name)
	if err != nil {
		return Episode{}, err
	}
	date := time.Now().Format("2006-01-02")
	return getEpisode(series, date)
}

func fetchLastEpisode(name string, ch chan<- string) {
	episode, err := lastEpisode(name)
	if err != nil {
		ch <- fmt.Sprintf("%s", err)
		return
	}
	ch <- fmt.Sprintf("%s", episode)
}

func (e Episode) String() string {
	show := strings.Replace(e.Show, " ", ".", -1)
	title := strings.Replace(e.Title, " ", ".", -1)
	return fmt.Sprintf("%s.S%02dE%02d.%s", show, e.Season, e.Episode, title)
}

func main() {

	if len(os.Args) <= 1 {
		fmt.Println("Usage: GoTV <filepath>")
		return
	}

	ch := make(chan string)

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	entries := 0
	for scanner.Scan() {
		go fetchLastEpisode(scanner.Text(), ch)
		entries++
	}

	for i := 0; i < entries; i++ {
		fmt.Println(<-ch)
	}
}
