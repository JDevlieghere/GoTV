package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type GetUrl func(query string) (string, error)

type Episode struct {
	Show    string
	Season  int
	Episode int
	Title   string
}

type Result struct {
	Episode Episode
	Err     error
}

func (e Episode) String() string {
	return fmt.Sprintf("%s S%02dE%02d", e.Show, e.Season, e.Episode)
}

func getData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func getFile(url string, fileName string) error {
	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("User-Agent", "GoTV")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func download(episode *Episode, quality string, dir string, getUrl GetUrl) *Result {
	query := fmt.Sprintf("%s", episode)
	url, err := getUrl(query + " " + quality)
	if err != nil {
		return &Result{*episode, err}
	}
	fileName := strings.Replace(query, " ", ".", -1)
	filePath := filepath.Join(dir, fileName+".torrent")
	err = getFile(url, filePath)
	if err == nil {
		log.Println("Downloaded", episode, "to", filePath)
	}
	return &Result{*episode, err}
}
