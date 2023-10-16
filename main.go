package main

import (
	"sort"

	"github.com/a-straus/skip-protocol-interview/calculator"
	"github.com/a-straus/skip-protocol-interview/cmd"
	"github.com/a-straus/skip-protocol-interview/fetcher"
	"github.com/a-straus/skip-protocol-interview/logger"
	"github.com/a-straus/skip-protocol-interview/output"
	"github.com/a-straus/skip-protocol-interview/types"
)

func main() {
	if err := cmd.ParseAndValidateFlags(); err != nil {
		logger.Error("Error: %v", err)
		return
	}

	collection := types.Collection{
		Count: cmd.CollectionCount,
		Url:   cmd.CollectionURL,
	}

	tokens, traits := fetcher.GetTokens(collection)
	rarities := calculator.CalculateTokenRarities(tokens, traits)

	sort.Slice(rarities, func(i, j int) bool {
		return rarities[i].Rarity > rarities[j].Rarity
	})

	if err := output.WriteCSV(rarities, cmd.TopTokens); err != nil {
		logger.Error("Error writing CSV output: %v", err)
	}
}
