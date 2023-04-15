package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	url_alchemy := "https://eth-mainnet.g.alchemy.com/v2/OEhsfiW248rLojHL9hJ1ivPIKo8dryFQ"
	url_quicknode := "https://snowy-wild-road.discover.quiknode.pro/df8aea5d494de1cf8dcfc7204601f86d2c360e36/"
	url_chainnode := "https://mainnet.chainnodes.org/1840bd6e-0b4c-4aaa-9772-e1b15c4f29ab"

	file, err := os.Create("rpc_speed.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"RPC Name", "Speed", "ETH Balance"}
	writer.Write(headers)

	for i := 0; i < 10; i++ {
		chainTime := requestRpc(url_chainnode, "chainnode")
		quickTime := requestRpc(url_quicknode, "quicknode")
		alchTime := requestRpc(url_alchemy, "alchemy")

		row1 := []string{"chainnode", chainTime}
		writer.Write(row1)
		row2 := []string{"quicknode", quickTime}
		writer.Write(row2)
		row3 := []string{"alchemy", alchTime}
		writer.Write(row3)
		// time.Sleep(1 * time.Second)

	}

}

func requestRpc(url string, name string) string {

	fmt.Printf("\nRequesting %s", name)
	payload := strings.NewReader("{\"id\":10,\"jsonrpc\":\"2.0\",\"params\":[\"0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045\",\"latest\"],\"method\":\"eth_getBalance\"}")
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
	return elapsedTime.String()
}

type ethGetBalanceResponse struct {
	Result string `json:"result"`
}

func convBodytoresult(body []byte) (ethGetBalanceResponse, error) {
	var response ethGetBalanceResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}
	if response.Result == "" {
		return response, fmt.Errorf("empty result")
	}

	convHextoEth(response.Result)
	return response, nil
}

func convHextoEth(result string) {
	// convert the hexadecimal string to a big integer
	wei, ok := new(big.Int).SetString(result, 0)
	if !ok {
		panic("invalid hexadecimal string")
	}
	ether := new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e18))
	fmt.Printf("\n%s eth", ether.Text('f', -1))
}
