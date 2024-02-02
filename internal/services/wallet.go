package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"ewallet/internal/database"
	"ewallet/internal/models"
)

func TransferMoney(dbh *database.DBHandler, senderID, receiverID string, amount float64) error {

	tx, err := dbh.Begin()
	if err != nil {
		return err
	}

	senderBalance, err := dbh.GetWalletBalanceWithTx(senderID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	receiverBalance, err := dbh.GetWalletBalanceWithTx(receiverID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if senderBalance < amount {
		return errors.New("insufficient funds")
	}

	senderBalance -= amount
	receiverBalance += amount

	err = dbh.UpdateWalletBalanceWithTx(senderID, senderBalance, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = dbh.UpdateWalletBalanceWithTx(receiverID, receiverBalance, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = dbh.LogTransactionWithTx(senderID, receiverID, amount, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func CreateWallet(dbh *database.DBHandler) (*models.Wallet, error) {
	walletID, err := generateWalletID()
	if err != nil {
		return nil, err
	}

	newWallet := &models.Wallet{
		ID:      walletID,
		Balance: 100.0,
	}

	err = dbh.InsertWalletToDB(newWallet)
	if err != nil {
		return nil, err
	}

	return newWallet, nil
}

func generateWalletID() (string, error) {
	idBytes := make([]byte, 16)
	_, err := rand.Read(idBytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(idBytes), nil
}
