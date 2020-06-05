package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"

	//cfg "github.com/slysterous/print-scrape/internal/config"
	printscrape "github.com/slysterous/print-scrape/internal/domain"
	//"github.com/slysterous/print-scrape/internal/postgres"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "print-scrape",
	Short: "Prntscr Scrapper",
	Long:  "A highly concurrent PrntScr Scrapper.",
}

// var purgeCmd = &cobra.Command{
// 	Use:   "purge",
// 	Short: "Purge postgres db of all data",
// 	Long:  "This command with purge the postgres database of all the prnt.sc data that are already scrapped",
// 	Run:   purgeFn,
// }

// var findCmd = &cobra.Command{
// 	Use:   "find",
// 	Short: "Searches for an already scrapped screenshot",
// 	Long:  "Searches postgres db for an already scrapped screenshot. Returns all data available for it.",
// 	Run:   findFn,
// }

// var fetchCmd = &cobra.Command{
// 	Use:   "fetch",
// 	Short: "Scrapes a specific screenshot",
// 	Long:  "Scrapes a specific screenshot based on a 6 digit alphanumeric code.",
// 	Run:   fetchFn,
// }

// var scrapeCmd = &cobra.Command{
// 	Use:   "scrape",
// 	Short: "Scrapes everything based on an algorithm",
// 	Run:   scrapeFn,
// }

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Nevermind this is a test",
	Run:   testFn,
}

func init() {
	// rootCmd.AddCommand(purgeCmd)
	// rootCmd.AddCommand(findCmd)
	// rootCmd.AddCommand(fetchCmd)
	// rootCmd.AddCommand(scrapeCmd)
	rootCmd.AddCommand(testCmd)

	testCmd.Flags().StringP("code", "c", "", "prntsc code")
	//findCmd.Flags().StringP("code", "c", "", "6 digit alphanumeric code to search against the database")
	//fetchCmd.Flags().StringP("code", "c", "", "6 digit alphanumeric code to scrape")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

//purgeFn is responsible for truncating the database and the filesystem
// func purgeFn(_ *cobra.Command, args []string) {
// 	pgClient, err := postgres.NewClient(getDataSource(cfg.FromEnv()), cfg.FromEnv().MaxDBConnections)
// 	if err != nil {
// 		log.Fatalf("could not connect to DB, err: %v", err)
// 	}

// 	fileClient, err := file.NewClient()
// 	if err != nil {
// 		log.Fatalf("could not get a file manager: %v", err)
// 	}

// 	storage := printscrape.Storage{
// 		fileClient,
// 		pgClient,
// 	}

// 	err = storage.Purge()
// 	if err !=nil {
// 		log.Fatalf("could not Purge Storage, err: %v", err)
// 	}
// }

// func findFn(_ *cobra.Command, args []string) {
// 	pg, err := postgres.NewClient(getDataSource(cfg.FromEnv()), cfg.FromEnv().MaxDBConnections)

// 	if err != nil {
// 		log.Fatalf("could not connect to DB, err: %v", err)
// 	}

// 	//find a scrap based on a 5 digit code
// 	ScreenshotGetter := printscrape.ScreenshotGetter(pg)

// 	ScreenshotGetter.GetByCode(args[0])
// }

// func fetchFn(_ *cobra.Command, args []string) {
// 	_, err := postgres.NewClient(getDataSource(cfg.FromEnv()), cfg.FromEnv().MaxDBConnections)
// 	if err != nil {
// 		log.Fatalf("could not connect to DB, err: %v", err)
// 	}

// 	//fetch a 5 digit screenshot from prnt.scr

// }

// func scrapeFn(_ *cobra.Command, args []string) {
// 	_, err := postgres.NewClient(getDataSource(cfg.FromEnv()), cfg.FromEnv().MaxDBConnections)
// 	if err != nil {
// 		log.Fatalf("could not connect to DB, err: %v", err)
// 	}

// 	//begin the async concurrent process of downloading files from prnt-scr

// }

func testFn(cmd *cobra.Command, args []string) {
	//get code
	_, err := cmd.Flags().GetString("code")

	if err != nil {
		log.Fatalf("error loading screenshot code, err: %v,err")
	}

	//url:=fmt.Sprintf("https://prnt.sc/%s",code)
	url := "https://prnt.sc"
	log.Printf("scraping page: %s\n", url)

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(doc.Text())

	// Find the review items
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		link := s.Find("src").Text()
		fmt.Printf("THIS: %s", link)
	})

	fmt.Println("SUCCESS")

}

func getDataSource(cfg printscrape.Config) string {
	user := cfg.DatabaseUser
	pass := cfg.DatabasePassword
	host := cfg.DatabaseHost
	port := cfg.DatabasePort
	name := cfg.DatabaseName

	return "host=" + host + " port=" + port + " user=" + user + " password=" + pass + " dbname=" + name + " sslmode=disable"
}
