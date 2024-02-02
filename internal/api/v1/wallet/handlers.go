package wallet

import (
	"encoding/json"
	"ewallet/internal/database"
	"ewallet/internal/services"
	"net/http"

	"github.com/gorilla/mux"
)

func GetWalletStateHandler(dbh *database.DBHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		walletID := mux.Vars(r)["walletId"]
		if walletID == "" {
			http.Error(w, "Invalid wallet ID", http.StatusBadRequest)
			return
		}

		balance, err := dbh.GetWalletBalance(walletID)
		if err != nil {
			http.Error(w, "Wallet Not Found", http.StatusNotFound)
			return
		}

		response := map[string]interface{}{
			"id":      walletID,
			"balance": balance,
		}

		json.NewEncoder(w).Encode(response)
	}
}

func GetTransactionHistoryHandler(dbh *database.DBHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		walletID := mux.Vars(r)["walletId"]
		if walletID == "" {
			http.Error(w, "Invalid wallet ID", http.StatusBadRequest)
			return
		}

		_, err := dbh.GetWalletBalance(walletID)
		if err != nil {
			http.Error(w, "Wallet Not Found", http.StatusNotFound)
			return
		}

		transactions, err := dbh.GetTransactionHistory(walletID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

		if transactions != nil {
			json.NewEncoder(w).Encode(transactions)
			return
		}

		json.NewEncoder(w).Encode("History is empty")
	}
}

func CreateWalletHandler(dbh *database.DBHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		newWallet, err := services.CreateWallet(dbh)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newWallet)
	}
}

func SendMoneyHandler(dbh *database.DBHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		senderWalletID := vars["walletId"]
		if senderWalletID == "" {
			http.Error(w, "Invalid wallet ID", http.StatusBadRequest)
			return
		}

		var sendRequest struct {
			To     string  `json:"to"`
			Amount float64 `json:"amount"`
		}

		err := json.NewDecoder(r.Body).Decode(&sendRequest)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = services.TransferMoney(dbh, senderWalletID, sendRequest.To, sendRequest.Amount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
