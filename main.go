package main

import (
	"sort"
	"sync"

	"github.com/a-straus/skip-protocol-interview/calculator"
	"github.com/a-straus/skip-protocol-interview/cmd"
	"github.com/a-straus/skip-protocol-interview/fetcher"
	"github.com/a-straus/skip-protocol-interview/logger"
	"github.com/a-straus/skip-protocol-interview/output"
	"github.com/a-straus/skip-protocol-interview/types"
)

const (
	COLOR_GREEN = "\033[32m"
	COLOR_RED   = "\033[31m"
	COLOR_RESET = "\033[0m"
	maxRetries  = 3
)

var (
	semaphore         chan struct{}
	mutex             sync.Mutex
)

func main() {
	if err := cmd.ParseAndValidateFlags(); err != nil {
		logger.Error("Error: %v", err)
		return
	}

	azuki := types.Collection{
		Count: cmd.CollectionCount,
		Url:   cmd.CollectionURL,
	}

	tokens, traits := fetcher.GetTokens(azuki)

	rarities := calculator.CalculateRarities(tokens, traits)
	sort.Slice(rarities, func(i, j int) bool {
		return rarities[i].Rarity > rarities[j].Rarity
	})

	if err := output.WriteCSV(rarities, cmd.TopTokens); err != nil {
		logger.Error("Error writing CSV output: %v", err)
	}
}
