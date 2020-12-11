package main

import (
	"fmt"
	"github.com/joho/godotenv"
	cobraClient "github.com/slysterous/print-scrape/internal/cobra"
	"github.com/slysterous/print-scrape/internal/config"
	cfg "github.com/slysterous/print-scrape/internal/config"
	file "github.com/slysterous/print-scrape/internal/file"
	phttp "github.com/slysterous/print-scrape/internal/http"
	"github.com/slysterous/print-scrape/internal/postgres"
	printscrape "github.com/slysterous/print-scrape/internal/printscrape"
	"log"
)

func main() {

	//load env.
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("could not load env file")
	}

	// fetch the config from env variables.
	config := config.FromEnv()

	pgClient, err := postgres.NewClient(getDataSource(cfg.FromEnv()), cfg.FromEnv().MaxDBConnections)
	if err != nil {
		log.Fatalf("could not connect to DB, err: %v", err)
	}

	// init a file manager.
	fileManager := file.NewManager(config.ScreenShotStorageFolder)
	if err != nil {
		log.Fatalf("could not get a file manager: %v", err)
	}

	//combine db and filestorage into generic storage.
	storage := printscrape.Storage{
		Fm: fileManager,
		Dm: pgClient,
	}

	//scrapper := phttp.NewProxyChainClient("127.0.0.1", "9050")
	scrapper := phttp.NewClient()

	commandManager := printscrape.CommandManager{
		Storage:  storage,
		Scrapper: scrapper,
	}

	cobraClient := cobraClient.NewClient(commandManager)
	cobraClient.RegisterStartCommand()
	cobraClient.RegisterPurgeCommand()

	if err := cobraClient.Execute(); err != nil {
		log.Fatalf("execution failed, err: %v", err)
	}

	fmt.Println("Execution has completed Successfuly!")
}

func getDataSource(cfg printscrape.Config) string {
	user := cfg.DatabaseUser
	pass := cfg.DatabasePassword
	host := cfg.DatabaseHost
	port := cfg.DatabasePort
	name := cfg.DatabaseName

	return "host=" + host + " port=" + port + " user=" + user + " password=" + pass + " dbname=" + name + " sslmode=disable"
}
