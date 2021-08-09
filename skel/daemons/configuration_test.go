package daemons

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/openware/pkg/mngapi/peatio"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

const peatioManagementURL = "api/v2/peatio/management"

func TestFetchConfigurationSuccess(t *testing.T) {
	mockedResponse := []byte(`{"currencies":[{"id":"btc","name":"Bitcoin","description":"","homepage":"","price":"39100.0","status":"enabled","type":"coin","precision":6,"position":4,"icon_url":"","networks":[]}],"markets":[{"id":"omgusdt","name":"OMG/USDT","base_unit":"omg","quote_unit":"usdt","state":"enabled","amount_precision":2,"price_precision":4,"min_price":"0.0037","max_price":"3687.4597","min_amount":"0.01","position":9}]}`)
	marketPeatioResponse := []byte(`{"id":"omgusdt","name":"OMG/USDT","base_unit":"omg","quote_unit":"usdt","state":"enabled","amount_precision":2,"price_precision":4,"min_price":"0.0037","max_price":"3687.4597","min_amount":"0.01","position":9}`)
	currencyPeatioResponse := []byte(`{"id":"btc","name":"Bitcoin","description":"","homepage":"","price":"39100.0","status":"enabled","type":"coin","precision":6,"position":4,"icon_url":""}`)
	networkPeatioResponse := []byte(`{"blockchain_key":"btc-testnet","currency_id":"btc","deposit_enabled":false,"withdrawal_enabled":true,"deposit_fee":"0.0","min_deposit_amount":"0.0","withdraw_fee":"0.0000000002557544","min_withdraw_amount":"0.0000000025575447","base_factor":1000000000000000000}`)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/opx/config", func(res http.ResponseWriter, req *http.Request) {
		res.Write(mockedResponse)
	})
	mux.HandleFunc("/api/v2/peatio/management/markets/new", func(res http.ResponseWriter, req *http.Request) {
		res.Write(marketPeatioResponse)
	})
	mux.HandleFunc("/api/v2/peatio/management/currencies/create", func(res http.ResponseWriter, req *http.Request) {
		res.Write(currencyPeatioResponse)
	})
	mux.HandleFunc("/api/v2/peatio/management/blockchain_currencies/new", func(res http.ResponseWriter, req *http.Request) {
		res.Write(networkPeatioResponse)
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	peatioClient, err := peatio.New(fmt.Sprintf("%s/%s", ts.URL, peatioManagementURL), jwtIssuer, jwtAlgo, jwtPrivateKey)
	require.NoError(t, err)

	app := initApp()
	app.Conf.Opendax.Addr = ts.URL
	res := FetchConfiguration(peatioClient, app.Conf.Opendax.Addr, "platformID")
	assert.Equal(t, res, nil)
}

func TestFetchConfigurationEmptyResponse(t *testing.T) {
	mockedResponse := []byte(`{}`)
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v2/opx/config", func(res http.ResponseWriter, req *http.Request) {
		res.Write(mockedResponse)
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	peatioClient, err := peatio.New(fmt.Sprintf("%s/%s", ts.URL, peatioManagementURL), jwtIssuer, jwtAlgo, jwtPrivateKey)
	require.NoError(t, err)

	app := initApp()
	app.Conf.Opendax.Addr = ts.URL
	res := FetchConfiguration(peatioClient, app.Conf.Opendax.Addr, "platformID")
	assert.Equal(t, res, nil)
}

func TestFetchConfigurationHostError(t *testing.T) {
	response := "Unexpected status: 404"
	mux := http.NewServeMux()

	ts := httptest.NewServer(mux)
	defer ts.Close()

	peatioClient, err := peatio.New(fmt.Sprintf("%s/%s", ts.URL, peatioManagementURL), jwtIssuer, jwtAlgo, jwtPrivateKey)
	require.NoError(t, err)

	app := initApp()
	app.Conf.Opendax.Addr = ts.URL
	res := FetchConfiguration(peatioClient, app.Conf.Opendax.Addr, "platfromID")
	require.Error(t, res)
	assert.Equal(t, res.Error(), response)
}

func TestDivideCurrenciesIntoGroups(t *testing.T) {
	response := []CurrencyResponse{
		{
			ID: "link", Name: "LINK", Description: "", Homepage: "",
			Price: "0.0", Type: "coin", Precision: 16, Position: 1, IconURL: "",
			Networks: []BlockchainCurrencyResponse{{CurrencyID: "link", BlockchainKey: "eth-rinkeby", ParentID: "eth",
				DepositEnabled: true, WithdrawEnabled: true, DepositFee: "0.0", MinDepositAmount: "0.0",
				WithdrawFee: "0.0", MinWithdrawAmount: "0.0",
				BaseFactor: 8}, {
				CurrencyID: "link", BlockchainKey: "eth-rinkeby", ParentID: "",
				DepositEnabled: true, WithdrawEnabled: true, DepositFee: "0.0", MinDepositAmount: "0.0",
				WithdrawFee: "0.0", MinWithdrawAmount: "0.0",
				BaseFactor: 8},
			},
		},
		{
			ID: "eth", Name: "ETH", Description: "", Homepage: "",
			Price: "0.0", Type: "coin", Precision: 16, Position: 1, IconURL: "",
			Networks: []BlockchainCurrencyResponse{{CurrencyID: "eth", BlockchainKey: "eth-rinkeby", ParentID: "",
				DepositEnabled: true, WithdrawEnabled: true, DepositFee: "0.0", MinDepositAmount: "0.0",
				WithdrawFee: "0.0", MinWithdrawAmount: "0.0",
				BaseFactor: 8}},
		},
		{
			ID: "usdt", Name: "USDT", Description: "", Homepage: "",
			Price: "0.0", Type: "coin", Precision: 16, Position: 1, IconURL: "",
			Networks: []BlockchainCurrencyResponse{{CurrencyID: "usdt", BlockchainKey: "eth-rinkeby", ParentID: "eth",
				DepositEnabled: true, WithdrawEnabled: true, DepositFee: "0.0", MinDepositAmount: "0.0",
				WithdrawFee: "0.0", MinWithdrawAmount: "0.0",
				BaseFactor: 8}},
		},
		{
			ID: "eur", Name: "EUR", Description: "", Homepage: "",
			Price: "0.0", Type: "fiat", Precision: 16, Position: 1, IconURL: "",
			Networks: []BlockchainCurrencyResponse{{CurrencyID: "eur", BlockchainKey: "", ParentID: "",
				DepositEnabled: true, WithdrawEnabled: true, DepositFee: "0.0", MinDepositAmount: "0.0",
				WithdrawFee: "0.0", MinWithdrawAmount: "0.0",
				BaseFactor: 8}},
		},
		{
			ID: "tron", Name: "TRON", Description: "", Homepage: "",
			Price: "0.0", Type: "coin", Precision: 16, Position: 1, IconURL: "",
			Networks: []BlockchainCurrencyResponse{{CurrencyID: "tron", BlockchainKey: "tron-testnet", ParentID: "",
				DepositEnabled: true, WithdrawEnabled: true, DepositFee: "0.0", MinDepositAmount: "0.0",
				WithdrawFee: "0.0", MinWithdrawAmount: "0.0",
				BaseFactor: 8}},
		},
		{
			ID: "xrp", Name: "XRP", Description: "", Homepage: "",
			Price: "0.0", Type: "coin", Precision: 16, Position: 1, IconURL: "",
			Networks: []BlockchainCurrencyResponse{{CurrencyID: "xrp", BlockchainKey: "xrp-testnet", ParentID: "",
				DepositEnabled: true, WithdrawEnabled: true, DepositFee: "0.0", MinDepositAmount: "0.0",
				WithdrawFee: "0.0", MinWithdrawAmount: "0.0",
				BaseFactor: 8}},
		},
		{
			ID: "txrp", Name: "TXRP", Description: "", Homepage: "",
			Price: "0.0", Type: "coin", Precision: 16, Position: 1, IconURL: "",
			Networks: []BlockchainCurrencyResponse{{CurrencyID: "txrp", BlockchainKey: "xrp-testnet", ParentID: "xrp",
				DepositEnabled: true, WithdrawEnabled: true, DepositFee: "0.0", MinDepositAmount: "0.0",
				WithdrawFee: "0.0", MinWithdrawAmount: "0.0",
				BaseFactor: 8}},
		},
	}

	expectedResult := make(map[string][]string)
	expectedResult["eth"] = []string{"eth", "link", "usdt"}
	expectedResult["tron"] = []string{"tron"}
	expectedResult["xrp"] = []string{"xrp", "txrp"}

	actualResult := divideCurrenciesIntoGroups(response)
	assert.Equal(t, reflect.DeepEqual(expectedResult, actualResult), true)
}

func TestFindCurrenciesInWallets(t *testing.T) {
	const blockchainKey = "opendax-cloud"
	wallets := []*peatio.Wallet{
		{
			ID: 1, Name: "BTC Deposit Wallet", Kind: "deposit",
			Currencies: []string{"btc"}, Address: "address", Gateway: "opendax_cloud",
			MaxBalance: "0.0", Balance: "0.0", BlockchainKey: blockchainKey, Status: "active",
		},
		{
			ID: 1, Name: "BTC Hot Wallet", Kind: "hot",
			Currencies: []string{"btc"}, Address: "address", Gateway: "opendax_cloud",
			MaxBalance: "0.0", Balance: "0.0", BlockchainKey: blockchainKey, Status: "active",
		},
		{
			ID: 1, Name: "ETH Deposit Wallet", Kind: "eth",
			Currencies: []string{"eth", "link"}, Address: "address", Gateway: "opendax_cloud",
			MaxBalance: "0.0", Balance: "0.0", BlockchainKey: blockchainKey, Status: "active",
		},
		{
			ID: 1, Name: "ETH Hot Wallet", Kind: "hot",
			Currencies: []string{"eth", "link"}, Address: "address", Gateway: "opendax_cloud",
			MaxBalance: "0.0", Balance: "0.0", BlockchainKey: blockchainKey, Status: "active",
		},
	}

	// Full match
	currencies := []string{"btc"}
	actualResult := findCurrenciesInWallets(wallets, currencies)
	expectedResult := map[string][]*peatio.Wallet{}
	expectedResult["full"] = append(expectedResult["full"], wallets[0])
	expectedResult["full"] = append(expectedResult["full"], wallets[1])
	assert.Equal(t, reflect.DeepEqual(actualResult, expectedResult), true)

	// Partial match
	currencies = []string{"eth"}
	actualResult = findCurrenciesInWallets(wallets, currencies)
	expectedResult = map[string][]*peatio.Wallet{}
	expectedResult["partial"] = append(expectedResult["partial"], wallets[2])
	expectedResult["partial"] = append(expectedResult["partial"], wallets[3])
	assert.Equal(t, reflect.DeepEqual(actualResult, expectedResult), true)

	// None Match
	currencies = []string{"xrp"}
	actualResult = findCurrenciesInWallets(wallets, currencies)
	expectedResult = map[string][]*peatio.Wallet{}
	expectedResult["none"] = []*peatio.Wallet{}
	assert.Equal(t, reflect.DeepEqual(actualResult, expectedResult), true)
}
