package kat

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	KAT_URL string = "https://kat.cr/json.php?q=%v&field=seeders&order=desc"
)

type KatQuery struct {
	Title    string
	Torrents []KatTorrent `json:"list"`
}

type KatTorrent struct {
	Title       string
	Category    string
	TorrentLink string
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

func getTorrents(query string) (KatQuery, error) {
	katUrl := fmt.Sprintf(KAT_URL, url.QueryEscape(query))
	jsonData, err := getData(katUrl)
	if err != nil {
		return KatQuery{}, err
	}

	var q KatQuery
	err = json.Unmarshal(jsonData, &q)
	if err != nil {
		return KatQuery{}, err
	}
	return q, nil
}

func GetUrl(query string) (string, error) {
	q, err := getTorrents(query)
	if err != nil {
		return "", err
	}

	if len(q.Torrents) == 0 {
		return "", errors.New("No torrent found for " + query)
	}

	urlParts := strings.Split(q.Torrents[0].TorrentLink, "?")
	url := strings.Replace(urlParts[0], "https", "http", 1)

	return url, nil
}
