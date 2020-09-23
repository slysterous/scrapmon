package http

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/slysterous/print-scrape/internal/domain"
	"strings"
	//	printscrape "github.com/slysterous/print-scrape/internal/domain"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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

// ScrapeImageByCode fetches an prnt.sc image stream an image type and an error.
func (c Client) ScrapeImageByCode(code string) (domain.ScrapedImage, error) {

	url := "https://i.imgur.com/" + code + ".png"

	fmt.Printf("URL: %s \n", url)

	//Get the response bytes from the url
	response, err := c.httpClient.Get(url)
	if err != nil {
		fmt.Printf("http: could not download image stream for url: %s, error %v \n", url, err)
		return domain.ScrapedImage{}, fmt.Errorf("http: could not download image stream for url: %s, error %v", url, err)
	}

	defer response.Body.Close()

	fmt.Printf("STATUS: %d ", response.StatusCode)

	if response.StatusCode == 404 || response.StatusCode == 302 {
		fmt.Printf("NOT FOUND! STATUS: %d ", response.StatusCode)
		return domain.ScrapedImage{}, nil
	}

	contents, err := ioutil.ReadAll(response.Body)

	//bodyString := string(contents)
	//fmt.Printf("RESPONSE: %v",bodyString)

	if err != nil {
		fmt.Printf("ERROR2 \n")
		return domain.ScrapedImage{}, fmt.Errorf("http: could not extract data from imagestream, err: %v", err)
	}

	contentType := response.Header.Get("Content-Type")
	imageType := strings.TrimLeft(contentType, "image/")
	scrapedImage := domain.ScrapedImage{
		Data: contents,
		Type: imageType,
		Code: code,
	}
	return scrapedImage, nil
}

// // GetImageUrlByCode fetches a ScreenShot's actual img url from a code.
// func (c Client) GetImageUrlByCode(code string) (string, error) {
// 	return "", nil
// }

//func (c Client) DownloadImage(url string) (io.Reader, error) {
//	fetchScreenShotSourceLinkByCode
//
//	return
//}

//func (c Client) Get(requestURL string) (io.ReadCloser, error) {
//	randomUA := RandomUserAgent()
//
//	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	// Set headers
//	req.Header.Set("User-Agent", randomUA)
//	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
//	req.Header.Set("cookie", "lang=en")
//
//	resp, err := c.httpClient.Do(req)
//	if err != nil {
//		return nil, fmt.Errorf("http: getting from url %s error: %v", requestURL, err)
//	}
//
//	return resp.Body, nil
//}

// func (c *Client) fetchScreenShotSourceLinkByCode(code string) (*string, error) {
// 	bow := surf.NewBrowser()
// 	err := bow.Open("https://prnt.sc/" + code)
// 	if err != nil {
// 		return nil, fmt.Errorf("http: could not fetch Screenshot image source link: %v", err)
// 	}

// 	var screenShotUrl string

// 	bow.Dom().Find("#ScreenShot-image").Each(func(i int, selection *goquery.Selection) {
// 		imgURL := selection.AttrOr(`src`, ``)
// 		if imgURL != "" {
// 			screenShotUrl = imgURL
// 			return
// 		}
// 	})

// 	//url was not found
// 	if screenShotUrl == "" {
// 		return nil, nil
// 	}

// 	return &screenShotUrl, nil
// }

// // DownloadImage attemps to download an image that belongs to a 6digit code in prnt.sc.
// func (c *Client) DownloadScreenShot(code string, filepath string, manager printscrape.FileManager) error {

// 	imgURL, err := c.fetchScreenShotSourceLinkByCode(code)
// 	if err != nil {
// 		return fmt.Errorf("http: could not get image screenshot source, err: %v", err)
// 	}
// 	if imgURL == nil {
// 		return nil
// 	}

// 	//Get the response bytes from the url
// 	response, err := http.Get(*imgURL)
// 	if err != nil {
// 		return fmt.Errorf("http: could not get image screenshot, err: %v", err)
// 	}
// 	defer response.Body.Close()

// 	err = SaveFile(filepath, &response.Body)
// 	if err != nil {
// 		return err
// 	}

// 	//randomUA := GenerateRandomUserAgent()
// 	//
// 	//requestURL:="https://prnt.sc/aaaaab"
// 	////doc, err := htmlquery.LoadURL("https://prnt.sc/aaaaab")
// 	////
// 	////if err != nil {
// 	////	return err
// 	////}
// 	//
// 	//req, err := http.NewRequest(http.MethodGet, requestURL, nil)
// 	//if err != nil {
// 	//	log.Fatalln(err)
// 	//}
// 	//
// 	//// Set headers
// 	//req.Header.Set("User-Agent", randomUA)
// 	//req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
// 	//req.Header.Set("cookie", "lang=en")
// 	//
// 	//resp, err := c.httpClient.Do(req)
// 	//
// 	//
// 	//	bodyBytes, err := ioutil.ReadAll(resp.Body)
// 	//	if err != nil {
// 	//		log.Fatal(err)
// 	//	}
// 	//	bodyString := string(bodyBytes)
// 	//	fmt.Printf("RESPONSE: %v",bodyString)
// 	//
// 	//
// 	//if err != nil {
// 	//	return nil, fmt.Errorf("http: getting from url %s error: %v", requestURL, err)
// 	//}
// 	//
// 	//bodyReader:=resp.Body
// 	//
// 	//defer func(bodyReader io.ReadCloser) {
// 	//	errC := bodyReader.Close()
// 	//	if errC != nil {
// 	//		log.Fatal("cargr: closing bodyReader")
// 	//	}
// 	//}(bodyReader)
// 	//
// 	////doc, err := goquery.NewDocumentFromReader(bodyReader)
// 	//
// 	//
// 	return nil
// }

// func SaveFile(filePath string, body *io.ReadCloser) error {
// 	//Create a empty file
// 	file, err := os.Create(filePath)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	//Write the bytes to the file
// 	_, err = io.Copy(file, *body)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func DownloadAndSaveFile(URL, filePath string) error {
// 	//Get the response bytes from the url
// 	response, err := http.Get(URL)
// 	if err != nil {
// 	}
// 	defer response.Body.Close()

// 	//Create a empty file
// 	file, err := os.Create(filePath)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	//Write the bytes to the file
// 	_, err = io.Copy(file, response.Body)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
