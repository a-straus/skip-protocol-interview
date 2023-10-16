package cmd

import (
	"flag"
	"fmt"
)

const (
	UrlBase                = "https://go-challenge.skip.money"
	defaultThreads         = 10
	defaultTopTokens       = 5
	defaultCollectionCount = 10000
	defaultCollectionURL   = "azuki1"
)

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
	flag.BoolVar(&HelpFlag, "help", false, "Display help information")
	flag.IntVar(&Threads, "threads", defaultThreads, fmt.Sprintf("Number of threads to use for fetching (default is %d)", defaultThreads))
	flag.IntVar(&TopTokens, "top", defaultTopTokens, fmt.Sprintf("Number of top rare tokens to display (default is %d)", defaultTopTokens))
	flag.IntVar(&CollectionCount, "count", defaultCollectionCount, fmt.Sprintf("Number of tokens in the collection (default is %d)", defaultCollectionCount))
	flag.StringVar(&CollectionURL, "url", defaultCollectionURL, fmt.Sprintf("URL of the collection (default is '%s')", defaultCollectionURL))
	flag.Parse()
}

func DisplayHelp() {
	fmt.Println(`Usage: mytool [OPTIONS]

This tool fetches token data and calculates the rarity of each token. The top N rarest tokens (specified by the -top option) are then displayed in both the console and CSV format.

Options:
	--help           Display help information
	--threads        Maximum number of threads to use for fetching (default is`, defaultThreads, `)
	--top            Number of top rarity tokens to ouput (default is`, defaultTopTokens, `)
	--count          Number of tokens in the collection to fetch (default is`, defaultCollectionCount, `)
	--url            URL of the collection (default is '`, defaultCollectionURL, `')

Note: The tool currently fetches data from the following URL base:`, UrlBase)
}
