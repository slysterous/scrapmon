package http

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/slysterous/scrapmon/internal/scrapmon"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	//"os"
)

// Client represents an http client.
type Client struct {
	httpClient *http.Client
}

// NewClient returns a new http client.
func NewClient() *Client {

	httpClient := &http.Client{
		Transport: &http.Transport{},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	var err error
	httpClient.Transport, err = NewRoundTripper(httpClient.Transport)

	if err != nil {
		log.Fatalf("http: error creating http roundtripper transport, err: %v", err)
	}

	return &Client{
		httpClient: httpClient,
	}
}

// NewProxyChainClient returns a new http client that utilizes a proxy chain.
func NewProxyChainClient(host, port string) *Client {

	torProxyString := fmt.Sprintf("socks5://%s:%s", host, port)
	//torProxyString := fmt.Sprintf("socks5://%s:%s", "127.0.0.1", "9050")
	torProxyURL, err := url.Parse(torProxyString)
	if err != nil {
		log.Fatal("http: error parsing Tor proxy URL:", torProxyString, ". ", err)
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(torProxyURL),
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	httpClient.Transport, err = NewRoundTripper(httpClient.Transport)

	if err != nil {
		log.Fatalf("http: error creating http roundtripper transport, err: %v", err)
	}

	return &Client{
		httpClient: httpClient,
	}
}

// scrapeScreenShotURLByCode fetches a pnt.sc image actual url.
func (c Client) scrapeScreenShotURLByCode(code string) (*string, error) {

	requestURL := "https://prnt.sc/" + code

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("htto: creating a new get request, error: %v", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http: fetching url image from: %s ,error: %v", requestURL, err)
	}

	bodyReader := resp.Body

	defer func(bodyReader io.ReadCloser) {
		errC := bodyReader.Close()
		if errC != nil {
			log.Fatal("http: closing bodyReader")
		}
	}(bodyReader)

	var screenShotURL string

	doc, err := goquery.NewDocumentFromReader(bodyReader)
	doc.Find("#screenshot-image").Each(func(i int, selection *goquery.Selection) {
		imgURL := selection.AttrOr(`src`, ``)
		if imgURL != "" {
			screenShotURL = imgURL
			return
		}
	})

	//url was not found
	if screenShotURL == "" {
		return nil, fmt.Errorf("http: could not find screenShotUrl")
	}

	return &screenShotURL, nil
}

// ScrapeByCode fetches an prnt.sc image stream an image type and an error.
func (c Client) ScrapeByCode(code string) (scrapmon.ScrapedFile, error) {

	url := "https://i.imgur.com/" + code + ".png"

	//fmt.Printf("URL: %s \n", url)

	//Get the response bytes from the url
	response, err := c.httpClient.Get(url)
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

	//bodyString := string(contents)
	//fmt.Printf("RESPONSE: %v",bodyString)

	if err != nil {
		//fmt.Printf("ERROR2 \n")
		return scrapmon.ScrapedFile{}, fmt.Errorf("http: could not extract data from imagestream, err: %v", err)
	}

	contentType := response.Header.Get("Content-Type")
	imageType := strings.TrimLeft(contentType, "image/")

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