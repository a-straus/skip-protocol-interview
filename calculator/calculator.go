package calculator

import (
	"container/heap"

	"github.com/a-straus/skip-protocol-interview/cmd"
	"github.com/a-straus/skip-protocol-interview/types"
)

type MinHeap []types.RarityScorecard

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].Rarity < h[j].Rarity }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(types.RarityScorecard))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}


func CalculateTokenRarity(token *types.Token, traits types.CollectionTraits) float64 {
	rarity := 0.0
	for category, value := range token.Attrs {
		countWithTraitValue := traits[category][value]
		numValuesInCategory := len(traits[category])
		rarity += 1 / float64(countWithTraitValue*numValuesInCategory)
	}
	return rarity
}

func CalculateRarities(tokens []*types.Token, traits types.CollectionTraits) []types.RarityScorecard {
	h := &MinHeap{}
	heap.Init(h)

	for _, token := range tokens {
		rarity := CalculateTokenRarity(token, traits)
		scorecard := types.RarityScorecard{
			Rarity: rarity,
			Id:     token.Id,
		}
		heap.Push(h, scorecard)
		if h.Len() > cmd.TopTokens {
			heap.Pop(h) // This pops the smallest rarity, ensuring only top rarities are retained
		}
	}

	var scorecards []types.RarityScorecard
	for h.Len() > 0 {
		scorecards = append(scorecards, heap.Pop(h).(types.RarityScorecard))
	}
	for i, j := 0, len(scorecards)-1; i < j; i, j = i+1, j-1 {
		scorecards[i], scorecards[j] = scorecards[j], scorecards[i]
	}

	return scorecards

}
