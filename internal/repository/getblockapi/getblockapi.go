package getblockapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/AndreyBrytkov/getblock.io-tz/internal/adapter"
	"github.com/AndreyBrytkov/getblock.io-tz/internal/models"
	"github.com/AndreyBrytkov/getblock.io-tz/pkg/utils"
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
	var blockNum big.Int

	// Form request
	requestId := "daegon376"

	body := []byte(fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"method": "eth_blockNumber",
		"params": [],
		"id": "%s"
	}`, requestId))

	req, err := http.NewRequest(http.MethodPost, GetBlockEndpointETH, bytes.NewBuffer(body))
	if err != nil {
		return blockNum, utils.WrapErr(caller, "Create request error", err)
	}
	req.Header.Set("x-api-key", gb.config.Key)
	req.Header.Set("Content-Type", "application/json")

	// Make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return blockNum, utils.WrapErr(caller, "Get response body error", err)
	}
	// gb.logger.Debug(caller, "eth_blockNumber response: %v", resp)
	defer resp.Body.Close()

	// Read response
	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return blockNum, utils.WrapErr(caller, "Read response body error", err)
	}

	result := RespBlockNum{}
	err = json.Unmarshal(respBodyBytes, &result)
	if err != nil {
		return blockNum, utils.WrapErr(caller, "Unmarshal response body error", err)
	}

	err = blockNum.UnmarshalText([]byte(result.Result))
	if err != nil {
		return blockNum, utils.WrapErr(caller, "Convert blockNum to big.Int error", err)
	}

	return blockNum, nil
}

func (gb *GetBlockApi) GetBlockByNum(n big.Int) (*models.Block, error) {
	// Form request
	requestId := "daegon376"

	body := []byte(fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"method": "eth_getBlockByNumber",
		"params": ["%x", true],
		"id": "%s"
	}`, n, requestId))

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
	}
	defer resp.Body.Close()

	// Read response
	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, utils.WrapErr(caller, "Read response body error", err)
	}
	result := RespBlockByNum{}
	err = json.Unmarshal(respBodyBytes, &result)
	if err != nil {
		return nil, utils.WrapErr(caller, "Unmarshal response body error", err)
	}

	block := result.Result

	err = block.Number.UnmarshalText([]byte(block.NumberHex))
	if err != nil {
		return nil, utils.WrapErr(caller, "Convert NumberHex to big.Int error", err)
	}

	// Convert all numbers
	for i, tx := range block.Transactions {
		tx.BlockNumber = block.Number

		err = tx.Value.UnmarshalText([]byte(tx.ValueHex))
		if err != nil {
			return nil, utils.WrapErr(caller, "Convert ValueHex to big.Int error", err)
		}

		err = tx.Idx.UnmarshalText([]byte(tx.IdxHex))
		if err != nil {
			return nil, utils.WrapErr(caller, "Convert IdxHex to big.Int error", err)
		}

		var gasAmount big.Int
		err = gasAmount.UnmarshalText([]byte(tx.GasHex))
		if err != nil {
			return nil, utils.WrapErr(caller, "Convert GasHex to big.Int error", err)
		}

		var gasPrice big.Int
		err = gasPrice.UnmarshalText([]byte(tx.GasPriceHex))
		if err != nil {
			return nil, utils.WrapErr(caller, "Convert GasPriceHex to big.Int error", err)
		}

		tx.GasTotal.Mul(&gasAmount, &gasPrice)

		block.Transactions[i] = tx
	}

	return &block, nil
}

type RespBlockNum struct {
	Result string `json:"result"`
}

type RespBlockByNum struct {
	Result models.Block `json:"result"`
}
