package getblockapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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

func (gb *GetBlockApi) GetHeadBlockNum() (uint64, error) {
	// Form request
	requestId := "daegon" + strconv.Itoa(int(time.Now().Unix()))

	body := fmt.Sprintf(`
	{
		"jsonrpc": "2.0",
		"method": "eth_blockNumber",
		"params": [],
		"id": "%s"
	}
	`, requestId)

	req, err := http.NewRequest(http.MethodPost, GetBlockEndpointETH, strings.NewReader(body))
	if err != nil {
		return 0, utils.WrapErr(caller, "Create request error", err)
	}
	req.Header.Add("x-api-key", gb.config.Key)
	req.Header.Add("Content-Type", "application/json'")

	// Make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, utils.WrapErr(caller, "Get response body error", err)
	}

	// Read response
	defer resp.Body.Close()
	result := &RespBlockNum{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return 0, utils.WrapErr(caller, "Read response body error", err)
	}

	blockNum, err := utils.HexToUint64(result.Result)
	if err != nil {
		return 0, utils.WrapErr(caller, "Convert hex to uint64 error", err)
	}

	return blockNum, nil
}

func (gb *GetBlockApi) GetBlockByNum(n uint64) (*models.Block, error) {
	// Form request
	requestId := "daegon" + strconv.Itoa(int(time.Now().Unix()))

	body := fmt.Sprintf(`
	{
		"jsonrpc": "2.0",
		"method": "eth_getBlockByNumber",
		"params": ["%s", true],
		"id": "%s"
	}
	`, utils.Uint64ToHex(n), requestId)

	req, err := http.NewRequest(http.MethodPost, GetBlockEndpointETH, strings.NewReader(body))
	if err != nil {
		return nil, utils.WrapErr(caller, "Create request error", err)
	}
	req.Header.Add("x-api-key", gb.config.Key)
	req.Header.Add("Content-Type", "application/json'")

	// Make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, utils.WrapErr(caller, "Get response body error", err)
	}

	// Read response
	defer resp.Body.Close()
	result := &RespBlockByNum{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, utils.WrapErr(caller, "Read response body error", err)
	}

	return &result.Result, nil
}

type RespBlockNum struct {
	Result string `json:"result"`
}

type RespBlockByNum struct {
	Result models.Block `json:"result"`
}
