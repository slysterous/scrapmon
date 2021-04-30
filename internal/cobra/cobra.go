package cobra

import (
	"fmt"
	"strconv"

	scrapmon "github.com/slysterous/scrapmon/internal/scrapmon"
	"github.com/spf13/cobra"
)

// Client is responsible for interacting with cobra.
type Client struct {
	rootCmd *cobra.Command
}

//NewClient constructs a new Client.
func NewClient() *Client {
	var rootCmd = &cobra.Command{
		Use:   "scrapmon",
		Short: "Prntscr Scrapper",
		Long:  "A highly concurrent PrntScr Scrapper.",
	}
	return &Client{
		rootCmd: rootCmd,
	}
}

//NewPurgeCommand creates a new purge command.
func (c Client) NewPurgeCommand(purgeFn scrapmon.PurgeLogic) *cobra.Command {
	purgeCommand := &cobra.Command{
		Use:   "purge",
		Short: "Purges db and filesystem storage",
		RunE: func(cmd *cobra.Command, args []string) error {
			return purgeFn()
		},
		SilenceErrors: true,
	}
	return purgeCommand
}

//Execute executes the root command.
func (c Client) Execute() error {
	return c.rootCmd.Execute()
}

//NewStartCommand creates a new start command.
func (c Client) NewStartCommand(startFn scrapmon.StartLogic) (*cobra.Command, error) {
	startCommand := &cobra.Command{
		Use:   "start",
		Short: "Starts scraping images from imgur",
		RunE: func(cmd *cobra.Command, args []string) error {
			from, err := handleFromParam(cmd)
			if err != nil {
				return fmt.Errorf("command validation error, err: %v", err)
			}
			iterations, err := handleIterationsParam(cmd)
			if err != nil {
				return fmt.Errorf("command validation error, err: %v", err)
			}

			workersNumber, err := handleWorkersNumberParam(cmd)
			if err != nil {
				return fmt.Errorf("command validation error, err: %v", err)
			}
			return startFn(from, iterations, workersNumber)
		},
		SilenceErrors: true,
	}
	startCommand.Flags().StringP("from", "f", "", "starting imgur image code")
	startCommand.Flags().StringP("iterations", "i", "", "how many images should be downloaded")
	startCommand.Flags().StringP("workers", "w", "", "the amount of workers to be utilized for async operations")
	err := startCommand.MarkFlagRequired("workers")
	if err != nil {
		return nil, err
	}
	return startCommand, nil
}

//RegisterCommand registers a command onto the cobra client.
func (c Client) RegisterCommand(cmd *cobra.Command) {
	c.rootCmd.AddCommand(cmd)
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

func handleWorkersNumberParam(cmd *cobra.Command) (int, error) {
	workersString, err := cmd.Flags().GetString("workers")
	if err != nil {
		return 0, fmt.Errorf("could not parse --workers command, err: %v", err)
	}

	workersInt, err := strconv.Atoi(workersString)
	if err != nil && workersString != "" {
		return 0, fmt.Errorf("workers provided was not a number, err: %v", err)
	}

	if workersInt <= 0 {
		return 0, fmt.Errorf("workers have to be at least 1, err")
	}

	return workersInt, nil
}
