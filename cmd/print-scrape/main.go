package main

import (
	"container/ring"
	cfg "github.com/slysterous/print-scrape/internal/config"
	printscrape "github.com/slysterous/print-scrape/internal/domain"
	"github.com/slysterous/print-scrape/internal/file"
	phttp "github.com/slysterous/print-scrape/internal/http"
	"github.com/slysterous/print-scrape/internal/postgres"
	"log"
	"os"
	//"github.com/slysterous/print-scrape/internal/postgres"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "print-scrape",
	Short: "Prntscr Scrapper",
	Long:  "A highly concurrent PrntScr Scrapper.",
}

//a linked list of linked lists
var test [6]ring.Ring

//for i, ring := range test {
//	test[i] = ring.New(36)
//	for j,cell :=range test[i]{
//
//	}
//}

func newRing()*ring.Ring{
	r:=ring.New(36)
	values:=[]rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
	for _,rn :=range values{
		r.Value=rn
		r.Next()
	}
	r.Next()
	return r
}

// => 0 1 2 3 4 5 6 7 8 9 a b c d e f g h i j k l m n o p q r s t u v w x y z  arithmetic system

// var value = "0a9esd"

// value.increment() => 0a8ese

//==================================

// var value = "0a9esz"

// value.increment() => 0a8et0
//=====================================


//start scrapping
// 000000
// .
// .
// 00000z
// 000010





//
// node represents a single item in the list, where each item has a value and points to another item
type node struct {
	value interface{}
	next  *node
}

// List has all the items of the list beginning at first and ending at last
type List struct {
	first *node
	last  *node
	size  int
}

// New instantiates a new list and adds values to the list if any
func New(values ...interface{}) *List {
	list := &List{}
	if len(values) > 0 {
		list.Append(values...)
	}

	return list
}

// Append appends values to the list
func (list *List) Append(values ...interface{}) {
	for _, value := range values {
		newNode := &node{value: value}
		if list.size == 0 {
			list.first = newNode
			list.last = newNode
		} else {
			list.last.next = newNode
			list.last = newNode
		}
		list.size++
	}
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
	client := phttp.NewClient("tor", "9051")

	//todo design the way to find codes to use
	code := "gae309"
	pathToSave := "gae309.png"

	// client fetch image url for specific code -- return url string
	url, err := client.GetImageUrlByCode(code)
	if err != nil {
		log.Fatalf("BUMMER!")
	}

	//download the image by image url --return image (stream)
	imageReader, err := client.DownloadImage(url)

	//save image to file system
	fileManager := file.NewManager()
	err = fileManager.SaveImage(imageReader, pathToSave)
	if err != nil {
		log.Fatalf("BUMMER!2")
	}

	//save state for specific code to database.
	dbClient, err := postgres.NewClient(getDataSource(cfg.FromEnv()), cfg.FromEnv().MaxDBConnections)
	screenshot:=printscrape.Screenshot{
		RefCode:code,
		FileURI:pathToSave,
	}
	_,err=dbClient.CreateScrap(screenshot)
	if err!=nil{
		log.Fatalf("ANTE GAMHSOU")
	}
}

func getDataSource(cfg printscrape.Config) string {
	user := cfg.DatabaseUser
	pass := cfg.DatabasePassword
	host := cfg.DatabaseHost
	port := cfg.DatabasePort
	name := cfg.DatabaseName

	return "host=" + host + " port=" + port + " user=" + user + " password=" + pass + " dbname=" + name + " sslmode=disable"
}
