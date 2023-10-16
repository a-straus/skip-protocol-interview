package main

import (
	"sync"
	"time"

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

	startTime := time.Now()
	rarities := calculator.CalculateRarities(tokens, traits)

	elapsedTime := time.Since(startTime)

	logger.Info("Rarity calculation took %s", elapsedTime)

	if err := output.WriteCSV(rarities, cmd.TopTokens); err != nil {
		logger.Error("Error writing CSV output: %v", err)
	}
}
