package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func geocodeAddress(address string) (Coordinates, error) {
	addressEncoded := url.QueryEscape(address)
	url := fmt.Sprintf("https://api.tomtom.com/search/2/geocode/%s.json?key=%s", addressEncoded, tomtomAPIKey)

	resp, err := http.Get(url)
	if err != nil {
		return Coordinates{}, fmt.Errorf("error making request to TomTom API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Coordinates{}, fmt.Errorf("error response from TomTom API: %s", resp.Status)
	}

	var result struct {
		Results []struct {
			Position Coordinates `json:"position"`
		} `json:"results"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return Coordinates{}, fmt.Errorf("error decoding TomTom API response: %v", err)
	}

	if len(result.Results) == 0 {
		return Coordinates{}, fmt.Errorf("no results found for address")
	}

	return result.Results[0].Position, nil
}
