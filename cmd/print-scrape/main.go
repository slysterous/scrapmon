package printscrape


import (
	"github.com/spf13/cobra"
	"github.com/slysterous.com/print-scrape/internal/postgres"
)

var rootCmd = &cobra.Command {
	Use: "print-scrape",
	Short: "Prntscr Scrapper",
	Long:"A highly concurrent PrntScr Scrapper.",
}

var purgeCmd = &cobra.Command {
	Use: "",
	Short: "",
	Long: "",
	Run: purgeFn,
}

var findCmd = &cobra.Command {
	Use: "",
	Short: "",
	Long: "",
	Run: findFn,
}

var fetchCmd = &cobra.Command {
	Use: "",
	Short: "",
	Long: "",
	Run: fetchFn,
}

var scrapeCmd = &cobra.Command {
	Use: "",
	Short: "",
	Long: "",
	Run: scrapeFn,
}


func init() {
	rootCmd.AddCommand(purgeCmd)
	rootCmd.AddCommand(findCmd)
	rootCmd.AddCommand(fetchCmd)
	rootCmd.AddCommand(scrapeCmd)
}


func main(){
	err := rootCmd.Execute()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func purgeFn(_ *cobra.Command, args []string){
	db := postgres.Client{}
	postgresClient, err := db.NewClient("")
}

func findFn(_ *cobra.Command, args []string){

}

func fetchFn(_ *cobra.Command, args []string){

}

func scrapeFn(_ *cobra.Command, args []string){

}