package restserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/AndreyBrytkov/getblock.io-tz/internal/adapter"
	"github.com/AndreyBrytkov/getblock.io-tz/internal/models"
	"github.com/AndreyBrytkov/getblock.io-tz/pkg/utils"
)

const caller = "server"

// Main object
type RestServer struct {
	logger  *utils.MyLogger
	config  *models.ServerConfig
	wg      *sync.WaitGroup
	usecase adapter.Usecase
	mux     *http.ServeMux
	address string
}

// Getter
func NewRestApi(logger *utils.MyLogger, config *models.ServerConfig, wg *sync.WaitGroup, usecase adapter.Usecase) adapter.Rest {
	rs := &RestServer{
		logger:  logger,
		config:  config,
		wg:      wg,
		usecase: usecase,
		mux:     http.NewServeMux(),
		address: fmt.Sprintf("%s:%d", config.Address, config.Port),
	}

	rs.mux.HandleFunc("/", rs.maxDeltaWalletHandler)

	return rs
}

func (rs *RestServer) Run() {
	defer rs.wg.Done()
	rs.logger.Info(caller, "started...")
	err := http.ListenAndServe(rs.address, rs.mux)
	if err != nil {
		rs.logger.Fatal(utils.WrapErr(caller, "", err))
	}
}

func (rs *RestServer) maxDeltaWalletHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		msg := fmt.Sprintf("method '%s' not allowed", req.Method)
		rs.logger.Error(utils.WrapErr(caller, "", errors.New(msg)))
		http.Error(w, msg + ", expect method GET", http.StatusMethodNotAllowed)
		return
	}

	wallet, delta, err := rs.usecase.GetMaxBalanceDeltaWallet()
	if err != nil {
		msg := "get max balance delta wallet error"
		rs.logger.Error(utils.WrapErr(caller, msg, err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	result := Result{
		BalanceDeltaWei: delta.String(),
		Wallet: wallet,
	}

	resultJson, err := json.Marshal(result)
	if err != nil {
		msg := "marshal result error"
		rs.logger.Error(utils.WrapErr(caller, msg, err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resultJson)
}

type Result struct {
	BalanceDeltaWei string `json:"balance_delta_wei"`
	Wallet          string `json:"wallet"`
}

