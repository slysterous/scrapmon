package main

import (
	"fmt"
	"github.com/joho/godotenv"
	cobraClient "github.com/slysterous/scrapmon/internal/cobra"
	config "github.com/slysterous/scrapmon/internal/config"
	file "github.com/slysterous/scrapmon/internal/file"
	phttp "github.com/slysterous/scrapmon/internal/http"
	"github.com/slysterous/scrapmon/internal/logger"
	"github.com/slysterous/scrapmon/internal/postgres"
	scrapmon "github.com/slysterous/scrapmon/internal/scrapmon"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	pgClient, err := postgres.NewClient(getDataSource(conf), conf.MaxDBConnections)
	if err != nil {
		log.Fatalf("could not connect to DB, err: %v", err)
	}
	defer pgClient.DB.Close()

	// init a file manager.
	fileManager := file.NewManager(conf.ScrapStorageFolder, writer{}, purger{})

	//combine db and filestorage into generic storage.
	storage := scrapmon.Storage{
		Fm: fileManager,
		Dm: pgClient,
	}

	//TODO fix the TOR client
	//scrapper := phttp.NewProxyChainClient("127.0.0.1", "9050")
	scrapper := phttp.NewClient("https://i.imgur.com/", reader{}, &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	})

	commandManager := scrapmon.ConcurrentCommandManager{
		Storage:       storage,
		Scrapper:      scrapper,
		CodeAuthority: scrapmon.ConcurrentCodeAuthority{
			Logger: logger.NewLogger(1,os.Stdout),
		},
	}

	cobraC := cobraClient.NewClient()

	startCommand, err := cobraC.NewStartCommand(commandManager.StartCommand)
	if err != nil {
		log.Fatalf("could not register start command, err: %v", err)
	}

	purgeCommand := cobraC.NewPurgeCommand(commandManager.PurgeCommand)

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

type writer struct{}

func (w writer) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(filename, data, perm)
}

type purger struct{}

func (p purger) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

func (p purger) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

type reader struct{}

func (r reader) ReadAll(re io.Reader) ([]byte, error) {
	return ioutil.ReadAll(re)
}
