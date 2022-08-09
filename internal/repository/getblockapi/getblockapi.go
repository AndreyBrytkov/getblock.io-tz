package getblockapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"

	"github.com/AndreyBrytkov/getblock.io-tz/internal/adapter"
	"github.com/AndreyBrytkov/getblock.io-tz/internal/models"
	"github.com/AndreyBrytkov/getblock.io-tz/pkg/utils"
	"github.com/ubiq/go-ubiq/common/hexutil"
)

const (
	caller = "getblock-api"

	GetBlockEndpointETH = "https://eth.getblock.io/mainnet/"
)

type GetBlockApi struct {
	logger *utils.MyLogger
	config *models.GetBlockConfig
}

func GetGetBlockApi(logger *utils.MyLogger, config *models.GetBlockConfig) adapter.GetBlockApi {
	return &GetBlockApi{
		logger: logger,
		config: config,
	}
}

func (gb *GetBlockApi) GetHeadBlockNum() (big.Int, error) {
	gb.logger.Debug(caller, "get lastest block...")
	blockNum := big.NewInt(0)

	// Form request
	requestId := time.Now().Unix()

	body := []byte(fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"method": "eth_blockNumber",
		"params": [],
		"id": "daegon%d"
	}`, requestId))

	req, err := http.NewRequest(http.MethodPost, GetBlockEndpointETH, bytes.NewBuffer(body))
	if err != nil {
		return *blockNum, utils.WrapErr(caller, "Create request error", err)
	}
	req.Header.Set("x-api-key", gb.config.Key)
	req.Header.Set("Content-Type", "application/json")

	// Make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return *blockNum, utils.WrapErr(caller, "Get response body error", err)
	} else if resp.StatusCode != http.StatusOK {
		return *blockNum, utils.WrapErr(caller, "", fmt.Errorf("response %d %s", resp.StatusCode, resp.Status))
	}
	defer resp.Body.Close()

	// Read response
	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return *blockNum, utils.WrapErr(caller, "Read response body error", err)
	}
	gb.logger.Debug(caller, "response:\n%s", string(respBodyBytes))

	result := make(map[string]string)
	err = json.Unmarshal(respBodyBytes, &result)
	if err != nil {
		return *blockNum, utils.WrapErr(caller, "Unmarshal response body error", err)
	}

	numStr, ok := result["result"]
	if !ok {
		return *blockNum, errors.New("unmarshal response body error")
	}

	blockNum, err = hexutil.DecodeBig(numStr)
	if err != nil {
		return *blockNum, utils.WrapErr(caller, "Decode blockNum to big.Int error", err)
	}

	gb.logger.Debug(caller, "last block num = %s", hexutil.EncodeBig(blockNum))
	return *blockNum, nil
}

func (gb *GetBlockApi) GetBlockByNum(n big.Int) (*models.Block, error) {
	// Form request
	requestId := time.Now().Unix()
	gb.logger.Debug(caller, "get block num %s", hexutil.EncodeBig(&n))

	body := []byte(fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"method": "eth_getBlockByNumber",
		"params": ["%s", true],
		"id": "daegon%d"
	}`, hexutil.EncodeBig(&n), requestId))

	gb.logger.Debug(caller, "request body:\n%s", string(body))

	req, err := http.NewRequest(http.MethodPost, GetBlockEndpointETH, bytes.NewBuffer(body))
	if err != nil {
		return nil, utils.WrapErr(caller, "Create request error", err)
	}
	req.Header.Set("x-api-key", gb.config.Key)
	req.Header.Set("Content-Type", "application/json")

	// Make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, utils.WrapErr(caller, "Get response body error", err)
	} else if resp.StatusCode != http.StatusOK {
		return nil, utils.WrapErr(caller, "", fmt.Errorf("response %d %s", resp.StatusCode, resp.Status))
	}

	// Read response
	result := RespBlock{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	// respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, utils.WrapErr(caller, "Read response body error", err)
	}
	resp.Body.Close()

	gb.logger.Debug(caller, "response:\n%v", result)
	// err = json.Unmarshal(respBodyBytes, &result)
	// if err != nil {
	// 	return nil, utils.WrapErr(caller, "Unmarshal response body error", err)
	// }

	// Decode to models.block's bigint fields
	gb.logger.Debug(caller, "result.Result.Number = %s", result.Block.Number)
	blockNum, err := hexutil.DecodeBig(result.Block.Number)
	if err != nil {
		return nil, utils.WrapErr(caller, "Decode Number to big.Int error", err)
	}

	block := models.Block{
		Number: *blockNum,
	}

	// Convert all numbers
	for _, tx := range result.Block.Transactions {
		// Index
		idx, err := hexutil.DecodeBig(tx.Idx)
		if err != nil {
			return nil, utils.WrapErr(caller, "Decode transaction Index to big.Int error", err)
		}
		// Value
		value, err := hexutil.DecodeBig(tx.Value)
		if err != nil {
			return nil, utils.WrapErr(caller, "Decode transaction Value to big.Int error", err)
		}
		// Gas
		gas, err := hexutil.DecodeBig(tx.Gas)
		if err != nil {
			return nil, utils.WrapErr(caller, "Decode transaction Gas to big.Int error", err)
		}
		// Price WEI per gas
		gasPrice, err := hexutil.DecodeBig(tx.GasPrice)
		if err != nil {
			return nil, utils.WrapErr(caller, "Decode transaction Gas Price to big.Int error", err)
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
			GasTotal:    *gasTotal,
		}

		block.Transactions = append(block.Transactions, newTx)
	}

	return &block, nil
}

type RespBlock struct {
	RequestId string    `json:"id"`
	Error     string    `json:"error"`
	Block     BlockJson `json:"result"`
}
type BlockJson struct {
	Number       string           `json:"number"`
	Transactions []TrasactionJson `json:"transactions"`
}
type TrasactionJson struct {
	Idx      string `json:"transactionIndex"`
	From     string `json:"from"`
	To       string `json:"to"`
	Value    string `json:"value"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
}
