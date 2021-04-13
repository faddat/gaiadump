package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/nitishm/go-rejson"
	"io"
	"log"
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
			LatestBlockHeight   int       `json:"latest_block_height"`
			LatestBlockTime     time.Time `json:"latest_block_time"`
			EarliestBlockHash   string    `json:"earliest_block_hash"`
			EarliestAppHash     string    `json:"earliest_app_hash"`
			EarliestBlockHeight int       `json:"earliest_block_height"`
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

	rh := rejson.NewReJSONHandler()
	cli := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	defer func() {
		if err := cli.FlushAll().Err(); err != nil {
			log.Fatalf("goredis - failed to flush: %v", err)
		}
		if err := cli.Close(); err != nil {
			log.Fatalf("goredis - failed to communicate to redis-server: %v", err)
		}
	}()
	rh.SetGoRedisClient(cli)

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
			GetBlock(i, rh)
		}
		time.Sleep(5 * time.Second)
	}

}

func GetBlock(block int, rh *rejson.Handler) {

	//http://localhost:26657/block?height=5272289
	//res, err := rh.JSONSet()

	res, err := rh.JSONSet("student", ".", student)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res)

}
