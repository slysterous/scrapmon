package main

import (
	"fmt"
	cobraClient "github.com/slysterous/print-scrape/internal/cobra"
	cfg "github.com/slysterous/print-scrape/internal/config"
	printscrape "github.com/slysterous/print-scrape/internal/domain"
	file "github.com/slysterous/print-scrape/internal/file"
	phttp "github.com/slysterous/print-scrape/internal/http"
	"github.com/slysterous/print-scrape/internal/postgres"
	"log"
)

func main() {
	//init a db client.
	pgClient, err := postgres.NewClient(getDataSource(cfg.FromEnv()), cfg.FromEnv().MaxDBConnections)
	if err != nil {
		log.Fatalf("could not connect to DB, err: %v", err)
	}
	// init a file manager.
	fileManager := file.NewManager()
	if err != nil {
		log.Fatalf("could not get a file manager: %v", err)
	}

	storage := printscrape.Storage{
		Fm: fileManager,
		Dm: pgClient,
	}

	scrapper := phttp.NewProxyChainClient("127.0.0.1", "9050")

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
