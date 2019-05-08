package types

import (
	"time"
)

//Receipt receipt information
type Receipt struct {
	tableName         struct{}      `sql:"rpt_tbl"`
	BlockHash         string        `sql:"blockHash"`
	BlockNumber       string        `sql:"blockNumber"`
	ContractAddress   string        `sql:"contractAddress"`
	CumulativeGasUsed string        `sql:"cumulativeGasUsed"`
	From              string        `sql:"from"`
	GasUsed           string        `sql:"gasUsed"`
	Logs              []interface{} `sql:"logs"`
	LogsBloom         string        `sql:"logsBloom"`
	Status            string        `sql:"status"`
	To                string        `sql:"to"`
	TransactionHash   string        `sql:"transactionHash"`
	TransactionIndex  string        `sql:"transactionIndex"`
}

//Transaction transaction information
type Transaction struct {
	tableName        struct{}  `sql:"txn_tbl"`
	ID               int64     `sql:"id" json:"id"`
	BlockNumber      string    `sql:"blocknumber" json:"blockNumber"`
	Gas              string    `sql:"gas" json:"gas"`
	Nonce            string    `sql:"nonce" json:"nonce"`
	R                string    `sql:"r" json:"r"`
	S                string    `sql:"s" json:"s"`
	GasPrice         string    `sql:"gasprice" json:"gasPrice"`
	Input            string    `sql:"input" json:"input"`
	Value            string    `sql:"value" json:"value"`
	To               string    `sql:"txn_to" json:"to"`
	TransactionIndex string    `sql:"transactionindex" json:"transactionIndex"`
	V                string    `sql:"v" json:"v"`
	BlockHash        string    `sql:"blockhash" json:"blockHash"`
	From             string    `sql:"txn_from" json:"from"`
	Hash             string    `sql:"hash" json:"hash"`
	Timestamp        time.Time `sql:"ts" json:"timestamp"`
	Status           string    `sql:"-" json:"status"`
}

//Block etherenum block information
type Block struct {
	Difficulty       string        `json:"difficulty"`
	ExtraData        string        `json:"extraData"`
	GasLimit         string        `json:"gasLimit"`
	GasUsed          string        `json:"gasUsed"`
	Hash             string        `json:"hash"`
	LogsBloom        string        `json:"logsBloom"`
	Miner            string        `json:"miner"`
	MixHash          string        `json:"mixHash"`
	Nonce            string        `json:"nonce"`
	Number           string        `json:"number"`
	ParentHash       string        `json:"parentHash"`
	ReceiptsRoot     string        `json:"receiptsRoot"`
	Sha3Uncles       string        `json:"sha3Uncles"`
	Size             string        `json:"size"`
	StateRoot        string        `json:"stateRoot"`
	Timestamp        string        `json:"timestamp"`
	TotalDifficulty  string        `json:"totalDifficulty"`
	Transactions     []Transaction `json:"transactions"`
	TransactionsRoot string        `json:"transactionsRoot"`
	Uncles           []interface{} `json:"uncles"`
	Validator        string        `json:"validator"`
}

//TransactionProcessRecord records last block each worker processed
type TransactionProcessRecord struct {
	Workers    int     `json:"workers"`
	LastBlocks []int64 `json:"lastBlocks"`
}
