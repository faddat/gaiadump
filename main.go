package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nitishm/go-rejson"
	goredis "github.com/go-redis/redis/v8"
	"io"
	"net/http"
	"time"
)

type Status struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		NodeInfo struct {
			ProtocolVersion struct {
				P2P   string `json:"p2p"`
				Block string `json:"block"`
				App   string `json:"app"`
			} `json:"protocol_version"`
			ID         string `json:"id"`
			ListenAddr string `json:"listen_addr"`
			Network    string `json:"network"`
			Version    string `json:"version"`
			Channels   string `json:"channels"`
			Moniker    string `json:"moniker"`
			Other      struct {
				TxIndex    string `json:"tx_index"`
				RPCAddress string `json:"rpc_address"`
			} `json:"other"`
		} `json:"node_info"`
		SyncInfo struct {
			LatestBlockHash     string    `json:"latest_block_hash"`
			LatestAppHash       string    `json:"latest_app_hash"`
			LatestBlockHeight   int    `json:"latest_block_height"`
			LatestBlockTime     time.Time `json:"latest_block_time"`
			EarliestBlockHash   string    `json:"earliest_block_hash"`
			EarliestAppHash     string    `json:"earliest_app_hash"`
			EarliestBlockHeight int    `json:"earliest_block_height"`
			EarliestBlockTime   time.Time `json:"earliest_block_time"`
			CatchingUp          bool      `json:"catching_up"`
		} `json:"sync_info"`
		ValidatorInfo struct {
			Address string `json:"address"`
			PubKey  struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"pub_key"`
			VotingPower string `json:"voting_power"`
		} `json:"validator_info"`
	} `json:"result"`
}






func main() {

	flag.Parse()

	for {
		resp, err := http.Get("http://localhost:26657/status")
		if err != nil {
			fmt.Println("ERROR!", err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("ERROR!", err)
		}
		var status Status
		json.Unmarshal(body, &status)
		first := status.Result.SyncInfo.EarliestBlockHeight
		last := status.Result.SyncInfo.LatestBlockHeight
		fmt.Println(last)

		for i := first; i < last; i++ {
			GetBlock(i)
		}

		time.Sleep(5 * time.Second)





	}

}

func GetBlock(block int) {

	rh := rejson.NewReJSONHandler()
	cli := goredis.NewClient







}