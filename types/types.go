package types

type Token struct {
	Id    int
	Attrs map[string]string
}

type Collection struct {
	Count int
	Url   string
}

type CollectionTraits map[string]map[string]int

type RarityScorecard struct {
	Rarity float64
	Id     int
}