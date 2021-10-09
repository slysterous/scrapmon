package http

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/slysterous/scrapmon/internal/scrapmon"
	//"os"
)

//go:generate mockgen -destination mock/http.go -package http_mock . Reader,Downloader

// Downloader downloads files via http
type Downloader interface {
	Get(url string) (resp *http.Response, err error)
}

// Reader is well, a reader
type Reader interface {
	ReadAll(r io.Reader) ([]byte, error)
}

// Client represents an http client.
type Client struct {
	baseURL    string
	reader     Reader
	downloader Downloader
}

// NewClient returns a new http client.
func NewClient(baseUrl string, reader Reader, downloader Downloader) *Client {
	return &Client{
		baseURL:    baseUrl,
		reader:     reader,
		downloader: downloader,
	}
}

// ScrapeByCode scrapes a file.
func (c Client) ScrapeByCode(code, ext string) (scrapmon.ScrapedFile, error) {
	url := c.baseURL + code + "." + ext

	//Get the response bytes from the url
	response, err := c.downloader.Get(url)
	if err != nil {
		fmt.Printf("http: could not download image stream for url: %s, error %v \n", url, err)
		return scrapmon.ScrapedFile{}, fmt.Errorf("http: could not download image stream for url: %s, error %v", url, err)
	}

	defer response.Body.Close()

	if response.StatusCode == 404 || response.StatusCode == 302 {
		//fmt.Printf("NOT FOUND! STATUS: %d \n", response.StatusCode)
		return scrapmon.ScrapedFile{}, nil
	}

	//fmt.Printf("STATUS: %d \n", response.StatusCode)
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return scrapmon.ScrapedFile{}, fmt.Errorf("http: could not extract data from imagestream, err: %v", err)
	}

	contentType := response.Header.Get("Content-Type")
	imageType := strings.TrimLeft(contentType, "image/")

	//TODO should i keep this?
	if imageType == "f" {
		imageType = "gif"
	}

	ScrapedFile := scrapmon.ScrapedFile{
		Data: contents,
		Type: imageType,
		Code: code,
	}
	return ScrapedFile, nil
}
