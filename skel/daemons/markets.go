package daemons

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/openware/pkg/mngapi/peatio"
)

// Define response data
type MarketResponse struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	BaseUnit        string `json:"base_unit"`
	QuoteUnit       string `json:"quote_unit"`
	State           string `json:"state"`
	AmountPrecision int64  `json:"amount_precision"`
	PricePrecision  int64  `json:"price_precision"`
	MinPrice        string `json:"min_price"`
	MaxPrice        string `json:"max_price"`
	MinAmount       string `json:"min_amount"`
	Position        int64  `json:"position"`
}

func FetchMarkets(peatioClient *peatio.Client, opendaxAddr string) {
	for {
		FetchMarketsFromOpenfinexCloud(peatioClient, opendaxAddr)
		<-time.After(5 * time.Minute)
	}
}

func FetchMarketsFromOpenfinexCloud(peatioClient *peatio.Client, opendaxAddr string) error {
	url := fmt.Sprintf("%s/api/v2/opx/markets", opendaxAddr)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("ERROR: FetchMarketsFromOpenfinexCloud: Can't fetch markets: %v", err.Error())
		return err
	}
	// Call HTTP request
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR: FetchMarketsFromOpenfinexCloud: Request failed: %v", err.Error())
		return err
	}
	defer resp.Body.Close()

	// Convert response body to []byte
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: FetchMarketsFromOpenfinexCloud: Can't convert body to []: %d -> %v", resp.StatusCode, err.Error())
		return err
	}
	// Check for API error
	if resp.StatusCode != http.StatusOK {
		log.Printf("ERROR: FetchMarketsFromOpenfinexCloud: Unexpected status: %d", resp.StatusCode)
		return errors.New(fmt.Sprintf("ERROR: FetchMarketsFromOpenfinexCloud: Unexpected status: %d", resp.StatusCode))
	}

	// Unmarshal response body result
	markets := []MarketResponse{}
	marshalErr := json.Unmarshal(resBody, &markets)
	if marshalErr != nil {
		log.Printf("ERROR: FetchMarketsFromOpenfinexCloud: Can't unmarshal response. %v", marshalErr)
		return marshalErr
	}

	// Iterate through all markets
	for _, market := range markets {
		// Find market by ID, if there is no system will create
		res, apiError := peatioClient.GetMarketByID(market.ID)
		// Check result here
		if res == nil && apiError != nil {
			marketParams := peatio.CreateMarketParams{
				BaseCurrency:    market.BaseUnit,
				QuoteCurrency:   market.QuoteUnit,
				State:           "disabled",
				EngineName:      "opendax-cloud-engine",
				AmountPrecision: market.AmountPrecision,
				PricePrecision:  market.PricePrecision,
				MinPrice:        market.MinPrice,
				MaxPrice:        market.MaxPrice,
				MinAmount:       market.MinAmount,
				Position:        market.Position,
			}

			_, apiError := peatioClient.CreateMarket(marketParams)
			if apiError != nil {
				log.Printf("ERROR: FetchMarketsFromOpenfinexCloud: Can't create market with id %s. Error: %v. Errors: %v", market.ID, apiError.Error, apiError.Errors)
			}
		}
	}

	return nil
}
