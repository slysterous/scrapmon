package cobra

import (
	"fmt"
	printscrape "github.com/slysterous/print-scrape/internal/domain"
	"github.com/spf13/cobra"
	"strconv"
)


// Client is responsible for interacting with cobra.
type Client struct {
	rootCmd        *cobra.Command
	store printscrape.Storage
	scrapper printscrape.ImageScrapper
}

// NewClient constructs a new Client.
func NewClient() *Client {
	var rootCmd = &cobra.Command{
		Use:   "print-scrape",
		Short: "Prntscr Scrapper",
		Long:  "A highly concurrent PrntScr Scrapper.",
	}
	return &Client{
		rootCmd: rootCmd,
	}

}

// RegisterStartCommand registers the start command to cobra
func (c Client) RegisterStartCommand() {
	startCmd := c.createStartCmd(c.store,c.scrapper)
	c.rootCmd.AddCommand(startCmd)
}

//RegisterPurgeCommand registers the purge command to cobra
func (c Client) RegisterPurgeCommand() {
	purgeCmd := c.createPurgeCmd(c.store)
	c.rootCmd.AddCommand(purgeCmd)
}

// Execute executes the root command.
func (c Client) Execute() error {
	return c.rootCmd.Execute()
}

func (c Client) createStartCmd(store printscrape.Storage, scrapper printscrape.ImageScrapper) *cobra.Command {
	startCommand := &cobra.Command{
		Use:   "start",
		Short: "Starts scraping images from imgur",
		RunE: func(cmd *cobra.Command, args []string) error {
			from,err :=handleFromParam(cmd)
			if err !=nil {
				return fmt.Errorf("command validation error, err: %v",err)
			}
			iterations,err:=handleIterationsParam(cmd)
			if err !=nil {
				return fmt.Errorf("command validation error, err: %v",err)
			}
			return printscrape.StartCommand(store,scrapper,from,iterations) 
		},
		SilenceErrors: true,
	}
	startCommand.Flags().StringP("from", "f", "", "starting imgur image code")
	startCommand.Flags().StringP("iterations", "i", "", "how many images should be downloaded")

	return startCommand
}

func (c Client) createPurgeCmd(store printscrape.Storage) *cobra.Command {
	purgeCommand := &cobra.Command{
		Use: "purge",
		Short: "Purges db and filesystem storage",
		RunE: func(cmd *cobra.Command, args []string) error {
			return printscrape.PurgeCommand(store)
		},
		SilenceErrors:true,
	}

	return purgeCommand
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
