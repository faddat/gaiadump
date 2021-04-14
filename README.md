# sdkeater

This will get every block of any comos sdk blockchain, then every transaction, and store JSON returns to [ReJSON](https://github.com/RedisJSON/RedisJSON).

If we can't do that, ecosystem-wide, that's bad.

The comparable approach is tracelistener by [@gsora](https://github.com/gsora).

If it is necessary to use a unix socket to pipe a terminal-based dump into cockroachdb, instead of simply using user-facing APIs, that is an illustration of a problem.

Sdkeater is very stupid.  It thinks only like this:

```go
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
```

Then, it focuses on only two things:

`EarliestBlockHeight` to `LatestBlockHeight`

Sdkeater is done eating when every block from `EarliestBlockHeight` to `LatestBlockHeight` is put into a waiting Redis with JSON module. 

It is not thinking at all, or worrying about the exact types of what it is eating.

![Screenshot from 2021-04-13 15-54-06](https://user-images.githubusercontent.com/7142025/114690148-a2d0e900-9d40-11eb-920a-bba6b4a1b285.png)

If it was worried about types, it would probably look more like https://github.com/forbole/bdjuno: a re-implementation of the chain being consumed to ensure that types can be known.


