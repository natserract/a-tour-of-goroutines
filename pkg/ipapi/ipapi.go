package ipapi

import (
	"encoding/json"
	"fmt"
	"goroutines/pkg/api"
)

type Response struct {
	Asn                string      `json:"asn"`
	City               string      `json:"city"`
	ContinentCode      string      `json:"continent_code"`
	Country            string      `json:"country"`
	CountryArea        float64     `json:"country_area"`
	CountryCallingCode string      `json:"country_calling_code"`
	CountryCapital     string      `json:"country_capital"`
	CountryCode        string      `json:"country_code"`
	CountryCodeIso3    string      `json:"country_code_iso3"`
	CountryName        string      `json:"country_name"`
	CountryPopulation  int64       `json:"country_population"`
	CountryTld         string      `json:"country_tld"`
	Currency           string      `json:"currency"`
	CurrencyName       string      `json:"currency_name"`
	InEu               bool        `json:"in_eu"`
	IP                 string      `json:"ip"`
	Languages          string      `json:"languages"`
	Latitude           float64     `json:"latitude"`
	Longitude          float64     `json:"longitude"`
	Network            string      `json:"network"`
	Org                string      `json:"org"`
	Postal             interface{} `json:"postal"`
	Region             string      `json:"region"`
	RegionCode         string      `json:"region_code"`
	Timezone           string      `json:"timezone"`
	UtcOffset          string      `json:"utc_offset"`
	Version            string      `json:"version"`
}

func Request() (*Response, error) {
	client, err := api.NewClient("https://ipapi.co")
	if err != nil {
		fmt.Printf("Client error: %s\n", err)
		return nil, err
	}

	req, err := client.NewRequest("GET", "/json", nil)
	if err != nil {
		fmt.Printf("Request error: %s\n", err)
		return nil, err
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()

	var body *Response
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
