package main

import (
	"fmt"
	"github.com/joho/godotenv"
	cobraClient "github.com/slysterous/scrapmon/internal/cobra"
	file "github.com/slysterous/scrapmon/internal/file"
	phttp "github.com/slysterous/scrapmon/internal/http"
	"github.com/slysterous/scrapmon/internal/postgres"
	scrapmon "github.com/slysterous/scrapmon/internal/scrapmon"
	config "github.com/slysterous/scrapmon/internal/config"
	"log"
)

func main() {

	//load env.
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("could not load env file")
	}

	// fetch the config from env variables.
	conf := config.FromEnv()

	// init database manager
	pgClient,err := postgres.NewClient(getDataSource(conf), conf.MaxDBConnections)
	if err != nil {
		log.Fatalf("could not connect to DB, err: %v", err)
	}
	defer pgClient.DB.Close()

	// init a file manager.
	fileManager := file.NewManager(conf.ScrapStorageFolder)

	//combine db and filestorage into generic storage.
	storage := scrapmon.Storage{
		Fm: fileManager,
		Dm: pgClient,
	}

	//TODO fix the TOR client
	//scrapper := phttp.NewProxyChainClient("127.0.0.1", "9050")
	scrapper := phttp.NewClient()

	commandManager := scrapmon.CommandManager{
		Storage:  storage,
		Scrapper: scrapper,
	}

	cobraC := cobraClient.NewClient()

	startCommand,err:=cobraC.NewStartCommand(commandManager.StartCommand)
	if err!=nil{
		log.Fatalf("could not register start command, err: %v",err)
	}
	purgeCommand:=cobraC.NewPurgeCommand(commandManager.PurgeCommand)

	cobraC.RegisterCommand(startCommand)
	cobraC.RegisterCommand(purgeCommand)

	if err := cobraC.Execute(); err != nil {
		log.Fatalf("execution failed, err: %v", err)
	}

	fmt.Println("Execution has completed Successfuly!")
}

func getDataSource(cfg scrapmon.Config) string {
	user := cfg.DatabaseUser
	pass := cfg.DatabasePassword
	host := cfg.DatabaseHost
	port := cfg.DatabasePort
	name := cfg.DatabaseName

	return "host=" + host + " port=" + port + " user=" + user + " password=" + pass + " dbname=" + name + " sslmode=disable"
}
