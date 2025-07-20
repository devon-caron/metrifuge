package handlers

import (
	"encoding/json"
	"github.com/devon-caron/metrifuge/api/errhandler"
	"net/http"

	"github.com/devon-caron/metrifuge/api/internal/tools"
	api "github.com/devon-caron/metrifuge/api/types"
	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
)

func GetCoinBalance(w http.ResponseWriter, r *http.Request) {
	var params = api.CoinBalanceParams{}
	var decoder = schema.NewDecoder()
	var err error

	err = decoder.Decode(&params, r.URL.Query())

	if err != nil {
		log.Error(err)
		errhandler.InternalErrorHandler(w)
		return
	}

	log.Println("url query decoded successfully")

	var database *tools.DatabaseInterface
	database, err = tools.NewDatabase()
	if err != nil {
		errhandler.InternalErrorHandler(w)
		return
	}

	log.Println("database created successfully")

	var tokenDetails *tools.CoinDetails = (*database).GetUserCoins(params.Username)
	if tokenDetails == nil {
		log.Error(err)
		errhandler.InternalErrorHandler(w)
		return
	}

	var response = api.CoinBalanceResponse{
		Balance: (*tokenDetails).Coins,
		Code:    http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		errhandler.InternalErrorHandler(w)
		return
	}

	log.Println("response created and encoded successfully")
}
