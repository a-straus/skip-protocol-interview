package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
)

const (
	urlBase       = "https://go-challenge.skip.money"
	collection    = "azuki"
	COLOR_GREEN = "\033[32m"
	COLOR_RED   = "\033[31m"
	COLOR_RESET = "\033[0m"
)

var logger *log.Logger = log.New(os.Stderr, "", log.Ldate|log.Ltime)

var helpFlag bool
var threads int

func init() {
	flag.BoolVar(&helpFlag, "h", false, "Display help information")
	flag.BoolVar(&helpFlag, "help", false, "Display help information")
	flag.IntVar(&threads, "t", 1, "Number of threads to use for fetching (default is 1)")
	flag.IntVar(&threads, "threads", 1, "Number of threads to use for fetching (default is 1)")
}

func displayHelp() {
	fmt.Println(`Usage: mytool [OPTIONS]

This tool fetches token data and calculates the rarity of each token. The top 5 rarest tokens are then displayed in both the console and CSV format.

Options:
  -h, --help      Display help information
	-t, --threads   Maximum number of threads

Note: The tool currently fetches data from the following URL base:`, urlBase)
}

type Token struct {
	Id    int
	Attrs map[string]string
}

type RarityScorecard struct {
	Rarity float64
	Id     int
}

type Collection struct {
	Count int
	Url   string
}

type CollectionTraits map[string]map[string]int

func getToken(tid int, colUrl string) (map[string]string, error) {
	url := fmt.Sprintf("%s/%s/%d.json", urlBase, colUrl, tid)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	attrs := make(map[string]string)
	if err := json.Unmarshal(body, &attrs); err != nil {
		return nil, err
	}

	return attrs, nil
}

func getTokens(col Collection) ([]*Token, CollectionTraits, error) {
	tokens := make([]*Token, col.Count)
	traits := make(CollectionTraits)

	for i := 0; i < col.Count; i++ {
		logger.Println(string(COLOR_GREEN), fmt.Sprintf("Getting token %d", i), string(COLOR_RESET))
		attrs, err := getToken(i, col.Url)
		if err != nil {
			logger.Println(fmt.Sprintf("Error getting token %d:", i), err)
			continue
		}

		tokens[i] = &Token{Id: i, Attrs: attrs}
		for trait, value := range tokens[i].Attrs {
			if _, exists := traits[trait]; !exists {
				traits[trait] = make(map[string]int)
			}
			traits[trait][value]++
		}
	}

	return tokens, traits, nil
}

func calculateTokenRarity(token *Token, traits CollectionTraits) float64 {
	rarity := 0.0
	for category, value := range token.Attrs {
		countWithTraitValue := traits[category][value]
		numValuesInCategory := len(traits[category])
		rarity += 1 / float64(countWithTraitValue*numValuesInCategory)
	}
	return rarity
}

func calculateRarities(tokens []*Token, traits CollectionTraits) []RarityScorecard {
	rarities := make([]RarityScorecard, len(tokens))
	for i, token := range tokens {
		rarities[i] = RarityScorecard{
			Rarity: calculateTokenRarity(token, traits),
			Id:     token.Id,
		}
	}
	return rarities
}

func writeCSV(rarities []RarityScorecard) error {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	// Write the header
	if err := w.Write([]string{"Rank", "Token ID", "Rarity"}); err != nil {
			return fmt.Errorf("error writing header to csv: %w", err)
	}

	// Write the 5 rarest tokens to stdout in CSV format
	for i := 0; i < 5 && i < len(rarities); i++ {
			rank := fmt.Sprintf("%d", i+1)
			tokenID := fmt.Sprintf("%d", rarities[i].Id)
			rarity := fmt.Sprintf("%.6f", rarities[i].Rarity)
			if err := w.Write([]string{rank, tokenID, rarity}); err != nil {
					return fmt.Errorf("error writing record to csv: %w", err)
			}
	}
	return nil
}

func main() {
	flag.Parse()

	if helpFlag {
		displayHelp()
		return
	}

	if threads <= 0 {
		logger.Println(string(COLOR_RED), "Error: Number of threads should be greater than 0", string(COLOR_RESET))
		return
	}

	azuki := Collection{
		Count: 10,
		Url:   "azuki1",
	}

	tokens, traits, err := getTokens(azuki)
	if err != nil {
		logger.Println(string(COLOR_RED), "Error:", err, string(COLOR_RESET))
		return
	}

	rarities := calculateRarities(tokens, traits)
	sort.Slice(rarities, func(i, j int) bool {
		return rarities[i].Rarity > rarities[j].Rarity
	})

	logger.Println("------ 5 Rarest Tokens ------")
	for i := 0; i < 5 && i < len(rarities); i++ {
		logger.Printf("Rank %d:\n", i+1)
		logger.Printf("Token ID: %d\n", rarities[i].Id)
		logger.Printf("Rarity: %.6f\n", rarities[i].Rarity)
		logger.Println("----------------------------")
	}

	if err := writeCSV(rarities); err != nil {
		logger.Println("Error writing CSV:", err)
	}
}
