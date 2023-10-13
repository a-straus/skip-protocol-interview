package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/a-straus/skip-protocol-interview/logger"
	"github.com/a-straus/skip-protocol-interview/types"
)

const (
	urlBase    = "https://go-challenge.skip.money"
	maxRetries = 3
)


var mutex sync.Mutex

func getToken(tid int, colUrl string, traits types.CollectionTraits, wg *sync.WaitGroup, semaphore chan struct{}) (map[string]string, error) {
	defer wg.Done()
	defer func() { <-semaphore }()

	for retry := 0; retry < maxRetries; retry++ {
		url := fmt.Sprintf("%s/%s/%d.json", urlBase, colUrl, tid)
		res, err := http.Get(url)
		if err == nil {
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				return nil, err
			}

			attrs := make(map[string]string)
			if err := json.Unmarshal(body, &attrs); err != nil {
				return nil, err
			}

			mutex.Lock()
			for trait, value := range attrs {
				if _, exists := traits[trait]; !exists {
					traits[trait] = make(map[string]int)
				}
				traits[trait][value]++
			}
			mutex.Unlock()

			return attrs, nil
		}
	}

	logger.Error("Error getting token %d after %d retries", tid, maxRetries)
	return nil, fmt.Errorf("Failed to fetch token %d after %d retries", tid, maxRetries)
}

func GetTokens(col types.Collection, semaphore chan struct{}) ([]*types.Token, types.CollectionTraits) {
	var wg sync.WaitGroup
	tokens := make([]*types.Token, col.Count)
	traits := make(types.CollectionTraits)

	for i := 0; i < col.Count; i++ {
		semaphore <- struct{}{}
		wg.Add(1)

		go func(i int) {
			logger.Info("%sGetting token %d%s", "\033[32m", i, "\033[0m")
			attrs, err := getToken(i, col.Url, traits, &wg, semaphore)
			if err != nil {
				logger.Error("Error getting token %d:", i, err)
			}
			tokens[i] = &types.Token{Id: i, Attrs: attrs}
		}(i)
	}
	wg.Wait()

	return tokens, traits
}
