package main

import (
	"fmt"
	cfg "github.com/slysterous/print-scrape/internal/config"
	printscrape "github.com/slysterous/print-scrape/internal/domain"
	cobraClient "github.com/slysterous/print-scrape/internal/cobra"
	file "github.com/slysterous/print-scrape/internal/file"
	phttp "github.com/slysterous/print-scrape/internal/http"
	"github.com/slysterous/print-scrape/internal/postgres"
	"log"
	"os"
	"strconv"
	"time"
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

	cobraClient := cobraClient.NewClient(storage,scrapper)
	cobraClient.RegisterStartCommand()
	cobraClient.RegisterPurgeCommand()

	if err:= cobraClient.Execute(); err !=nil {
		log.Fatalf("execution failed, err: %v",err)
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

func handleFromParam(cmd *cobra.Command) (string, error) {
	fromCode, err := cmd.Flags().GetString("from")
	if err != nil {
		return "", fmt.Errorf("could not parse --from command, err: %v", err)
	}

	if fromCode != "" {
		return fromCode, nil
	}

	fmt.Print("--from was not provided. This will start downloading from code 0 Are you ok with that? (Y/n)")
	res := askForConfirmation()
	if !res {
		return fromCode, fmt.Errorf("main: User requested to abort procedure")
	}
	return "", nil
}

func handleIterationsParam(cmd *cobra.Command) (int, error) {
	iterationsString, err := cmd.Flags().GetString("iterations")
	if err != nil {
		return 0, fmt.Errorf("could not parse --from command, err: %v", err)
	}

	iterationsInt, err := strconv.Atoi(iterationsString)
	if err != nil && iterationsString != "" {
		return 0, fmt.Errorf("count provided was not a number, err: %v", err)
	}

	if iterationsString == "" {
		iterationsInt = -1
		fmt.Print("--iterations was not provided. This will continue downloading until you cancel the operation. Is that ok?  (Y/n)")
		res := askForConfirmation()
		if !res {
			return iterationsInt, fmt.Errorf("main: User requested to abort procedure")
		}
	}

	return iterationsInt, nil
}
