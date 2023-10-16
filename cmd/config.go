package cmd

import (
	"flag"
	"fmt"
)

const UrlBase = "https://go-challenge.skip.money"

var (
	HelpFlag        bool
	Threads         int
	TopTokens       int
	CollectionCount int
	CollectionURL   string
)

func ParseAndValidateFlags() error {
	ParseFlags()

	if HelpFlag {
		DisplayHelp()
		return fmt.Errorf("help flag provided")
	}

	if Threads <= 0 {
		return fmt.Errorf("number of threads should be greater than 0")
	}

	return nil
}

func ParseFlags() {
	flag.BoolVar(&HelpFlag, "h", false, "Display help information")
	flag.BoolVar(&HelpFlag, "help", false, "Display help information")
	flag.IntVar(&Threads, "t", 10, "Number of threads to use for fetching (default is 10)")
	flag.IntVar(&Threads, "threads", 10, "Number of threads to use for fetching (default is 10)")
	flag.IntVar(&TopTokens, "top", 5, "Number of top rare tokens to display (default is 5)")
	flag.IntVar(&CollectionCount, "count", 10000, "Number of tokens in the collection (default is 10000)")
	flag.StringVar(&CollectionURL, "url", "azuki1", "URL of the collection (default is 'azuki1')")
	flag.Parse()
}

func DisplayHelp() {
	fmt.Println(`Usage: mytool [OPTIONS]

This tool fetches token data and calculates the rarity of each token. The top N rarest tokens (specified by the -top option) are then displayed in both the console and CSV format.

Options:
	-h, --help      Display help information
	-threads   Maximum number of threads to use for fetching (default is 1)
	-top            Number of top rarity tokens to ouput (default is 5)
	-count          Number of tokens in the collection to fetch (default is 10000)
	-url            URL of the collection (default is 'azuki1')

Note: The tool currently fetches data from the following URL base:`, UrlBase)
}
