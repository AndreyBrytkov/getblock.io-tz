package getblockapi

import (
	"context"
	"math/big"
	"net/http"

	"github.com/AndreyBrytkov/getblock.io-tz/internal/adapter"
	"github.com/AndreyBrytkov/getblock.io-tz/internal/models"
	"github.com/AndreyBrytkov/getblock.io-tz/pkg/utils"
	"github.com/ubiq/go-ubiq/common/hexutil"

	jsonrpc "github.com/ybbus/jsonrpc/v3"
)

const (
	caller = "getblock-api"

	GetBlockEndpointETH  = "https://eth.getblock.io/mainnet/"
	RpcMethodBlockNumber = "eth_blockNumber"
	RpcMethodGetBlock    = "eth_getBlockByNumber"
)

type GetBlockApi struct {
	client jsonrpc.RPCClient
	logger *utils.MyLogger
	config *models.GetBlockConfig
}

func GetGetBlockApi(logger *utils.MyLogger, config *models.GetBlockConfig) adapter.GetBlockApi {
	options := &jsonrpc.RPCClientOpts{
		HTTPClient: http.DefaultClient,
		CustomHeaders: map[string]string{
			"x-api-key":    config.Key,
			"Content-Type": "application/json",
		},
	}

	client := jsonrpc.NewClientWithOpts(GetBlockEndpointETH, options)

	return &GetBlockApi{
		client: client,
		logger: logger,
		config: config,
	}
}

func (gb *GetBlockApi) GetHeadBlockNum() (big.Int, error) {
	// gb.logger.Info(caller, "getting lastest block number...")
	blockNum := big.NewInt(0)

	resp, err := gb.client.Call(context.Background(), RpcMethodBlockNumber)
	if err != nil {
		return *blockNum, utils.WrapErr(caller, "JSON-RPC call 'eth_blockNumber' error", err)
	}

	if resp.Error != nil {
		return *blockNum, utils.WrapErr(caller, "JSON-RPC call 'eth_blockNumber' error", err)
	}

	numStr, err := resp.GetString()
	if resp.Error != nil {
		return *blockNum, utils.WrapErr(caller, "get block number from response error", err)
	}

	blockNum, err = hexutil.DecodeBig(numStr)
	if err != nil {
		return *blockNum, utils.WrapErr(caller, "Decode blockNum to big.Int error", err)
	}

	gb.logger.Debug(caller, "last block number is '%s'", hexutil.EncodeBig(blockNum))
	return *blockNum, nil
}

func (gb *GetBlockApi) GetBlockByNum(n big.Int) (*models.Block, error) {
	// gb.logger.Info(caller, "getting block number '%s'", n.String())
	resp, err := gb.client.Call(context.Background(), RpcMethodGetBlock, hexutil.EncodeBig(&n), true)
	if err != nil {
		return nil, utils.WrapErr(caller, "JSON-RPC call 'eth_getBlockByNumber' error", err)
	}

	if resp.Error != nil {
		return nil, utils.WrapErr(caller, "JSON-RPC call 'eth_getBlockByNumber' error", err)
	}

	blockJson := new(BlockJson)
	err = resp.GetObject(blockJson)
	if resp.Error != nil {
		return nil, utils.WrapErr(caller, "get BlockJson object from response error", err)
	}

	block := new(models.Block)
	err = decodeBlockJSON(blockJson, block)
	if err != nil {
		return nil, utils.WrapErr(caller, "decode block error", err)
	}

	return block, nil
}

type BlockJson struct {
	BaseFeePerGas    string            `json:"baseFeePerGas"`
	Difficulty       string            `json:"difficulty"`
	ExtraData        string            `json:"extraData"`
	GasLimit         string            `json:"gasLimit"`
	GasUsed          string            `json:"gasUsed"`
	Hash             string            `json:"hash"`
	LogsBloom        string            `json:"logsBloom"`
	Miner            string            `json:"miner"`
	MixHash          string            `json:"mixHash"`
	Nonce            string            `json:"nonce"`
	Number           string            `json:"number"`
	ParentHash       string            `json:"parentHash"`
	ReceiptsRoot     string            `json:"receiptsRoot"`
	Sha3Uncles       string            `json:"sha3Uncles"`
	Size             string            `json:"size"`
	StateRoot        string            `json:"stateRoot"`
	Timestamp        string            `json:"timestamp"`
	TotalDifficulty  string            `json:"totalDifficulty"`
	Transactions     []TransactionJson `json:"transactions"`
	TransactionsRoot string            `json:"transactionsRoot"`
	Uncles           []interface{}     `json:"uncles"`
}

type TransactionJson struct {
	BlockHash            string        `json:"blockHash"`
	BlockNumber          string        `json:"blockNumber"`
	From                 string        `json:"from"`
	Gas                  string        `json:"gas"`
	GasPrice             string        `json:"gasPrice"`
	MaxFeePerGas         string        `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas string        `json:"maxPriorityFeePerGas,omitempty"`
	Hash                 string        `json:"hash"`
	Input                string        `json:"input"`
	Nonce                string        `json:"nonce"`
	To                   string        `json:"to"`
	TransactionIndex     string        `json:"transactionIndex"`
	Value                string        `json:"value"`
	Type                 string        `json:"type"`
	AccessList           []interface{} `json:"accessList,omitempty"`
	ChainID              string        `json:"chainId,omitempty"`
	V                    string        `json:"v"`
	R                    string        `json:"r"`
	S                    string        `json:"s"`
}

func decodeBlockJSON(blockJson *BlockJson, dst *models.Block) error {
	blockNum, err := hexutil.DecodeBig(blockJson.Number)
	if err != nil {
		return utils.WrapErr(caller, "Decode Number to big.Int error", err)
	}
	dst.Number = *blockNum

	// Convert all numbers
	for _, tx := range blockJson.Transactions {
		// Index
		idx, err := hexutil.DecodeBig(tx.TransactionIndex)
		if err != nil {
			return utils.WrapErr(caller, "Decode transaction Index to big.Int error", err)
		}
		// Value
		value, err := hexutil.DecodeBig(tx.Value)
		if err != nil {
			return utils.WrapErr(caller, "Decode transaction Value to big.Int error", err)
		}
		// Gas
		gas, err := hexutil.DecodeBig(tx.Gas)
		if err != nil {
			return utils.WrapErr(caller, "Decode transaction Gas to big.Int error", err)
		}
		// Price WEI per gas
		gasPrice, err := hexutil.DecodeBig(tx.GasPrice)
		if err != nil {
			return utils.WrapErr(caller, "Decode transaction Gas Price to big.Int error", err)
		}

		// Calculate total gas price
		gasTotal := big.NewInt(0)
		gasTotal.Mul(gas, gasPrice)

		// Form new tx
		newTx := models.Trasaction{
			BlockNumber: *blockNum,
			Idx:         *idx,
			From:        tx.From,
			To:          tx.To,
			Value:       *value,
			Gas:         *gas,
			GasPrice:    *gasPrice,
			GasTotal:    *gasTotal,
		}

		dst.Transactions = append(dst.Transactions, newTx)
	}

	return nil
}
