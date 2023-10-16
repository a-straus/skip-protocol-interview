package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/a-straus/skip-protocol-interview/cmd"
	"github.com/a-straus/skip-protocol-interview/logger"
	"github.com/a-straus/skip-protocol-interview/types"
)

const (
	maxRetries = 3
)

func updateTraits(traits types.CollectionTraits, attrs map[string]string, mutex *sync.Mutex) {
	mutex.Lock()
	defer mutex.Unlock()

	for trait, value := range attrs {
		if _, exists := traits[trait]; !exists {
			traits[trait] = make(map[string]int)
		}
		traits[trait][value]++
	}
}

func getToken(tid int, colUrl string, traits types.CollectionTraits, wg *sync.WaitGroup, semaphore chan struct{}, mutex *sync.Mutex) (map[string]string, error) {
	defer wg.Done()
	defer func() { <-semaphore }()

	var lastError error
	for retry := 0; retry < maxRetries; retry++ {
		url := fmt.Sprintf("%s/%s/%d.json", cmd.UrlBase, colUrl, tid)
		res, err := http.Get(url)
		if err != nil {
			lastError = fmt.Errorf("error fetching token %d (attempt %d): %v", tid, retry+1, err)
			continue
		}
		defer res.Body.Close()
		if err == nil {
			body, err := io.ReadAll(res.Body)
			if err != nil {
				lastError = fmt.Errorf("error reading token %d data (attempt %d): %v", tid, retry+1, err)
				continue
			}

			attrs := make(map[string]string)
			if err := json.Unmarshal(body, &attrs); err != nil {
				lastError = fmt.Errorf("error unmarshalling token %d data (attempt %d): %v", tid, retry+1, err)
				continue
			}

			updateTraits(traits, attrs, mutex)

			return attrs, nil
		}
	}

	return nil, fmt.Errorf("Failed to fetch token %d after %d retries: %v", tid, maxRetries, lastError)
}

func GetTokens(col types.Collection) ([]*types.Token, types.CollectionTraits) {
	var wg sync.WaitGroup
	var semaphore = make(chan struct{}, cmd.Threads)
	var mutex sync.Mutex

	tokens := make([]*types.Token, col.Count)
	traits := make(types.CollectionTraits)

	for i := 0; i < col.Count; i++ {
		semaphore <- struct{}{}
		wg.Add(1)

		go func(i int) {
			logger.Info("%sGetting token %d%s", "\033[32m", i, "\033[0m")
			attrs, err := getToken(i, col.Url, traits, &wg, semaphore, &mutex)
			if err != nil {
				logger.Error("Error getting token %d:", i, err)
			}
			tokens[i] = &types.Token{Id: i, Attrs: attrs}
		}(i)
	}
	wg.Wait()

	return tokens, traits
}
