package daemons

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openware/pkg/mngapi/peatio"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func TestFetchMarketsFromOpenfinexCloudSuccess(t *testing.T) {
	mockedResponse := []byte(`[{"id":"omgusdt","name":"OMG/USDT","base_unit":"omg","quote_unit":"usdt","state":"enabled","amount_precision":2,"price_precision":4,"min_price":"0.0037","max_price":"3687.4597","min_amount":"0.01","position":9},{"id":"uniusdt","name":"UNI/USDT","base_unit":"uni","quote_unit":"usdt","state":"enabled","amount_precision":2,"price_precision":4,"min_price":"0.0037","max_price":"3670","min_amount":"0.01","position":12}]`)
	peatioMockedResponse := []byte(`{"id":"omgusdt","name":"OMG/USDT","base_unit":"omg","quote_unit":"usdt","state":"enabled","amount_precision":2,"price_precision":4,"min_price":"0.0037","max_price":"3687.4597","min_amount":"0.01","position":9}`)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/opx/markets", func(res http.ResponseWriter, req *http.Request) {
		res.Write(mockedResponse)
	})
	mux.HandleFunc("/api/v2/peatio/management/markets/new", func(res http.ResponseWriter, req *http.Request) {
		res.Write(peatioMockedResponse)
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	peatioClient, err := peatio.New(fmt.Sprintf("%s/api/v2/peatio/management", ts.URL), jwtIssuer, jwtAlgo, jwtPrivateKey)
	require.NoError(t, err)

	app := initApp()
	app.Conf.Opendax.Addr = ts.URL
	res := FetchMarketsFromOpenfinexCloud(peatioClient, app.Conf.Opendax)
	assert.Equal(t, res, nil)
}

func TestFetchMarketsFromOpenfinexCloudUnmarshalError(t *testing.T) {
	mockedResponse := []byte(`{"error":"json: cannot unmarshal object into Go value of type []deamons.MarketResponse"}`)
	response := "json: cannot unmarshal object into Go value of type []daemons.MarketResponse"
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v2/opx/markets", func(res http.ResponseWriter, req *http.Request) {
		res.Write(mockedResponse)
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	peatioClient, err := peatio.New(fmt.Sprintf("%s/api/v2/peatio/management", ts.URL), jwtIssuer, jwtAlgo, jwtPrivateKey)
	require.NoError(t, err)

	app := initApp()
	app.Conf.Opendax.Addr = ts.URL
	res := FetchMarketsFromOpenfinexCloud(peatioClient, app.Conf.Opendax)
	require.Error(t, res)
	assert.Equal(t, res.Error(), response)
}

func TestFetchMarketsFromOpenfinexCloudHostError(t *testing.T) {
	response := "Unexpected status: 404"
	mux := http.NewServeMux()

	ts := httptest.NewServer(mux)
	defer ts.Close()

	peatioClient, err := peatio.New(fmt.Sprintf("%s/api/v2/peatio/management", ts.URL), jwtIssuer, jwtAlgo, jwtPrivateKey)
	require.NoError(t, err)

	app := initApp()
	app.Conf.Opendax.Addr = ts.URL
	res := FetchMarketsFromOpenfinexCloud(peatioClient, app.Conf.Opendax)
	require.Error(t, res)
	assert.Equal(t, res.Error(), response)
}
