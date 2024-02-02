package main

import (
	"ewallet/internal/api/v1/wallet"
	"ewallet/internal/database"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	DBConfig := database.NewDBConfig()

	dbHandler, err := database.ConnectDB(DBConfig)
	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return
	}

	defer dbHandler.Close()

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/wallet", wallet.CreateWalletHandler(dbHandler)).Methods("POST")
	router.HandleFunc("/api/v1/wallet/{walletId}/send", wallet.SendMoneyHandler(dbHandler)).Methods("POST")
	router.HandleFunc("/api/v1/wallet/{walletId}/history", wallet.GetTransactionHistoryHandler(dbHandler))
	router.HandleFunc("/api/v1/wallet/{walletId}", wallet.GetWalletStateHandler(dbHandler))

	http.Handle("/", router)

	port := os.Getenv("APP_PORT")

	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

	if err != nil {
		fmt.Println("Error when starting the server:", err)
	}
}
