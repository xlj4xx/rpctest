package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"
)

func main() {

	// url_alchemy := "https://eth-mainnet.g.alchemy.com/v2/OEhsfiW248rLojHL9hJ1ivPIKo8dryFQ"
	// url_quicknode := "https://snowy-wild-road.discover.quiknode.pro/df8aea5d494de1cf8dcfc7204601f86d2c360e36/"
	url := "https://mainnet.chainnodes.org/1840bd6e-0b4c-4aaa-9772-e1b15c4f29ab"

	payload := strings.NewReader("{\"id\":10,\"jsonrpc\":\"2.0\",\"params\":[\"0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045\",\"latest\"],\"method\":\"eth_getBalance\"}")
	for i := 0; i < 3; i++ {
		req, _ := http.NewRequest("POST", url, payload)
		req.Header.Add("accept", "application/json")
		req.Header.Add("content-type", "application/json")

		startTime := time.Now()

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		body, _ := ioutil.ReadAll(res.Body)
		convBodytoresult(body)

		elapsedTime := time.Since(startTime)
		fmt.Printf("\nTime taken: %s", elapsedTime)
		time.Sleep(1 * time.Second)
	}

}

type ethGetBalanceResponse struct {
	Result string `json:"result"`
}

func convBodytoresult(body []byte) {
	var response ethGetBalanceResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}
	convHextoEth(response.Result)
}

func convHextoEth(result string) {
	// convert the hexadecimal string to a big integer
	wei, ok := new(big.Int).SetString(result, 0)
	if !ok {
		panic("invalid hexadecimal string")
	}
	ether := new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e18))
	fmt.Printf("%s eth", ether.Text('f', -1))
}
