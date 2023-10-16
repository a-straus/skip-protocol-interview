package calculator

import (
	"github.com/a-straus/skip-protocol-interview/types"
)

func CalculateTokenRarity(token *types.Token, traits types.CollectionTraits) float64 {
	rarity := 0.0
	for category, value := range token.Attrs {
		countWithTraitValue := traits[category][value]
		numValuesInCategory := len(traits[category])
		rarity += 1 / float64(countWithTraitValue*numValuesInCategory)
	}
	return rarity
}

func CalculateTokenRarities(tokens []*types.Token, traits types.CollectionTraits) []types.RarityScorecard {
	var scorecards []types.RarityScorecard

	for _, token := range tokens {
		rarity := CalculateTokenRarity(token, traits)
		scorecards = append(scorecards, types.RarityScorecard{
			Rarity: rarity,
			Id:     token.Id,
		})
	}

	return scorecards
}
