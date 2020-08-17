package main

import (
	//"time"
	"fmt"
	"github.com/slysterous/print-scrape/internal/file"
	"github.com/slysterous/print-scrape/internal/postgres"
	"log"
	"os"
	"time"

	cfg "github.com/slysterous/print-scrape/internal/config"
	printscrape "github.com/slysterous/print-scrape/internal/domain"
	phttp "github.com/slysterous/print-scrape/internal/http"
	customNumber "github.com/slysterous/print-scrape/pkg/customnumber"
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
// 	Short: "Searches for an already scrapped ScreenShot",
// 	Long:  "Searches postgres db for an already scrapped ScreenShot. Returns all data available for it.",
// 	Run:   findFn,
// }

// var fetchCmd = &cobra.Command{
// 	Use:   "fetch",
// 	Short: "Scrapes a specific ScreenShot",
// 	Long:  "Scrapes a specific ScreenShot based on a 6 digit alphanumeric code.",
// 	Run:   fetchFn,
// }

var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Scrapes everything based on an algorithm",
	Run:   scrapeFn,
}

//var testCmd = &cobra.Command{
//	Use:   "test",
//	Short: "Nevermind this is a test",
//	Run:   testFn,
//}

func init() {
	// rootCmd.AddCommand(purgeCmd)
	// rootCmd.AddCommand(findCmd)
	// rootCmd.AddCommand(fetchCmd)
	rootCmd.AddCommand(scrapeCmd)
	//rootCmd.AddCommand(testCmd)

	//testCmd.Flags().StringP("code", "c", "", "prntsc code")
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
// 	ScreenShotGetter := printscrape.ScreenShotGetter(pg)

// 	ScreenShotGetter.GetByCode(args[0])
// }

// func fetchFn(_ *cobra.Command, args []string) {
// 	_, err := postgres.NewClient(getDataSource(cfg.FromEnv()), cfg.FromEnv().MaxDBConnections)
// 	if err != nil {
// 		log.Fatalf("could not connect to DB, err: %v", err)
// 	}

// 	//fetch a 5 digit ScreenShot from prnt.scr

// }

func scrapeFn(_ *cobra.Command, _ []string) {
	// init db client
	db, err := postgres.NewClient(getDataSource(cfg.FromEnv()), cfg.FromEnv().MaxDBConnections)
	if err != nil {
		log.Fatalf("could not connect to DB, err: %v", err)
	}
	// find the resume point, if any.
	lastCode, err := db.GetLatestCreatedScrapCode()
	if err != nil {
		log.Fatalf("could not get latest prnt.sc code, err: %v", err)
	}

	// create a screenShot item
	screenShot := printscrape.ScreenShot{
		CodeCreatedAt: time.Now(),
		RefCode:       createResumeCodeNumber(lastCode).String(),
		FileURI:       "",
	}

	// start saving item to db with status pending
	_, err = db.CreateScreenShot(screenShot)
	if err != nil {
		log.Fatalf("could not save scrap, err: %v", err)
	}

	torClient := phttp.NewClient("tor", "9051")

	// client fetch image url for specific code -- return url string
	//url, err := torClient.FetchScreenShotSourceLinkByCode(screenShot.RefCode)
	//if err != nil {
	//	log.Fatalf("could not fetch ScreenShot image link, err: %v",err)
	//}
	//
	//if url == nil {
	//	log.Fatalf("image link was not found")
	//}

	// update screenShot status to pending since the process has started.
	err = db.UpdateScreenShotStatusByCode(screenShot.RefCode,printscrape.StatusOngoing)
	if err !=nil {
		log.Fatalf("could not update scrap status, err: %v",err)
	}

	// create the url where the actual image will be saved.
	fileUrl:= fmt.Sprintf("%s/%s.jpg",cfg.FromEnv().ScreenShotStorageFolder,screenShot.RefCode)

	//save file to filesystem
	fileManager := file.NewManager()

	//download and save the image
	err = torClient.DownloadScreenShot(screenShot.RefCode,fileUrl,fileManager)
	if err !=nil {
		log.Fatalf("could not ")
	}

	screenShot.FileURI=fileUrl
	screenShot.Status = printscrape.StatusOngoing


	err = db.UpdateScreenShotByCode(screenShot)
	if err != nil {
		log.Fatalf("could not update scrap status and file url, err: %v",err)
	}

	//save file to filesystem
	fileManager := file.NewManager()

	err = fileManager.SaveImage(imageReader, fileUrl)
	if err != nil {
		log.Fatalf("could not save image to file system, err: %v",err)
	}

	err = db.UpdateScreenShotStatusByCode(screenShot.RefCode,printscrape.StatusSuccess)
	if err != nil {
		log.Fatalf("could not update scrap status to success, err: %v",err)
	}
}

//func testFn(cmd *cobra.Command, args []string) {
//
//	//start := time.Now()
//
//	number := customNumber.NewNumber(values, "150000")
//	fmt.Printf("initial number: %s \n", number.String())
//
//	for number.String() != "z00000" {
//		number.Increment()
//		fmt.Printf("initial number: %s \n", number.String())
//	}
//	//fmt.Println(time.Since(start))
//
//	// client := phttp.NewClient("tor", "9051")
//
//	// //todo design the way to find codes to use
//	// code := "gae309"
//	// pathToSave := "gae309.png"
//
//	// // client fetch image url for specific code -- return url string
//	// url, err := client.GetImageUrlByCode(code)
//	// if err != nil {
//	// 	log.Fatalf("BUMMER!")
//	// }
//
//	// //download the image by image url --return image (stream)
//	// imageReader, err := client.DownloadImage(url)
//
//	// //save image to file system
//	// fileManager := file.NewManager()
//	// err = fileManager.SaveImage(imageReader, pathToSave)
//	// if err != nil {
//	// 	log.Fatalf("BUMMER!2")
//	// }
//
//	// //save state for specific code to database.
//	// dbClient, err := postgres.NewClient(getDataSource(cfg.FromEnv()), cfg.FromEnv().MaxDBConnections)
//	// ScreenShot:=printscrape.ScreenShot{
//	// 	RefCode:code,
//	// 	FileURI:pathToSave,
//	// }
//	// _,err=dbClient.CreateScreenShot(ScreenShot)
//	// if err!=nil{
//	// 	log.Fatalf("ANTE GAMHSOU")
//	// }
//}

func getDataSource(cfg printscrape.Config) string {
	user := cfg.DatabaseUser
	pass := cfg.DatabasePassword
	host := cfg.DatabaseHost
	port := cfg.DatabasePort
	name := cfg.DatabaseName

	return "host=" + host + " port=" + port + " user=" + user + " password=" + pass + " dbname=" + name + " sslmode=disable"
}

func createResumeCodeNumber(code *string) customNumber.Number {
	// if no code was found then start from the beginning.
	if code == nil {
		return customNumber.NewNumber(printscrape.CustomNumberDigitValues, "000000")
	}

	number := customNumber.NewNumber(printscrape.CustomNumberDigitValues, *code)
	number.Increment()
	return number
}
