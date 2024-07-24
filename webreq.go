package main

import (
	"encoding/json"
	"fmt"
	"io"
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

func unmarshallPokeLocationResult(jsonByteSlc []byte) (pokeLocationResult, error) {
	plResult := pokeLocationResult{}
	err := json.Unmarshal(jsonByteSlc, &plResult)
	if err != nil {
		return pokeLocationResult{}, fmt.Errorf("failed to convert json to location data: %s", err)
	}
	return plResult, nil
}
