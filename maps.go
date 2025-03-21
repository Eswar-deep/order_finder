package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func geocodeAddress(address string) (Coordinates, error) {
	apiKey := googleAPIKey
	if apiKey == "" {
		return Coordinates{}, fmt.Errorf("google Maps API key is not set")
	}

	addressEncoded := url.QueryEscape(address)
	apiURL := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s", addressEncoded, apiKey)

	resp, err := http.Get(apiURL)
	if err != nil {
		return Coordinates{}, fmt.Errorf("error making request to Google Maps API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Coordinates{}, fmt.Errorf("error response from Google Maps API: %s", resp.Status)
	}

	var result struct {
		Results []struct {
			Geometry struct {
				Location Coordinates `json:"location"`
			} `json:"geometry"`
		} `json:"results"`
		Status string `json:"status"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return Coordinates{}, fmt.Errorf("error decoding Google Maps API response: %v", err)
	}

	if result.Status != "OK" || len(result.Results) == 0 {
		return Coordinates{}, fmt.Errorf("no results found for address: %s", address)
	}

	return result.Results[0].Geometry.Location, nil
}
