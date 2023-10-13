package output

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/a-straus/skip-protocol-interview/types"
)

func WriteCSV(rarities []types.RarityScorecard, top int) error {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	if err := w.Write([]string{"Rank", "Token ID", "Rarity"}); err != nil {
		return fmt.Errorf("error writing header to csv: %w", err)
	}

	for i := 0; i < top && i < len(rarities); i++ {
		rank := fmt.Sprintf("%d", i+1)
		tokenID := fmt.Sprintf("%d", rarities[i].Id)
		rarity := fmt.Sprintf("%.6f", rarities[i].Rarity)
		if err := w.Write([]string{rank, tokenID, rarity}); err != nil {
			return fmt.Errorf("error writing record to csv: %w", err)
		}
	}
	return nil
}
