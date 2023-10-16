package calculator

import (
	"math"
	"testing"

	"github.com/a-straus/skip-protocol-interview/types"
)

const epsilon = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= epsilon
}

func TestCalculateTokenRarity(t *testing.T) {
	tests := []struct {
		name   string
		token  *types.Token
		traits types.CollectionTraits
		want   float64
	}{
		{
			name: "test single trait with high rarity",
			token: &types.Token{
				Attrs: map[string]string{"color": "blue"},
			},
			traits: types.CollectionTraits{
				"color": {"blue": 1, "red": 10, "green": 50},
			},
			want: 1.0 / (1 * 3), // 1 occurrence, 3 possible colors
		},
		{
			name: "multiple traits",
			token: &types.Token{
				Attrs: map[string]string{"color": "red", "shape": "circle"},
			},
			traits: types.CollectionTraits{
				"color": {"blue": 1, "red": 10, "green": 50},
				"shape": {"circle": 10, "square": 5},
			},
			want: 1.0/(10*3) + 1.0/(10*2),
		},
		{
			name: "trait not in collection traits",
			token: &types.Token{
				Attrs: map[string]string{"color": "yellow"},
			},
			traits: types.CollectionTraits{
				"color": {"blue": 1, "red": 10, "green": 50},
			},
			want: 0, // Not in the collection traits
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateTokenRarity(tt.token, tt.traits)
			if !almostEqual(got, tt.want) {
				t.Errorf("CalculateTokenRarity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateTokenRarities(t *testing.T) {
	tests := []struct {
		name   string
		tokens []*types.Token
		traits types.CollectionTraits
		want   []types.RarityScorecard
	}{
		{
			name: "test multiple tokens",
			tokens: []*types.Token{
				{Id: 1, Attrs: map[string]string{"color": "blue"}},
				{Id: 2, Attrs: map[string]string{"color": "red"}},
			},
			traits: types.CollectionTraits{
				"color": {"blue": 1, "red": 10, "green": 50},
			},
			want: []types.RarityScorecard{
				{Id: 1, Rarity: 1.0 / (1 * 3)},
				{Id: 2, Rarity: 1.0 / (10 * 3)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateTokenRarities(tt.tokens, tt.traits)
			if len(got) != len(tt.want) {
				t.Fatalf("Length mismatch: CalculateTokenRarities() = %v, want %v", len(got), len(tt.want))
			}
			for i, scorecard := range got {
				if scorecard.Id != tt.want[i].Id || scorecard.Rarity != tt.want[i].Rarity {
					t.Errorf("CalculateTokenRarities() = %v, want %v", scorecard, tt.want[i])
				}
			}
		})
	}
}
