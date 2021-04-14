package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/nitishm/go-rejson"
)

// Struct generated here: https://mholt.github.io/json-to-go/ because the clients don't quite do the trick.
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
			LatestBlockHeight   string    `json:"latest_block_height"`
			LatestBlockTime     time.Time `json:"latest_block_time"`
			EarliestBlockHash   string    `json:"earliest_block_hash"`
			EarliestAppHash     string    `json:"earliest_app_hash"`
			EarliestBlockHeight string    `json:"earliest_block_height"`
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
		status := GetStatus()
		res, err := rh.JSONSet(status.Result.NodeInfo.Network, ".", status)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Redis response", res)

		ebs := status.Result.SyncInfo.EarliestBlockHeight
		lbs := status.Result.SyncInfo.LatestBlockHeight

		eb, err := strconv.Atoi(ebs)
		lb, err := strconv.Atoi(lbs)

		fmt.Println(eb, "to", lb)
		for i := eb; i < lb; i++ {
			GetBlock(i, rh)
		}

		// this is a rudimentary way of checking for new blocks and should likely be improved.
		time.Sleep(5 * time.Second)
	}

}

// GetBlock grabs a block from a cosmos-sdk chain.
func GetBlock(block int, rh *rejson.Handler) {

	// http://localhost:26657/block?height=5272289
	// res, err := rh.JSONSet()
	// CHAIN should be like: https://rpc.testnet1.test.gravitydex.io/
	base := os.Getenv("CHAIN") + "block?height="
	fmt.Println(block)
	requrl := base + strconv.Itoa(block)

	fmt.Println("request url", requrl)

	resp, err := http.Get(requrl)
	if err != nil {
		fmt.Println("ERROR!", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR!", err)
	}

	fmt.Println("body", body)
	res, err := rh.JSONSet("height", "status", body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res)

}

// GetStatus() grabs the status from a Cosmos-SDK chain.
func GetStatus() (status Status) {
	resp, err := http.Get(os.Getenv(CHAIN))
	if err != nil {
		fmt.Println("ERROR!", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR!", err)
	}

	json.Unmarshal(body, &status)
	return
}
