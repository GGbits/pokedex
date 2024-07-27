package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func getApiResponse(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("Response failed with status code: %d and\nbody:%s\n", res.StatusCode, body)
	}
	return body, nil
}

func queryCache(cfg *config, url string) ([]byte, error) {
	webResult, ok := cfg.cache.Get(url)
	if !ok {
		log.Println("Missed cache, trying page directly..")
		wr, err := getApiResponse(url)
		if err != nil {
			return nil, fmt.Errorf("an error occured while getting API response: %s", err)
		}
		log.Println("adding data to cache...")
		cfg.cache.Add(url, wr)
		webResult = wr
	}
	return webResult, nil
}

func unmarshallPokeLocationResult(jsonByteSlc []byte) (pokeLocationResult, error) {
	plResult := pokeLocationResult{}
	err := json.Unmarshal(jsonByteSlc, &plResult)
	if err != nil {
		return pokeLocationResult{}, fmt.Errorf("failed to convert json to location data: %s", err)
	}
	return plResult, nil
}

func unmarshallPokeExploreResult(jsonByteSlc []byte) (pokeExploreResult, error) {
	pExResult := pokeExploreResult{}
	err := json.Unmarshal(jsonByteSlc, &pExResult)
	if err != nil {
		return pokeExploreResult{}, fmt.Errorf("failed to convert json to explore data: %s", err)
	}
	return pExResult, nil
}
