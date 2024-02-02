package database

import (
	"database/sql"
	"errors"
	"ewallet/internal/models"
	"fmt"
	"time"
)

type DBHandler struct {
	db *sql.DB
}

func ConnectDB(config *DBConfig) (*DBHandler, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.host, config.port, config.username, config.password, config.database)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Failed to connect to database")
		return nil, err
	}

	fmt.Println("Connected to the database")

	return &DBHandler{db: db}, nil
}

func (dbh *DBHandler) Close() {
	if dbh.db != nil {
		dbh.db.Close()
	}
}

func (dbh *DBHandler) Begin() (*sql.Tx, error) {
	if dbh.db != nil {
		return dbh.db.Begin()
	}

	return nil, errors.New("database does not exist")
}

func (dbh *DBHandler) InsertWalletToDB(wallet *models.Wallet) error {
	query := "INSERT INTO Wallet (id, balance) VALUES ($1, $2)"
	_, err := dbh.db.Exec(query, wallet.ID, wallet.Balance)
	if err != nil {
		return fmt.Errorf("failed to insert wallet into DB: %v", err)
	}
	return nil
}

func (dbh *DBHandler) GetWalletBalance(walletID string) (float64, error) {
	var balance float64
	err := dbh.db.QueryRow("SELECT balance FROM Wallet WHERE ID = $1", walletID).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (dbh *DBHandler) GetWalletBalanceWithTx(walletID string, tx *sql.Tx) (float64, error) {
	var balance float64
	err := tx.QueryRow("SELECT balance FROM Wallet WHERE ID = $1", walletID).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (dbh *DBHandler) UpdateWalletBalance(walletID string, newBalance float64) error {
	_, err := dbh.db.Exec("UPDATE Wallet SET balance = $1 WHERE ID = $2", newBalance, walletID)
	return err
}

func (dbh *DBHandler) UpdateWalletBalanceWithTx(walletID string, newBalance float64, tx *sql.Tx) error {
	_, err := tx.Exec("UPDATE Wallet SET balance = $1 WHERE ID = $2", newBalance, walletID)
	return err
}

func (dbh *DBHandler) LogTransaction(from, to string, amount float64) error {
	_, err := dbh.db.Exec("INSERT INTO History (time, sender, recipient, amount) VALUES ($1, $2, $3, $4)",
		time.Now(), from, to, amount)
	return err
}

func (dbh *DBHandler) LogTransactionWithTx(from, to string, amount float64, tx *sql.Tx) error {
	_, err := tx.Exec("INSERT INTO History (time, sender, recipient, amount) VALUES ($1, $2, $3, $4)",
		time.Now(), from, to, amount)
	return err
}

func (dbh *DBHandler) GetTransactionHistory(walletID string) ([]models.Transaction, error) {
	query := `
		SELECT time, sender, recipient, amount
		FROM History
		WHERE sender = $1 OR recipient = $1
	`

	rows, err := dbh.db.Query(query, walletID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		if err := rows.Scan(&transaction.Time, &transaction.From, &transaction.To, &transaction.Amount); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
